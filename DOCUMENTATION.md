# Linea Documentation

Complete documentation for Linea CLI tool.

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [YAML File Format](#yaml-file-format)
4. [Command Reference](#command-reference)
5. [Variables](#variables)
6. [Cross-Platform Support](#cross-platform-support)
7. [Examples](#examples)
8. [Troubleshooting](#troubleshooting)

## Installation

### Prerequisites

- Go 1.18 or higher
- Git (for cloning the repository)

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd linea

# Build for your platform
go build -o bin/linea

# On Windows
go build -o bin/linea.exe
```

### Installing via Go

```bash
go install
```

## Quick Start

1. Create a YAML file:

```yaml
# hello.yml
command: echo
args:
  - "Hello, World!"
```

2. Run it:

```bash
linea run hello.yml
```

3. Output:
```
Hello, World!
```

## YAML File Format

### Basic Structure

```yaml
command: <command-name>
subcommand: <optional-subcommand>
args:
  - <argument1>
  - <argument2>
variables:
  <var1>: <value1>
  <var2>: <value2>
```

### Field Descriptions

#### `command` (required)
The main command to execute. Must be a valid executable or shell built-in.

**Example:**
```yaml
command: docker
```

#### `subcommand` (optional)
A subcommand for the main command.

**Example:**
```yaml
command: docker
subcommand: ps
```

#### `args` (optional)
List of arguments to pass to the command.

**Example:**
```yaml
args:
  - -a
  - -l
  - "/path/to/file"
```

#### `variables` (optional)
Key-value pairs for variable substitution.

**Example:**
```yaml
variables:
  name: "John"
  path: "/home/user"
```

## Command Reference

### `run`

Execute a command defined in a YAML file.

**Syntax:**
```bash
linea run [options] <yaml-file>
```

**Options:**
- `-v, --verbose`: Show the command before executing
- `--args <var>=<value>`: Provide variable values

**Examples:**
```bash
# Basic execution
linea run config.yml

# With verbose output
linea run -v config.yml

# With command-line variables
linea run config.yml --args name="John" --args age=30
```

### `test`

Perform a dry-run of a command without executing it.

**Syntax:**
```bash
linea test [options] <yaml-file>
```

**Options:**
- `--args <var>=<value>`: Provide variable values for testing

**Examples:**
```bash
# Dry-run a command
linea test config.yml

# Test with variables
linea test config.yml --args variable="test"
```

**Output:**
```
Dry run - would execute:
<full-command>
```

### `help`

Display information about a command defined in a YAML file.

**Syntax:**
```bash
linea help <yaml-file>
```

**Example:**
```bash
linea help config.yml
```

**Output:**
```
Command: docker
Subcommand: ps
Arguments: [-a]

Full command: docker ps -a
```

### `init`

Initialize a new workflow YAML file with template and documentation.

**Syntax:**
```bash
linea init <file-name>
```

**Examples:**
```bash
# Create a new workflow file
linea init workflow.yml

# Create in a specific directory
linea init examples/my-workflow.yml

# Create with custom name
linea init deploy-commands.yml
```

**What it creates:**
- Template YAML structure
- Documentation comments explaining each field
- Example usage instructions
- Variable examples
- Multiple commands example (commented)

**Note:** The command will fail if the file already exists to prevent overwriting.

## Variables

### Variable Syntax

Linea supports two variable syntaxes:

1. **Curly Brace Syntax:** `{variable}`
2. **Dollar Sign Syntax:** `$variable`

### Variable Sources

Variables can be defined in two ways:

1. **In YAML file:**
```yaml
variables:
  name: "John"
```

2. **Via command-line:**
```bash
linea run config.yml --args name="John"
```

**Priority:** Command-line variables override YAML variables.

### Variable Substitution

Variables are substituted in:
- Command arguments
- Variable values (nested substitution)

**Example:**
```yaml
command: echo
args:
  - "{greeting}, {name}!"
variables:
  greeting: "Hello"
  name: "World"
```

Result: `echo Hello, World!`

### Variable Validation

Linea validates that all referenced variables are defined:

```bash
$ linea run config.yml
Error: undefined variables: name (use --args to provide values)
```

## Cross-Platform Support

### Path Normalization

Linea automatically normalizes paths based on the operating system:

- **Windows:** Converts `/` to `\`
- **Unix/Linux/macOS:** Converts `\` to `/`

**Example:**
```yaml
args:
  - "C:/Users/File.txt"  # Becomes C:\Users\File.txt on Windows
```

### Flag Preservation

Flags are not normalized:
- `/?` (Windows help) remains `/?`
- `-v` remains `-v`
- `--help` remains `--help`

### Windows Shell Built-ins

On Windows, shell built-ins (like `echo`, `dir`) are automatically executed through `cmd.exe` when not found in PATH.

## Examples

### Simple Command

```yaml
# simple.yml
command: echo
args:
  - "Hello, Linea!"
```

```bash
linea run simple.yml
```

### Command with Subcommand

```yaml
# docker-ps.yml
command: docker
subcommand: ps
args:
  - -a
```

```bash
linea run docker-ps.yml
```

### Using Variables

```yaml
# greet.yml
command: echo
args:
  - "Hello, {name}!"
variables:
  name: "World"
```

```bash
linea run greet.yml
# Output: Hello, World!

# Override variable
linea run greet.yml --args name="John"
# Output: Hello, John!
```

### Complex Example

```yaml
# build.yml
command: docker
subcommand: build
args:
  - -t
  - "{image_name}:{tag}"
  - -f
  - "{dockerfile}"
  - "{context}"
variables:
  image_name: "myapp"
  tag: "latest"
  dockerfile: "./Dockerfile"
  context: "."
```

```bash
linea run build.yml --args tag="v1.0.0"
```

## Troubleshooting

### Common Issues

#### "executable file not found in %PATH%"

**Problem:** Command not found on Windows.

**Solution:** Linea automatically handles Windows shell built-ins. For other commands, ensure they're in your PATH or use full paths.

#### "undefined variables: variable"

**Problem:** Variable referenced but not defined.

**Solution:** 
- Define the variable in the YAML `variables` section, or
- Provide it via `--args` flag: `--args variable="value"`

#### "failed to parse YAML file"

**Problem:** Invalid YAML syntax.

**Solution:** 
- Check YAML syntax (indentation, quotes, etc.)
- Validate YAML with an online validator
- Ensure `command` field is present

#### Path Issues on Windows

**Problem:** Paths not working correctly.

**Solution:** Linea automatically normalizes paths. Use forward slashes in YAML; they'll be converted on Windows.

### Getting Help

- Check the [README.md](README.md) for quick reference
- Review [FEATURES.md](FEATURES.md) for feature details
- See [examples/](examples/) directory for example YAML files
- Open an issue on GitHub for bugs or feature requests

## Best Practices

1. **Use descriptive variable names:** `image_name` instead of `img`
2. **Document complex commands:** Add comments in YAML files
3. **Version control YAML files:** Track command configurations
4. **Test before running:** Use `linea test` to verify commands
5. **Use variables for paths:** Makes commands portable across systems
6. **Validate early:** Check for undefined variables before execution

## Advanced Usage

### Nested Variables

Variables can reference other variables:

```yaml
variables:
  base_path: "/home/user"
  config_path: "{base_path}/config"
```

### Multiple Variable Sources

Combine YAML and command-line variables:

```yaml
# config.yml
variables:
  env: "development"
```

```bash
linea run config.yml --args env="production"
```

The command-line value (`production`) overrides the YAML value (`development`).

