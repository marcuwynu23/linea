package cmd

import (
	"fmt"
	"os"
	"strings"
)

// InitCommand creates a new workflow YAML file with template and documentation
func InitCommand(yamlFile string) error {
	// Check if file already exists
	if _, err := os.Stat(yamlFile); err == nil {
		return fmt.Errorf("file %s already exists", yamlFile)
	}

	// Generate template content
	template := `# Linea Workflow Configuration
# This file defines commands that can be executed using: linea run <this-file>

# Main command to execute
command: echo

# Optional subcommand (for commands like: docker ps, git status, etc.)
# subcommand: ps

# Arguments to pass to the command
args:
  - "Hello, Linea!"
  - "This is a template workflow file"

# Variables for substitution
# Use {variable} or $variable syntax in args
variables:
  # message: "Custom message"
  # name: "Your Name"

# Example usage:
#   linea run <this-file>
#   linea run <this-file> --args message="Custom"
#   linea test <this-file>
#   linea help <this-file>

# Multiple Commands:
# You can define multiple commands in one file by separating them with ---
# 
# command: echo
# args:
#   - "First command"
# ---
# command: echo
# args:
#   - "Second command"
`

	// Write template to file
	err := os.WriteFile(yamlFile, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	fmt.Printf("✅ Created workflow file: %s\n", yamlFile)
	fmt.Printf("\n")
	fmt.Printf("You can now:\n")
	fmt.Printf("  • Edit the file to customize your workflow\n")
	fmt.Printf("  • Test it: linea test %s\n", yamlFile)
	fmt.Printf("  • Run it: linea run %s\n", yamlFile)
	fmt.Printf("\n")

	return nil
}

// InitCommandMain is the entry point for the init subcommand
func InitCommandMain(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no file name specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea init <file-name>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    linea init workflow.yml\n")
		fmt.Fprintf(os.Stderr, "    linea init my-commands.yml\n")
		fmt.Fprintf(os.Stderr, "    linea init examples/new-workflow.yml\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	yamlFile := args[0]

	// Ensure file has .yml or .yaml extension
	if !strings.HasSuffix(yamlFile, ".yml") && !strings.HasSuffix(yamlFile, ".yaml") {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ⚠️  Warning: file should have .yml or .yaml extension\n")
		fmt.Fprintf(os.Stderr, "  Continuing anyway...\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	if err := InitCommand(yamlFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

