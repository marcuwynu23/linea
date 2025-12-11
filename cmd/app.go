package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// AppCreateCommand creates a new Linea App folder structure
func AppCreateCommand(appName string) error {
	// Check if directory already exists
	if _, err := os.Stat(appName); err == nil {
		return fmt.Errorf("directory %s already exists", appName)
	}

	// Create directory structure
	workflowsDir := filepath.Join(appName, ".linea", "workflows")
	scriptsDir := filepath.Join(appName, "scripts")

	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create scripts directory: %w", err)
	}

	// Create example workflow files
	createVMWorkflow := `# Create VM Workflow
# Usage: linea run .linea/workflows/create-vm.yml -s name="vm-name"

command: echo
args:
  - "Creating VM: {name}"
variables:
  name: "default-vm"
`

	lsWorkflow := `# List Directory Workflow
# Usage: linea run .linea/workflows/ls.yml

command: ls
args:
  - -l
  - -a
`

	// Write workflow files
	if err := os.WriteFile(filepath.Join(workflowsDir, "create-vm.yml"), []byte(createVMWorkflow), 0644); err != nil {
		return fmt.Errorf("failed to create create-vm.yml: %w", err)
	}

	if err := os.WriteFile(filepath.Join(workflowsDir, "ls.yml"), []byte(lsWorkflow), 0644); err != nil {
		return fmt.Errorf("failed to create ls.yml: %w", err)
	}

	// Create example script
	exampleScript := `#!/bin/lineash
# Linea Script Example with bash-like features
# This script demonstrates variables, conditionals, and loops
# Note: Use $variable syntax in lineash (not {variable} which is for YAML)

# Variables
VM_NAME="my-vm"
VM_OS="alpine"

echo "Starting VM creation..."

# Conditional execution
if [ "$VM_OS" = "alpine" ]
then
    echo "Using Alpine Linux"
    # Pass variables to workflows using $variable syntax
    create-vm -s name="$VM_NAME"
else
    echo "Using different OS"
fi

# For loop
for item in workflows scripts
do
    echo "Checking $item..."
    ls
done

echo "Script completed!"
`

	if err := os.WriteFile(filepath.Join(scriptsDir, "script.lnsh"), []byte(exampleScript), 0755); err != nil {
		return fmt.Errorf("failed to create script.lnsh: %w", err)
	}

	// Create README
	readme := "# " + appName + "\n\n" +
		"This is a Linea App directory structure.\n\n" +
		"## Directory Structure\n\n" +
		"- `.linea/workflows/` - Workflow YAML files that can be executed as commands\n" +
		"- `scripts/` - Lineash scripts (`.lnsh` files) that can use workflows as commands\n\n" +
		"## Usage\n\n" +
		"### Running Workflows\n\n" +
		"```bash\n" +
		"# Run a workflow directly\n" +
		"linea run .linea/workflows/create-vm.yml -s name=\"my-vm\"\n\n" +
		"# Or use lineash to run workflows as commands\n" +
		"lineash scripts/script.lnsh\n" +
		"```\n\n" +
		"### Creating New Workflows\n\n" +
		"1. Create a new YAML file in `.linea/workflows/`\n" +
		"2. Define your command structure\n" +
		"3. Use it in scripts or run directly with `linea run`\n\n" +
		"### Writing Scripts\n\n" +
		"Scripts in `scripts/` can:\n" +
		"- Execute workflows as commands (if they exist in `.linea/workflows/`)\n" +
		"- Use bash-like syntax (variables, conditions, loops)\n" +
		"- Call system commands\n\n" +
		"Example:\n" +
		"```bash\n" +
		"#!/bin/lineash\n" +
		"echo \"Hello\"\n" +
		"create-vm -s name=\"test\"\n" +
		"ls\n" +
		"```\n"

	if err := os.WriteFile(filepath.Join(appName, "README.md"), []byte(readme), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	fmt.Printf("✅ Created Linea App: %s\n", appName)
	fmt.Printf("\n")
	fmt.Printf("Directory structure:\n")
	fmt.Printf("  %s/\n", appName)
	fmt.Printf("  ├─ .linea/workflows/\n")
	fmt.Printf("  │   ├─ create-vm.yml\n")
	fmt.Printf("  │   └─ ls.yml\n")
	fmt.Printf("  ├─ scripts/\n")
		fmt.Printf("  │   └─ script.lnsh\n")
	fmt.Printf("  └─ README.md\n")
	fmt.Printf("\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  • Edit workflows in .linea/workflows/\n")
	fmt.Printf("  • Create scripts in scripts/\n")
		fmt.Printf("  • Run scripts: lineash scripts/script.lnsh\n")
	fmt.Printf("\n")

	return nil
}

// AppCreateCommandMain is the entry point for the app create subcommand
func AppCreateCommandMain(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: no app name specified\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea app create <app-name>\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  EXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    linea app create my-app\n")
		fmt.Fprintf(os.Stderr, "    linea app create deployment\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	if args[0] != "create" {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  ❌ Error: unknown app subcommand '%s'\n", args[0])
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  USAGE:\n")
		fmt.Fprintf(os.Stderr, "    linea app create <app-name>\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	appName := args[1]

	if err := AppCreateCommand(appName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

