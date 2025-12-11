<div align="center">

  <h1>Linea - Commandline Workflow Tool</h1>

  <p><strong>Cross-Platform YAML-Driven Command Execution Tool</strong></p>

  <p>
    <img src="https://img.shields.io/github/v/release/marcuwynu23/linea?include_prereleases&style=flat-square" alt="Release"/>
    <img src="https://img.shields.io/github/go-mod/go-version/marcuwynu23/linea?style=flat-square" alt="Go Version"/>
    <img src="https://img.shields.io/github/stars/marcuwynu23/linea?style=flat-square" alt="GitHub Stars"/>
    <img src="https://img.shields.io/github/forks/marcuwynu23/linea?style=flat-square" alt="GitHub Forks"/>
    <img src="https://img.shields.io/github/license/marcuwynu23/linea?style=flat-square" alt="License"/>
    <img src="https://img.shields.io/github/issues/marcuwynu23/linea?style=flat-square" alt="GitHub Issues"/>
  </p>

</div>

Linea is an open-source CLI tool that allows developers and sysadmins to define commands, subcommands, arguments, and execution parameters in **YAML files**. It reads these files and executes the commands in a **cross-platform** way, handling OS-specific quirks like path separators, help flags, and other environment differences.


## Why Linea?

### The Problem with Traditional Scripting

Traditional shell scripts (Bash, PowerShell, etc.) have limitations:
- **Platform-specific**: Different syntax for Windows vs Unix
- **Hard to version control**: Complex scripts are difficult to review and maintain
- **No validation**: Errors only surface at runtime
- **Limited reusability**: Hard to parameterize and share
- **Poor documentation**: Command intent is buried in code

### The Linea Solution

Linea provides a **declarative, cross-platform approach** to command execution:

✅ **Universal YAML format** - Works identically on Windows, Linux, and macOS  
✅ **Version control friendly** - Easy to review, diff, and merge  
✅ **Built-in validation** - Catches errors before execution  
✅ **Parameterized execution** - Variables and overrides without code changes  
✅ **Self-documenting** - YAML structure makes intent clear  
✅ **Team collaboration** - Share command configurations effortlessly  

### Key Advantages Over Scripting

| Feature | Traditional Scripts | Linea |
|---------|---------------------|------------|
| **Cross-platform** | Requires separate scripts | Single YAML works everywhere |
| **Version control** | Hard to review changes | Easy to diff and merge |
| **Validation** | Runtime errors only | Pre-execution validation |
| **Documentation** | Comments in code | Self-documenting structure |
| **Reusability** | Copy-paste or functions | Template-based with variables |
| **Testing** | Manual execution | Built-in dry-run mode |
| **Maintenance** | Code complexity grows | Simple YAML structure |

## Features

- **YAML-Driven Execution**: Commands, subcommands, and arguments are defined in YAML files
- **Cross-Platform Support**: Automatically detects OS and converts paths to proper separators
- **Variable Substitution**: Use `{variable}` or `$variable` placeholders with validation
- **Command-Line Overrides**: Pass variables at runtime without modifying files
- **Dry-Run Mode**: Test commands without executing them using the `test` subcommand
- **Help Command**: Display information about commands defined in YAML files
- **Variable Validation**: Ensures all required variables are defined before execution

## Installation

```bash
go build -o bin/linea
```

On Windows:
```bash
go build -o bin/linea.exe
```

Or install directly:

```bash
go install
```

## Usage

```bash
linea <subcommand> <yaml-file>
```

### Subcommands

- `run` - Execute the command defined in the YAML file
- `test` - Dry-run the command (print without executing)
- `help` - Display information about the command
- `init` - Initialize a new workflow YAML file with template and documentation

## Example YAML Files

### Simple Echo Command

```yaml
# examples/echo-simple.yml
command: echo
args:
  - "Hello, Linea!"
```

Run it:
```bash
linea run examples/echo-simple.yml
```

### Command with Variables

```yaml
# examples/echo-variables.yml
command: echo
args:
  - "Message: {message}"
  - "User: {user}"
variables:
  message: "Welcome to Linea CLI"
  user: "Developer"
```

### Command with Command-Line Variables (--args)

You can override or provide variables at runtime using the `--args` flag and reference them with `$variable` syntax:

```yaml
# examples/greet.yml
command: echo
args:
  - "Hello, $name! Welcome to $platform."
variables:
  platform: "Linea CLI"
```

Run with command-line variables:
```bash
linea run examples/greet.yml --args name="John"
```

Output:
```
Hello, John! Welcome to Linea CLI.
```

You can also override YAML variables:
```bash
linea run examples/greet.yml --args name="John" --args platform="Commandline Workflow"
```

Output:
```
Hello, John! Welcome to Commandline Workflow.
```

### Docker Command

```yaml
# examples/docker-ps.yml
command: docker
subcommand: ps
args:
  - -a
```

### List Directory

```yaml
# examples/ls-directory.yml
command: ls
args:
  - -l
  - -a
  - "{directory}"
variables:
  directory: "."
```

## YAML File Structure

### Single Command

```yaml
command: <main-command>
subcommand: <optional-subcommand>
args:
  - <arg1>
  - <arg2>
  - <...>
variables:
  <var1>: <value1>
  <var2>: <value2>
```

### Multiple Commands

Separate multiple commands with `---`:

```yaml
command: <main-command>
args:
  - <arg1>
variables:
  <var1>: <value1>
---
command: <another-command>
args:
  - <arg2>
variables:
  <var2>: <value2>
---
# More commands...
```

## Examples

### Initialize a New Workflow

Create a new workflow file with a template:

```bash
linea init workflow.yml
```

This creates a new file with:
- Template structure
- Documentation comments
- Example usage
- Variable examples

Output:
```
✅ Created workflow file: workflow.yml

You can now:
  • Edit the file to customize your workflow
  • Test it: linea test workflow.yml
  • Run it: linea run workflow.yml
```

### Dry-Run a Command

```bash
linea test examples/docker-ps.yml
```

Output:
```
Dry run - would execute:
docker ps -a
```

### Get Help for a Command

```bash
linea help examples/echo-variables.yml
```

### Multiple Commands in One File (Advanced)

You can define multiple commands in a single YAML file by separating them with `---`:

```yaml
# examples/multi.yml
command: echo
args:
  - "First command: {message}"
variables:
  message: "Hello from first command"
---
command: echo
args:
  - "Second command: {message}"
variables:
  message: "Hello from second command"
---
command: echo
args:
  - "Third command: {message}"
variables:
  message: "Hello from third command"
```

Run all commands sequentially:
```bash
linea run -v examples/multi.yml
```

Output:
```
Found 3 commands in YAML file

[1/3] Executing: echo First command: Hello from first command
First command: Hello from first command

[2/3] Executing: echo Second command: Hello from second command
Second command: Hello from second command

[3/3] Executing: echo Third command: Hello from third command
Third command: Hello from third command
```

Test multiple commands:
```bash
linea test examples/multi.yml
```

Get help for all commands:
```bash
linea help examples/multi.yml
```

**Use Cases:**
- Execute a sequence of related commands
- Build workflows with multiple steps
- Run commands in a specific order
- Share command sequences with your team

## Running Tests

```bash
go test ./tests/...
```

## Project Structure

```
linea/
  main.go              # CLI entry point
  cmd/                 # Subcommands (run, test, help)
  internal/            # Core logic (parser, executor, utils)
  examples/            # Example YAML files
  tests/               # Test files
  bin/                 # Compiled executable (gitignored)
  go.mod               # Go module definition
  README.md            # This file
```

## Use Cases

### DevOps & Infrastructure
- **Deployment automation**: Define deployment commands in version-controlled YAML files
- **Multi-environment management**: Use variables to adapt commands for dev/staging/prod
- **Team standardization**: Share consistent command configurations across team members
- **CI/CD integration**: Parameterize pipeline commands for different environments

### Development Workflows
- **Build processes**: Standardize build commands across projects
- **Database migrations**: Version control database command sequences
- **Testing automation**: Define test execution commands declaratively
- **Code generation**: Template-based code generation with variable substitution

### System Administration
- **Cross-platform scripts**: Single YAML file works on all operating systems
- **Configuration management**: Document and version control system commands
- **Backup automation**: Define backup procedures in YAML
- **Monitoring setup**: Configure monitoring commands consistently

### Documentation & Training
- **Command libraries**: Build reusable command templates
- **Onboarding**: New team members can understand commands from YAML structure
- **Knowledge sharing**: Share complex command invocations easily
- **Documentation**: Self-documenting command configurations

## Documentation

### [FEATURES.md](FEATURES.md)
Complete feature documentation including:
- **Core Features**: YAML-driven execution, cross-platform support, variable substitution
- **Advanced Features**: Nested variables, path normalization, Windows shell support
- **Use Cases**: DevOps automation, development workflows, system administration
- **Roadmap**: Planned features and future enhancements

### [DOCUMENTATION.md](DOCUMENTATION.md)
Comprehensive user guide covering:
- **Installation**: Building from source and installation methods
- **Quick Start**: Get up and running in minutes
- **YAML Format**: Complete reference for YAML file structure
- **Command Reference**: Detailed documentation for `run`, `test`, and `help` commands
- **Variables**: Variable syntax, sources, substitution, and validation
- **Cross-Platform**: Path normalization, flag preservation, OS-specific handling
- **Examples**: Real-world examples and use cases
- **Troubleshooting**: Common issues and solutions

### [CONTRIBUTING.md](CONTRIBUTING.md)
Guidelines for contributing to Linea:
- **Code of Conduct**: Community standards and expectations
- **How to Contribute**: Reporting bugs, suggesting features, pull requests
- **Development Setup**: Prerequisites and getting started
- **Coding Guidelines**: Code style, project structure, testing requirements
- **Commit Messages**: Format and conventions
- **Review Process**: What to expect when submitting contributions

### [GUIDELINES.md](GUIDELINES.md)
Development standards and best practices:
- **Code Style**: Go conventions, naming, error handling, comments
- **Project Structure**: Directory organization and package layout
- **Testing Guidelines**: Test structure, coverage, and best practices
- **Code Review**: Checklist for pull requests
- **Git Workflow**: Branch naming and commit message conventions
- **Performance**: Optimization guidelines and considerations
- **Security**: Best practices for secure code

## Quick Comparison: Script vs Linea

### Traditional Shell Script
```bash
#!/bin/bash
# Hard to maintain, platform-specific
if [[ "$OSTYPE" == "msys" ]]; then
    docker build -t myapp:latest -f .\Dockerfile .
else
    docker build -t myapp:latest -f ./Dockerfile .
fi
```

### Linea YAML
```yaml
# Simple, cross-platform, version-controlled
command: docker
subcommand: build
args:
  - -t
  - "{image}:{tag}"
  - -f
  - "{dockerfile}"
  - "{context}"
variables:
  image: "myapp"
  tag: "latest"
  dockerfile: "./Dockerfile"
  context: "."
```

**Benefits:**
- ✅ Works on Windows, Linux, macOS without modification
- ✅ Easy to review in pull requests
- ✅ Variables can be overridden: `linea run build.yml --args tag="v1.0.0"`
- ✅ Self-documenting structure
- ✅ Built-in validation prevents errors

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

**Quick Start:**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./tests/...`)
6. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for complete contribution guidelines, code style standards, and development setup instructions.

## Support

- **Issues:** Report bugs or request features on GitHub
- **Discussions:** Ask questions and share ideas
- **Funding:** Support the project via [PayPal](https://www.paypal.com/paypalme/wynumarcu23)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors who help improve Linea
- Built with [Go](https://golang.org/)
- YAML parsing powered by [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml)

---

**Made with ❤️ by the Linea community**
