# Linea Features

This document provides a comprehensive overview of all features available in Linea CLI.

## Core Features

### 1. YAML-Driven Command Execution

Linea allows you to define commands entirely in YAML files, making command execution declarative and reproducible.

**Benefits:**
- Version control your commands
- Share command configurations easily
- Document complex command invocations
- Reproduce commands across different environments

### 2. Cross-Platform Support

Linea automatically handles OS-specific differences, making your YAML files work seamlessly on Windows, Linux, and macOS.

**Features:**
- Automatic path separator conversion (`/` vs `\`)
- OS-aware help flag detection (`/?` vs `--help`)
- Windows shell built-in support (echo, dir, etc.)
- Unix command compatibility

### 3. Variable Substitution

Two syntaxes are supported for variable substitution with different override behaviors:

**Curly Brace Syntax:**
```yaml
args:
  - "{message}"
```

**Dollar Sign Syntax:**
```yaml
args:
  - "$message"
```

**Variable Sources:**
- Defined in YAML `variables` section
- Provided via command-line `-s/--set` flag
- Command-line variables override YAML variables

### 4. Command-Line Variable Override

Pass variables at runtime without modifying YAML files:

```bash
linea run config.yml -s/--set variable="value"
```

**Use Cases:**
- Environment-specific configurations
- Dynamic parameter injection
- Testing with different values
- CI/CD pipeline integration

### 5. Variable Validation

Linea validates that all referenced variables are defined before execution, preventing runtime errors.

**Features:**
- Detects undefined variables in args and variable values
- Clear error messages indicating missing variables
- Suggests using `-s/--set` to provide missing values
- Validates both `{variable}` and `$variable` syntaxes

### 6. Dry-Run Mode

Test commands without executing them using the `test` subcommand:

```bash
linea test config.yml
```

**Benefits:**
- Preview commands before execution
- Debug variable substitution
- Verify command construction
- Safe command testing

### 7. Verbose Mode

Control output verbosity with the `-v` or `--verbose` flag:

```bash
linea run -v config.yml
```

**Features:**
- Shows full command before execution
- Default: silent execution (no "Executing:" message)
- Verbose: displays command preview

### 8. Help Command

Get information about commands defined in YAML files:

```bash
linea help config.yml
```

**Displays:**
- Command name
- Subcommand (if any)
- Arguments
- Variables and their values
- Full constructed command

### 9. Init Command

Quickly create new workflow files with templates:

```bash
linea init workflow.yml
```

**Features:**
- Creates template YAML file
- Includes comprehensive documentation
- Example usage comments
- Variable examples
- Multiple commands example
- Prevents overwriting existing files

**Benefits:**
- Quick start for new workflows
- Learn by example
- Consistent file structure
- Built-in documentation

### 10. Subcommand Support

Support for commands with subcommands:

```yaml
command: docker
subcommand: ps
args:
  - -a
```

### 11. Flexible Argument Handling

**Features:**
- Multiple arguments support
- Variable substitution in arguments
- Path normalization for file paths
- Flag preservation (doesn't normalize flags like `/?`, `-v`, etc.)

## Advanced Features

### Linea App

Create structured application directories with workflows and scripts:

```bash
linea app create my-app
```

Creates a directory structure:
```
my-app/
├─ .linea/workflows/    # Workflow YAML files (executable as commands)
├─ scripts/             # Lineash scripts (.lnsh files)
└─ README.md
```

**Benefits:**
- Organize workflows in a structured directory
- Execute workflows as commands from scripts
- Share app configurations with teams
- Version control entire app structures

### Lineash Scripts

Execute bash-like scripts that can run Linea workflows as first-class commands:

```bash
lineash scripts/deploy.lnsh [args...]
```

**Features:**
- **Friendly Syntax**: Simplified conditionals and loops with `end` keyword
- **Variables**: `VAR="value"` and `$VAR` substitution
- **Positional Parameters**: `$1`, `$2`, etc. from command-line arguments
- **Arithmetic Expressions**: `$((expression))` for calculations
- **Conditionals**: `if condition ... else ... end` with operators `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Loops**: `for VAR in list ... end` and `while condition ... end`
- **Workflow Commands**: Workflows in `.linea/workflows/` become executable commands
- **System Commands**: Unknown commands forwarded to system shell
- **No Shebang Required**: Scripts can run without `#!/bin/lineash` at the top

**Example (Friendly Syntax):**
```bash
# No shebang required!
VM_NAME="my-vm"
VM_OS="alpine"

if $VM_OS == "alpine"
    create-vm -s name=$VM_NAME
else
    echo "Using different OS"
end

for env in dev staging prod
    deploy -s environment=$env
end

counter=1
while $counter <= 3
    echo "Iteration $counter"
    counter=$((counter + 1))
end
```

**How It Works:**
1. Lineash scans `.linea/workflows/` for available workflows
2. Workflows become executable commands in scripts
3. Unknown commands are forwarded to the system shell
4. Variables, conditionals, and loops work with friendly syntax
5. Positional parameters (`$1`, `$2`, etc.) are available from command-line arguments
6. Arithmetic expressions are evaluated before variable substitution

**Backward Compatibility:**
Traditional bash syntax (`if/then/else/fi`, `for/do/done`) is still fully supported.

### Nested Variable References

Variables can reference other variables:

```yaml
variables:
  name: "Hello, $variable"
  message: "{name} World"
```

### Path Normalization

Automatic path normalization ensures:
- Windows paths use backslashes
- Unix paths use forward slashes
- Flags are not accidentally normalized

### Windows Shell Built-in Support

On Windows, Linea automatically detects shell built-ins (like `echo`, `dir`) and executes them through `cmd.exe` when needed.

### Error Handling

**Comprehensive error messages for:**
- Missing YAML files
- Invalid YAML syntax
- Undefined variables
- Command execution failures
- Missing required fields

## Use Cases

### DevOps Automation
- Define deployment commands in YAML
- Version control infrastructure commands
- Share commands across team members

### Development Workflows
- Standardize build processes
- Document complex command invocations
- Create reusable command templates

### System Administration
- Cross-platform script execution
- Consistent command execution
- Automated task management

### CI/CD Integration
- Define pipeline commands in YAML
- Parameterize commands for different environments
- Version control pipeline configurations

### 12. Multiple Commands in Single File

Define multiple commands in one YAML file separated by `---`:

```yaml
command: echo
args:
  - "First step"
---
command: echo
args:
  - "Second step"
---
command: echo
args:
  - "Third step"
```

**Benefits:**
- Execute command sequences
- Build workflows
- Share multi-step processes
- Version control entire workflows

## Future Features (Roadmap)

- Environment variable support
- Command dependencies and conditional execution
- Command templates and inheritance
- Plugin system for custom command types
- Cloud-init integration
- Terraform command support
- Command history and logging
- Interactive mode
- Parallel command execution

