package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// LineashContext holds the execution context for lineash scripts
type LineashContext struct {
	Variables    map[string]string
	WorkflowsDir string
	ScriptDir    string
	LineaPath    string
}

// NewLineashContext creates a new lineash context
func NewLineashContext(scriptPath string) (*LineashContext, error) {
	scriptDir := filepath.Dir(scriptPath)
	
	// Find .linea/workflows directory by walking up from script
	var workflowsDir string
	currentDir := scriptDir
	
	for {
		potentialDir := filepath.Join(currentDir, ".linea", "workflows")
		if info, err := os.Stat(potentialDir); err == nil && info.IsDir() {
			workflowsDir = potentialDir
			break
		}
		
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached root
			break
		}
		currentDir = parent
	}
	
	if workflowsDir == "" {
		return nil, fmt.Errorf("could not find .linea/workflows directory")
	}
	
	// Find linea executable - try multiple locations
	lineaExe := "linea"
	if runtime.GOOS == "windows" {
		lineaExe = "linea.exe"
	}
	
	lineaPath, err := findLineaExecutable(scriptPath, lineaExe)
	if err != nil {
		return nil, fmt.Errorf("linea executable not found: %w", err)
	}
	
	return &LineashContext{
		Variables:    make(map[string]string),
		WorkflowsDir: workflowsDir,
		ScriptDir:    scriptDir,
		LineaPath:    lineaPath,
	}, nil
}

// findLineaExecutable searches for the linea executable in multiple locations
func findLineaExecutable(scriptPath, lineaExe string) (string, error) {
	scriptDir := filepath.Dir(scriptPath)
	
	// 1. Try PATH
	if path, err := exec.LookPath(lineaExe); err == nil {
		return path, nil
	}
	
	// 2. Try relative to script directory (../../bin/linea.exe)
	potentialPath := filepath.Join(scriptDir, "..", "..", "bin", lineaExe)
	if absPath, err := filepath.Abs(potentialPath); err == nil {
		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}
	
	// 3. Try in same directory as lineash (if lineash is in bin/)
	// Get the executable path
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		potentialPath = filepath.Join(execDir, lineaExe)
		if _, err := os.Stat(potentialPath); err == nil {
			return potentialPath, nil
		}
	}
	
	// 4. Try current working directory/bin/
	if cwd, err := os.Getwd(); err == nil {
		potentialPath = filepath.Join(cwd, "bin", lineaExe)
		if _, err := os.Stat(potentialPath); err == nil {
			return potentialPath, nil
		}
	}
	
	return "", fmt.Errorf("could not find %s in PATH or common locations", lineaExe)
}

// GetAvailableWorkflows returns a list of available workflow names
func (ctx *LineashContext) GetAvailableWorkflows() ([]string, error) {
	entries, err := os.ReadDir(ctx.WorkflowsDir)
	if err != nil {
		return nil, err
	}
	
	var workflows []string
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".yml") || strings.HasSuffix(entry.Name(), ".yaml")) {
			// Remove extension to get workflow name
			name := strings.TrimSuffix(entry.Name(), ".yml")
			name = strings.TrimSuffix(name, ".yaml")
			workflows = append(workflows, name)
		}
	}
	
	return workflows, nil
}

// IsWorkflowCommand checks if a command is an available workflow
func (ctx *LineashContext) IsWorkflowCommand(cmdName string) bool {
	workflows, err := ctx.GetAvailableWorkflows()
	if err != nil {
		return false
	}
	
	for _, wf := range workflows {
		if wf == cmdName {
			return true
		}
	}
	
	return false
}

// ExecuteWorkflowCommand executes a workflow command via Linea
func (ctx *LineashContext) ExecuteWorkflowCommand(workflowName string, args []string) error {
	workflowFile := filepath.Join(ctx.WorkflowsDir, workflowName+".yml")
	
	// Check if .yaml extension exists
	if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
		workflowFile = filepath.Join(ctx.WorkflowsDir, workflowName+".yaml")
		if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
			return fmt.Errorf("workflow %s not found", workflowName)
		}
	}
	
	// Build linea command: linea run <workflow-file> [remaining args]
	// The args are already parsed and have quotes stripped by parseCommand
	// So we can pass them directly
	lineaArgs := []string{"run", workflowFile}
	lineaArgs = append(lineaArgs, args...)
	
	execCmd := exec.Command(ctx.LineaPath, lineaArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	
	return execCmd.Run()
}

// ExecuteSystemCommand executes a command via system shell
func (ctx *LineashContext) ExecuteSystemCommand(cmdLine string) error {
	// For echo command, strip quotes from arguments to avoid escaped quotes in output
	if strings.HasPrefix(strings.TrimSpace(cmdLine), "echo ") {
		cmdLine = stripEchoQuotes(cmdLine)
	}
	
	if runtime.GOOS == "windows" {
		execCmd := exec.Command("cmd.exe", "/c", cmdLine)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin
		return execCmd.Run()
	}
	
	// Unix-like systems
	execCmd := exec.Command("sh", "-c", cmdLine)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	return execCmd.Run()
}

// stripEchoQuotes removes quotes from echo command arguments
func stripEchoQuotes(cmdLine string) string {
	// Parse the echo command and rebuild without quotes
	parts := parseCommand(cmdLine)
	if len(parts) == 0 {
		return cmdLine
	}
	
	// Rebuild command: echo arg1 arg2 arg3 (without quotes)
	result := parts[0] // "echo"
	for i := 1; i < len(parts); i++ {
		result += " " + parts[i]
	}
	
	return result
}

// SubstituteVariables replaces $variable references in a string
func (ctx *LineashContext) SubstituteVariables(line string) string {
	result := line
	
	// Sort variables by length (longest first) to avoid partial replacements
	type varEntry struct {
		key   string
		value string
	}
	vars := make([]varEntry, 0, len(ctx.Variables))
	for key, value := range ctx.Variables {
		vars = append(vars, varEntry{key, value})
	}
	
	// Sort by key length descending
	for i := 0; i < len(vars); i++ {
		for j := i + 1; j < len(vars); j++ {
			if len(vars[i].key) < len(vars[j].key) {
				vars[i], vars[j] = vars[j], vars[i]
			}
		}
	}
	
	// Replace variables
	for _, v := range vars {
		// Replace ${VAR} first (more specific, handles braces)
		placeholder2 := "${" + v.key + "}"
		result = strings.ReplaceAll(result, placeholder2, v.value)
		
		// Replace $VAR (but not if it's part of a longer variable name)
		placeholder1 := "$" + v.key
		for {
			idx := strings.Index(result, placeholder1)
			if idx == -1 {
				break
			}
			
			// Check if it's a valid variable reference (not part of a longer variable)
			afterIdx := idx + len(placeholder1)
			if afterIdx >= len(result) {
				// End of string, valid replacement
				result = result[:idx] + v.value + result[afterIdx:]
			} else {
				nextChar := result[afterIdx]
				// Valid if next char is not alphanumeric or underscore
				if !((nextChar >= 'a' && nextChar <= 'z') || 
					 (nextChar >= 'A' && nextChar <= 'Z') || 
					 (nextChar >= '0' && nextChar <= '9') || 
					 nextChar == '_') {
					result = result[:idx] + v.value + result[afterIdx:]
				} else {
					// Skip this occurrence, it's part of a longer variable
					// Move past the $ to avoid infinite loop
					result = result[:idx+1] + result[idx+1:]
				}
			}
		}
	}
	
	return result
}

// ExecuteLines executes script lines with bash-like control flow using a simple parser
func ExecuteLines(ctx *LineashContext, scriptContent string) error {
	lines := strings.Split(scriptContent, "\n")
	
	// Simple bash-like parser
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			i++
			continue
		}
		
		// Skip shebang
		if strings.HasPrefix(line, "#!/") {
			i++
			continue
		}
		
		// Handle variable assignment: VAR=value
		if key, value, ok := parseVariableAssignment(line); ok {
			ctx.Variables[key] = value
			i++
			continue
		}
		
		// Handle if statement
		if strings.HasPrefix(line, "if ") {
			i = handleIfStatement(ctx, lines, i)
			continue
		}
		
		// Handle for loop
		if strings.HasPrefix(line, "for ") {
			i = handleForLoop(ctx, lines, i)
			continue
		}
		
		// Substitute variables in line BEFORE parsing
		// This ensures $variables in arguments are properly substituted
		line = ctx.SubstituteVariables(line)
		
		// Parse command (after variable substitution)
		parts := parseCommand(line)
		if len(parts) == 0 {
			i++
			continue
		}
		
		cmdName := parts[0]
		args := parts[1:]
		
		// Check if it's a workflow command
		if ctx.IsWorkflowCommand(cmdName) {
			// For workflow commands, args are already substituted
			// They will be passed as-is to linea run command
			if err := ctx.ExecuteWorkflowCommand(cmdName, args); err != nil {
				return fmt.Errorf("error executing workflow at line %d: %w", i+1, err)
			}
		} else {
			// Execute as system command
			if err := ctx.ExecuteSystemCommand(line); err != nil {
				return fmt.Errorf("error executing command at line %d: %w", i+1, err)
			}
		}
		
		i++
	}
	
	return nil
}

// parseVariableAssignment parses variable assignment: VAR=value
func parseVariableAssignment(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, "=") {
		return "", "", false
	}
	
	// Don't parse if it's part of a command (like -s/--set var=value)
	if strings.Contains(line, "-s") || strings.Contains(line, "--set") || strings.Contains(line, "--args") {
		return "", "", false
	}
	
	// Check if it's a standalone assignment (no spaces before =)
	eqIndex := strings.Index(line, "=")
	if eqIndex > 0 && !strings.Contains(line[:eqIndex], " ") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"'")
			return key, value, true
		}
	}
	
	return "", "", false
}

// parseCommand parses a command line into parts, handling quotes
func parseCommand(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	
	for i := 0; i < len(line); i++ {
		char := line[i]
		
		if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(char)
			}
		} else if char == ' ' && !inQuotes {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(char)
		}
	}
	
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	
	return parts
}

// handleIfStatement handles if/else/fi blocks
func handleIfStatement(ctx *LineashContext, lines []string, startIndex int) int {
	line := strings.TrimSpace(lines[startIndex])
	condition := strings.TrimSpace(strings.TrimPrefix(line, "if "))
	
	// Remove brackets if present
	condition = strings.Trim(condition, "[]")
	condition = strings.TrimSpace(condition)
	
	// Evaluate condition
	conditionMet := evaluateCondition(ctx, condition)
	
	// Find matching fi
	fiIndex := findMatchingFi(lines, startIndex)
	
	if conditionMet {
		// Execute if block
		for i := startIndex + 1; i < fiIndex; i++ {
			line := strings.TrimSpace(lines[i])
			
			// Skip "then" keyword
			if line == "then" {
				continue
			}
			
			if line == "else" {
				// Skip else block
				break
			}
			
			if line == "fi" {
				break
			}
			
			// Execute line
			executeLine(ctx, line, i)
		}
		} else {
			// Find else block if present
			elseIndex := -1
			for i := startIndex + 1; i < fiIndex; i++ {
				if strings.TrimSpace(lines[i]) == "else" {
					elseIndex = i
					break
				}
			}
			
			if elseIndex > 0 {
				// Execute else block
				for i := elseIndex + 1; i < fiIndex; i++ {
					line := strings.TrimSpace(lines[i])
					if line == "fi" {
						break
					}
					executeLine(ctx, line, i)
				}
			}
		}
	
	return fiIndex + 1
}

// handleForLoop handles for loops
func handleForLoop(ctx *LineashContext, lines []string, startIndex int) int {
	line := strings.TrimSpace(lines[startIndex])
	// Parse: for VAR in value1 value2 value3; do
	parts := strings.Fields(line)
	if len(parts) < 4 || parts[1] == "" || parts[2] != "in" {
		return startIndex + 1
	}
	
	varName := parts[1]
	values := parts[3:]
	
	// Remove "do" if present in the same line
	if len(values) > 0 && values[len(values)-1] == "do" {
		values = values[:len(values)-1]
	}
	
	// Find matching done
	doneIndex := findMatchingDone(lines, startIndex)
	
	// Find "do" keyword (might be on next line)
	doIndex := startIndex + 1
	for doIndex < doneIndex && strings.TrimSpace(lines[doIndex]) != "do" {
		doIndex++
	}
	if doIndex >= doneIndex {
		doIndex = startIndex + 1 // If no "do" found, start after for line
	} else {
		doIndex++ // Skip the "do" line
	}
	
	// Execute loop body for each value
	for _, value := range values {
		ctx.Variables[varName] = strings.Trim(value, "\"'")
		
		// Execute loop body (from doIndex to doneIndex)
		for i := doIndex; i < doneIndex; i++ {
			line := strings.TrimSpace(lines[i])
			
			// Skip "do" keyword
			if line == "do" {
				continue
			}
			
			if line == "done" {
				break
			}
			if err := executeLine(ctx, line, i); err != nil {
				// Continue on error for now
			}
		}
	}
	
	return doneIndex + 1
}

// executeLine executes a single line
func executeLine(ctx *LineashContext, line string, lineNum int) error {
	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}
	
	// Handle variable assignment
	if key, value, ok := parseVariableAssignment(line); ok {
		ctx.Variables[key] = value
		return nil
	}
	
	// Substitute variables
	line = ctx.SubstituteVariables(line)
	
	// Parse command
	parts := parseCommand(line)
	if len(parts) == 0 {
		return nil
	}
	
	cmdName := parts[0]
	args := parts[1:]
	
	// Check if it's a workflow command
	if ctx.IsWorkflowCommand(cmdName) {
		return ctx.ExecuteWorkflowCommand(cmdName, args)
	}
	
	// Execute as system command
	return ctx.ExecuteSystemCommand(line)
}

// evaluateCondition evaluates a bash-like condition
func evaluateCondition(ctx *LineashContext, condition string) bool {
	condition = strings.TrimSpace(condition)
	
	// Simple string comparison: [ "$VAR" = "value" ]
	if strings.Contains(condition, "=") {
		parts := strings.SplitN(condition, "=", 2)
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			
			// Remove quotes and $ from left
			left = strings.Trim(left, "\"'")
			if strings.HasPrefix(left, "$") {
				left = strings.TrimPrefix(left, "$")
				left = strings.Trim(left, "{}")
			}
			
			// Remove quotes from right
			right = strings.Trim(right, "\"'")
			
			// Get variable value
			leftVal := ctx.Variables[left]
			if leftVal == "" && !strings.HasPrefix(parts[0], "$") {
				leftVal = left
			}
			
			return leftVal == right
		}
	}
	
	// Check if variable is set: [ -n "$VAR" ] or [ "$VAR" ]
	if strings.HasPrefix(condition, "-n ") {
		varName := strings.TrimSpace(strings.TrimPrefix(condition, "-n"))
		varName = strings.Trim(varName, "\"'$")
		varName = strings.Trim(varName, "{}")
		val := ctx.Variables[varName]
		return val != ""
	}
	
	// Simple variable check
	if strings.HasPrefix(condition, "$") {
		varName := strings.Trim(condition, "$\"'")
		varName = strings.Trim(varName, "{}")
		val := ctx.Variables[varName]
		return val != ""
	}
	
	return false
}

// findMatchingFi finds the matching 'fi' for an 'if' statement
func findMatchingFi(lines []string, ifIndex int) int {
	depth := 1
	for i := ifIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "if ") {
			depth++
		}
		if line == "fi" {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(lines)
}

// findMatchingDone finds the matching 'done' for a 'for' statement
func findMatchingDone(lines []string, forIndex int) int {
	depth := 1
	for i := forIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "for ") {
			depth++
		}
		if line == "done" {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(lines)
}
