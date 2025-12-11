package cmd

import (
	"fmt"
	"os"

	"linea/internal"
)

// HelpCommand displays help information for a YAML command file (supports single or multiple commands)
func HelpCommand(yamlFile string) error {
	configs, err := internal.ParseMultiYAML(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to parse YAML file: %w", err)
	}

	if len(configs) == 1 {
		config := configs[0]
		fmt.Printf("Command: %s\n", config.Command)
		if config.Subcommand != "" {
			fmt.Printf("Subcommand: %s\n", config.Subcommand)
		}
		if len(config.Args) > 0 {
			fmt.Printf("Arguments: %v\n", config.Args)
		}
		if len(config.Variables) > 0 {
			fmt.Printf("Variables:\n")
			for key, value := range config.Variables {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		cmd, err := internal.BuildCommand(config, nil)
		if err != nil {
			return err
		}
		fmt.Printf("\nFull command: %s\n", internal.FormatCommand(cmd))
		return nil
	}

	// Multiple commands
	fmt.Printf("Found %d commands in YAML file:\n\n", len(configs))
	for i, config := range configs {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Command %d/%d:\n", i+1, len(configs))
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Command: %s\n", config.Command)
		if config.Subcommand != "" {
			fmt.Printf("Subcommand: %s\n", config.Subcommand)
		}
		if len(config.Args) > 0 {
			fmt.Printf("Arguments: %v\n", config.Args)
		}
		if len(config.Variables) > 0 {
			fmt.Printf("Variables:\n")
			for key, value := range config.Variables {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}

		cmd, err := internal.BuildCommand(config, nil)
		if err != nil {
			return fmt.Errorf("error building command %d: %w", i+1, err)
		}
		fmt.Printf("Full command: %s\n", internal.FormatCommand(cmd))
		if i < len(configs)-1 {
			fmt.Println()
		}
	}

	return nil
}

// HelpCommandMain is the entry point for the help subcommand
func HelpCommandMain(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no YAML file specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea help <yaml-file>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    linea help config.yml\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	yamlFile := args[0]
	if err := HelpCommand(yamlFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

