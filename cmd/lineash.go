package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"linea/internal"
)

// ExecuteLineashScript executes a .lnsh script file with bash-like features
func ExecuteLineashScript(scriptPath string) error {
	// Check if file exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script file not found: %s", scriptPath)
	}

	// Create lineash context
	ctx, err := internal.NewLineashContext(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to initialize lineash context: %w", err)
	}

	// Read script content
	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read script: %w", err)
	}

	// Execute script with bash-like features
	return internal.ExecuteLines(ctx, string(scriptContent))
}

// LineashMain is the entry point for the lineash command
func LineashMain(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no script file specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    lineash <script.lnsh>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    lineash scripts/script.lnsh\n")
		fmt.Fprintf(os.Stderr, "    lineash .linea/scripts/deploy.lnsh\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  NOTE:\n")
		fmt.Fprintf(os.Stderr, "    Scripts must be in a directory with .linea/workflows/ available\n")
		fmt.Fprintf(os.Stderr, "    Workflows in .linea/workflows/ can be called as commands\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	scriptPath := args[0]

	// Resolve absolute path
	if !filepath.IsAbs(scriptPath) {
		cwd, err := os.Getwd()
		if err == nil {
			scriptPath = filepath.Join(cwd, scriptPath)
		}
	}

	if err := ExecuteLineashScript(scriptPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

