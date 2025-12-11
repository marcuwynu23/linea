package cmd

import (
	"fmt"
	"os"
	"strings"

	"linea/internal"
)

// TestCommand performs a dry-run of a YAML command file (supports single or multiple commands)
func TestCommand(yamlFile string, overrideVars map[string]string) error {
	configs, err := internal.ParseMultiYAML(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	if len(configs) == 1 {
		cmd, err := internal.BuildCommand(configs[0], overrideVars)
		if err != nil {
			return err
		}
		internal.DryRun(cmd)
		return nil
	}

	// Multiple commands
	fmt.Printf("Found %d commands in YAML file:\n\n", len(configs))
	for i, config := range configs {
		fmt.Printf("[%d/%d] ", i+1, len(configs))
		cmd, err := internal.BuildCommand(config, overrideVars)
		if err != nil {
			return fmt.Errorf("error building command %d: %w", i+1, err)
		}
		internal.DryRun(cmd)
		if i < len(configs)-1 {
			fmt.Println()
		}
	}

	return nil
}

// TestCommandMain is the entry point for the test subcommand
func TestCommandMain(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no YAML file specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea test [options] <yaml-file>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "    --args <var>=<value>       Provide variable values for testing\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    linea test config.yml\n")
		fmt.Fprintf(os.Stderr, "    linea test config.yml --args variable=\"test\"\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	// Parse --args flags
	overrideVars, remainingArgs := ParseArgs(args)
	
	yamlFile := ""
	for _, arg := range remainingArgs {
		if !strings.HasPrefix(arg, "-") {
			yamlFile = arg
			break
		}
	}

	if yamlFile == "" {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no YAML file specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea test [options] <yaml-file>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  OPTIONS:\n")
		fmt.Fprintf(os.Stderr, "    --args <var>=<value>       Provide variable values for testing\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	if err := TestCommand(yamlFile, overrideVars); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

