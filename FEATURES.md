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

Two syntaxes are supported for variable substitution:

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
- Provided via command-line `--args` flag
- Command-line variables override YAML variables

### 4. Command-Line Variable Override

Pass variables at runtime without modifying YAML files:

```bash
linea run config.yml --args variable="value"
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
- Suggests using `--args` to provide missing values
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

