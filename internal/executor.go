package internal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// BuildCommand constructs the full command with subcommand and arguments
func BuildCommand(config *CommandConfig, overrideVars map[string]string) ([]string, error) {
	// Merge YAML variables with override variables (override takes precedence)
	mergedVars := make(map[string]string)
	if config.Variables != nil {
		for k, v := range config.Variables {
			mergedVars[k] = v
		}
	}
	if overrideVars != nil {
		for k, v := range overrideVars {
			mergedVars[k] = v
		}
	}
	
	// Collect all strings that need validation (args + variable values)
	stringsToValidate := make([]string, 0, len(config.Args))
	stringsToValidate = append(stringsToValidate, config.Args...)
	for _, v := range mergedVars {
		stringsToValidate = append(stringsToValidate, v)
	}
	
	// Validate that all referenced variables are defined
	if err := ValidateVariables(stringsToValidate, mergedVars); err != nil {
		return nil, err
	}
	
	cmd := []string{config.Command}
	
	if config.Subcommand != "" {
		cmd = append(cmd, config.Subcommand)
	}
	
	// Apply variable substitution to arguments
	args := SubstituteVariablesInArgs(config.Args, mergedVars)
	cmd = append(cmd, args...)
	
	return cmd, nil
}

// FormatCommand returns a string representation of the command for display
func FormatCommand(cmd []string) string {
	return strings.Join(cmd, " ")
}

// ExecuteCommand runs the command and returns the output
func ExecuteCommand(cmd []string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command is empty")
	}

	// On Windows, check if command exists in PATH
	// If not, try executing through cmd.exe (for shell built-ins like echo, dir, etc.)
	if runtime.GOOS == "windows" {
		_, err := exec.LookPath(cmd[0])
		if err != nil {
			// Command not found in PATH, try shell execution
			return executeWindowsShell(cmd)
		}
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

// ExecuteMultipleCommands executes multiple commands sequentially
// Stops on first error unless continueOnError is true
func ExecuteMultipleCommands(configs []*CommandConfig, overrideVars map[string]string, continueOnError bool, verbose bool) error {
	for i, config := range configs {
		if verbose {
			fmt.Printf("\n[%d/%d] ", i+1, len(configs))
		}

		cmd, err := BuildCommand(config, overrideVars)
		if err != nil {
			if continueOnError {
				fmt.Fprintf(os.Stderr, "Error building command %d: %v\n", i+1, err)
				continue
			}
			return fmt.Errorf("error building command %d: %w", i+1, err)
		}

		if verbose {
			fmt.Printf("Executing: %s\n", FormatCommand(cmd))
		}

		if err := ExecuteCommand(cmd); err != nil {
			if continueOnError {
				fmt.Fprintf(os.Stderr, "Error executing command %d: %v\n", i+1, err)
				continue
			}
			return fmt.Errorf("command %d execution failed: %w", i+1, err)
		}
	}

	return nil
}

// executeWindowsShell executes a command through cmd.exe on Windows
// This is used for shell built-ins like echo, dir, etc.
func executeWindowsShell(cmd []string) error {
	// Build the command string for cmd.exe /c
	cmdStr := FormatCommand(cmd)
	
	execCmd := exec.Command("cmd.exe", "/c", cmdStr)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	return execCmd.Run()
}

// DryRun prints the command without executing it
func DryRun(cmd []string) {
	fmt.Println("Dry run - would execute:")
	fmt.Println(FormatCommand(cmd))
}

