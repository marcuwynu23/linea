package cmd

import (
	"fmt"
	"os"
	"strings"

	"linea/internal"
)

// RunCommand executes a YAML command file (supports single or multiple commands)
func RunCommand(yamlFile string, verbose bool, overrideVars map[string]string) error {
	configs, err := internal.ParseMultiYAML(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	// If single command, execute normally for backward compatibility
	if len(configs) == 1 {
		cmd, err := internal.BuildCommand(configs[0], overrideVars)
		if err != nil {
			return err
		}
		
		if verbose {
			fmt.Printf("Executing: %s\n", internal.FormatCommand(cmd))
		}
		
		if err := internal.ExecuteCommand(cmd); err != nil {
			return fmt.Errorf("command execution failed: %w", err)
		}
		return nil
	}

	// Multiple commands - execute sequentially
	if verbose {
		fmt.Printf("Found %d commands in YAML file\n", len(configs))
	}

	return internal.ExecuteMultipleCommands(configs, overrideVars, false, verbose)
}

// ParseArgs parses -s/--set flags from command line arguments
// Format: -s variable="value" or --set variable=value
// Also supports --args for backward compatibility
func ParseArgs(args []string) (map[string]string, []string) {
	vars := make(map[string]string)
	remainingArgs := []string{}
	
	i := 0
	for i < len(args) {
		if args[i] == "-s" || args[i] == "--set" || args[i] == "--args" {
			if i+1 < len(args) {
				argPair := args[i+1]
				// Parse variable=value format
				parts := strings.SplitN(argPair, "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]
					// Remove quotes if present
					value = strings.Trim(value, "\"'")
					vars[key] = value
				}
				i += 2
			} else {
				i++
			}
		} else {
			remainingArgs = append(remainingArgs, args[i])
			i++
		}
	}
	
	return vars, remainingArgs
}

// RunCommandMain is the entry point for the run subcommand
func RunCommandMain(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no YAML file specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea run [options] <yaml-file>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "    -v, --verbose              Show the command before executing\n")
		fmt.Fprintf(os.Stderr, "    -s, --set <var>=<value>     Set variable values (can be used multiple times)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    linea run config.yml\n")
		fmt.Fprintf(os.Stderr, "    linea run -v config.yml\n")
		fmt.Fprintf(os.Stderr, "    linea run config.yml -s name=\"John\"\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	// Parse -s/--set flags first
	overrideVars, remainingArgs := ParseArgs(args)
	
	verbose := false
	yamlFile := ""
	
	// Parse other flags
	for _, arg := range remainingArgs {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else if !strings.HasPrefix(arg, "-") {
			yamlFile = arg
		}
	}

	if yamlFile == "" {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no YAML file specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea run [options] <yaml-file>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "    -v, --verbose              Show the command before executing\n")
		fmt.Fprintf(os.Stderr, "    -s, --set <var>=<value>     Set variable values (can be used multiple times)\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	if err := RunCommand(yamlFile, verbose, overrideVars); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

