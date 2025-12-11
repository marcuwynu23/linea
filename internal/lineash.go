package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// LineashContext holds the execution context for lineash scripts
type LineashContext struct {
	Variables    map[string]string
	WorkflowsDir string
	ScriptDir    string
	LineaPath    string
	Args         []string // Positional parameters $1, $2, etc.
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
		Args:         []string{}, // Empty by default, can be set if args are passed
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

// SubstituteVariables replaces $variable references, positional parameters, and arithmetic expressions
func (ctx *LineashContext) SubstituteVariables(line string) string {
	result := line
	
	// First, handle arithmetic expressions $((...))
	result = substituteArithmetic(result, ctx)
	
	// Handle positional parameters $1, $2, etc.
	result = substitutePositionalParams(result, ctx)
	
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

// substitutePositionalParams replaces $1, $2, etc. with actual arguments
func substitutePositionalParams(line string, ctx *LineashContext) string {
	result := line
	
	// Match $1, $2, etc. (handle multi-digit numbers properly)
	re := regexp.MustCompile(`\$(\d+)`)
	matches := re.FindAllStringSubmatch(result, -1)
	
	for _, match := range matches {
		paramNum, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		
		// $1 is first arg, $2 is second, etc.
		var replacement string
		if paramNum > 0 && paramNum <= len(ctx.Args) {
			replacement = ctx.Args[paramNum-1]
		} else {
			replacement = "" // Undefined parameter
		}
		
		result = strings.Replace(result, match[0], replacement, 1)
	}
	
	return result
}

// substituteArithmetic replaces $((expression)) with evaluated result
func substituteArithmetic(line string, ctx *LineashContext) string {
	result := line
	
	// Match $((...)) expressions
	re := regexp.MustCompile(`\$\(\(([^)]+)\)\)`)
	matches := re.FindAllStringSubmatch(result, -1)
	
	for _, match := range matches {
		expr := strings.TrimSpace(match[1])
		
		// Substitute variables in expression first
		// Handle $variable and ${variable} syntax
		expr = substitutePositionalParams(expr, ctx)
		
		// Substitute variables - need to handle variable names that might be part of the expression
		// Sort by length to avoid partial matches
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
		
		// Replace variables (without $ prefix in arithmetic expressions)
		for _, v := range vars {
			// Replace variable name with its value
			// Use word boundaries to avoid partial matches
			expr = regexp.MustCompile(`\b`+regexp.QuoteMeta(v.key)+`\b`).ReplaceAllString(expr, v.value)
		}
		
		// Also handle $variable syntax in arithmetic
		for _, v := range vars {
			expr = strings.ReplaceAll(expr, "$"+v.key, v.value)
			expr = strings.ReplaceAll(expr, "${"+v.key+"}", v.value)
		}
		
		// Evaluate arithmetic expression
		value := evaluateArithmetic(expr)
		result = strings.Replace(result, match[0], value, 1)
	}
	
	return result
}

// evaluateArithmetic evaluates a simple arithmetic expression
func evaluateArithmetic(expr string) string {
	expr = strings.TrimSpace(expr)
	
	// Try to parse as integer
	if val, err := strconv.Atoi(expr); err == nil {
		return strconv.Itoa(val)
	}
	
	// Handle basic arithmetic: +, -, *, /, %
	// Split by operators while preserving them
	parts := []string{}
	current := ""
	
	for i := 0; i < len(expr); i++ {
		char := expr[i]
		if char == '+' || char == '-' || char == '*' || char == '/' || char == '%' {
			if current != "" {
				parts = append(parts, strings.TrimSpace(current))
				current = ""
			}
			parts = append(parts, string(char))
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, strings.TrimSpace(current))
	}
	
	if len(parts) == 0 {
		return "0"
	}
	
	// Parse first number
	result, err := strconv.Atoi(parts[0])
	if err != nil {
		return "0"
	}
	
	// Apply operations
	for i := 1; i < len(parts); i += 2 {
		if i+1 >= len(parts) {
			break
		}
		
		op := parts[i]
		val, err := strconv.Atoi(parts[i+1])
		if err != nil {
			continue
		}
		
		switch op {
		case "+":
			result += val
		case "-":
			result -= val
		case "*":
			result *= val
		case "/":
			if val != 0 {
				result /= val
			}
		case "%":
			if val != 0 {
				result %= val
			}
		}
	}
	
	return strconv.Itoa(result)
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
			// Substitute variables in value before assignment
			value = ctx.SubstituteVariables(value)
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
		
		// Handle while loop (friendly syntax)
		if strings.HasPrefix(line, "while ") {
			i = handleWhileLoop(ctx, lines, i)
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

// handleIfStatement handles if/else/end blocks (friendly syntax)
// Also supports backward compatibility with if/then/else/fi
func handleIfStatement(ctx *LineashContext, lines []string, startIndex int) int {
	line := strings.TrimSpace(lines[startIndex])
	condition := strings.TrimSpace(strings.TrimPrefix(line, "if "))
	
	// Remove brackets if present (for backward compatibility)
	condition = strings.Trim(condition, "[]")
	condition = strings.TrimSpace(condition)
	
	// Substitute variables in condition
	condition = ctx.SubstituteVariables(condition)
	
	// Evaluate condition
	conditionMet := evaluateCondition(ctx, condition)
	
	// Find matching end or fi (for backward compatibility)
	endIndex := findMatchingEndOrFi(lines, startIndex)
	
	if conditionMet {
		// Execute if block
		for i := startIndex + 1; i < endIndex; i++ {
			line := strings.TrimSpace(lines[i])
			
			// Skip "then" keyword (for backward compatibility)
			if line == "then" {
				continue
			}
			
			if line == "else" || line == "elif" {
				// Skip else/elif block
				break
			}
			
			if line == "end" || line == "fi" {
				break
			}
			
			// Execute line
			if err := executeLine(ctx, line, i); err != nil {
				// Continue on error for now
			}
		}
	} else {
		// Find else block if present
		elseIndex := -1
		for i := startIndex + 1; i < endIndex; i++ {
			trimmed := strings.TrimSpace(lines[i])
			if trimmed == "else" {
				elseIndex = i
				break
			}
		}
		
		if elseIndex > 0 {
			// Execute else block
			for i := elseIndex + 1; i < endIndex; i++ {
				line := strings.TrimSpace(lines[i])
				if line == "end" || line == "fi" {
					break
				}
				if err := executeLine(ctx, line, i); err != nil {
					// Continue on error for now
				}
			}
		}
	}
	
	return endIndex + 1
}

// handleForLoop handles for loops with friendly syntax (for ... in ... end)
// Also supports backward compatibility with for ... do ... done
func handleForLoop(ctx *LineashContext, lines []string, startIndex int) int {
	line := strings.TrimSpace(lines[startIndex])
	// Parse: for VAR in value1 value2 value3
	parts := strings.Fields(line)
	if len(parts) < 4 || parts[1] == "" || parts[2] != "in" {
		return startIndex + 1
	}
	
	varName := parts[1]
	values := parts[3:]
	
	// Remove "do" if present in the same line (for backward compatibility)
	if len(values) > 0 && values[len(values)-1] == "do" {
		values = values[:len(values)-1]
	}
	
	// Substitute variables in values
	for i, val := range values {
		values[i] = ctx.SubstituteVariables(strings.Trim(val, "\"'"))
	}
	
	// Find matching end or done (for backward compatibility)
	endIndex := findMatchingEndOrDone(lines, startIndex)
	
	// Find body start (skip "do" if present for backward compatibility)
	bodyStart := startIndex + 1
	for bodyStart < endIndex && strings.TrimSpace(lines[bodyStart]) == "do" {
		bodyStart++
	}
	
	// Execute loop body for each value
	for _, value := range values {
		ctx.Variables[varName] = value
		
		// Execute loop body
		for i := bodyStart; i < endIndex; i++ {
			line := strings.TrimSpace(lines[i])
			
			// Skip "do" keyword (for backward compatibility)
			if line == "do" {
				continue
			}
			
			if line == "end" || line == "done" {
				break
			}
			
			if err := executeLine(ctx, line, i); err != nil {
				// Continue on error for now
			}
		}
	}
	
	return endIndex + 1
}

// handleWhileLoop handles while loops with friendly syntax (while ... end)
func handleWhileLoop(ctx *LineashContext, lines []string, startIndex int) int {
	line := strings.TrimSpace(lines[startIndex])
	condition := strings.TrimSpace(strings.TrimPrefix(line, "while "))
	
	// Find matching end
	endIndex := findMatchingEnd(lines, startIndex)
	
	// Find body start (skip "do" if present for backward compatibility)
	bodyStart := startIndex + 1
	for bodyStart < endIndex && strings.TrimSpace(lines[bodyStart]) == "do" {
		bodyStart++
	}
	
	// Execute loop while condition is true
	for {
		// Evaluate condition (substitute variables first)
		cond := ctx.SubstituteVariables(condition)
		if !evaluateCondition(ctx, cond) {
			break
		}
		
		// Execute loop body
		for i := bodyStart; i < endIndex; i++ {
			line := strings.TrimSpace(lines[i])
			
			if line == "do" {
				continue
			}
			
			if line == "end" || line == "done" {
				break
			}
			
			if err := executeLine(ctx, line, i); err != nil {
				// Continue on error for now
			}
		}
	}
	
	return endIndex + 1
}

// executeLine executes a single line
func executeLine(ctx *LineashContext, line string, lineNum int) error {
	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}
	
	// Handle variable assignment
	if key, value, ok := parseVariableAssignment(line); ok {
		// Substitute variables in value before assignment
		value = ctx.SubstituteVariables(value)
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

// evaluateCondition evaluates a condition with friendly operators (==, !=, <, >, <=, >=)
func evaluateCondition(ctx *LineashContext, condition string) bool {
	condition = strings.TrimSpace(condition)
	
	// Handle comparison operators: ==, !=, <, >, <=, >=
	// Check longer operators first to avoid matching shorter ones
	operators := []struct {
		op   string
		len  int
		fn   func(left, right string) bool
	}{
		{"<=", 2, func(l, r string) bool {
			left, err1 := strconv.Atoi(strings.TrimSpace(l))
			right, err2 := strconv.Atoi(strings.TrimSpace(r))
			if err1 == nil && err2 == nil {
				return left <= right
			}
			return strings.TrimSpace(l) <= strings.TrimSpace(r)
		}},
		{">=", 2, func(l, r string) bool {
			left, err1 := strconv.Atoi(strings.TrimSpace(l))
			right, err2 := strconv.Atoi(strings.TrimSpace(r))
			if err1 == nil && err2 == nil {
				return left >= right
			}
			return strings.TrimSpace(l) >= strings.TrimSpace(r)
		}},
		{"==", 2, func(l, r string) bool {
			return strings.TrimSpace(l) == strings.TrimSpace(r)
		}},
		{"!=", 2, func(l, r string) bool {
			return strings.TrimSpace(l) != strings.TrimSpace(r)
		}},
		{"<", 1, func(l, r string) bool {
			left, err1 := strconv.Atoi(strings.TrimSpace(l))
			right, err2 := strconv.Atoi(strings.TrimSpace(r))
			if err1 == nil && err2 == nil {
				return left < right
			}
			return strings.TrimSpace(l) < strings.TrimSpace(r)
		}},
		{">", 1, func(l, r string) bool {
			left, err1 := strconv.Atoi(strings.TrimSpace(l))
			right, err2 := strconv.Atoi(strings.TrimSpace(r))
			if err1 == nil && err2 == nil {
				return left > right
			}
			return strings.TrimSpace(l) > strings.TrimSpace(r)
		}},
		{"=", 1, func(l, r string) bool {
			return strings.TrimSpace(l) == strings.TrimSpace(r)
		}},
	}
	
	for _, opInfo := range operators {
		if idx := strings.Index(condition, opInfo.op); idx > 0 {
			left := strings.TrimSpace(condition[:idx])
			right := strings.TrimSpace(condition[idx+opInfo.len:])
			
			// Remove quotes
			left = strings.Trim(left, "\"'")
			right = strings.Trim(right, "\"'")
			
			// Variables should already be substituted by SubstituteVariables, but handle just in case
			// Handle $variable references (in case substitution didn't happen)
			if strings.HasPrefix(left, "$") {
				varName := strings.TrimPrefix(left, "$")
				varName = strings.Trim(varName, "{}")
				if val, ok := ctx.Variables[varName]; ok {
					left = val
				} else {
					left = "" // Undefined variable
				}
			}
			if strings.HasPrefix(right, "$") {
				varName := strings.TrimPrefix(right, "$")
				varName = strings.Trim(varName, "{}")
				if val, ok := ctx.Variables[varName]; ok {
					right = val
				} else {
					right = "" // Undefined variable
				}
			}
			
			return opInfo.fn(left, right)
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

// findMatchingEnd finds the matching 'end' for friendly syntax
func findMatchingEnd(lines []string, startIndex int) int {
	depth := 1
	for i := startIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// Check for nested blocks
		if strings.HasPrefix(line, "if ") || strings.HasPrefix(line, "for ") || strings.HasPrefix(line, "while ") {
			depth++
		}
		if line == "end" {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(lines)
}

// findMatchingEndOrFi finds matching 'end' or 'fi' (for backward compatibility)
func findMatchingEndOrFi(lines []string, ifIndex int) int {
	depth := 1
	for i := ifIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "if ") {
			depth++
		}
		if line == "end" || line == "fi" {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(lines)
}

// findMatchingEndOrDone finds matching 'end' or 'done' (for backward compatibility)
func findMatchingEndOrDone(lines []string, forIndex int) int {
	depth := 1
	for i := forIndex + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "for ") || strings.HasPrefix(line, "while ") {
			depth++
		}
		if line == "end" || line == "done" {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return len(lines)
}
