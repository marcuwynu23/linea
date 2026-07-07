<div align="center">

# User Guide

Complete reference for Linea — the cross-platform YAML-driven command execution tool.

</div>

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Templates](#templates)
- [Concepts](#concepts)
- [CI/CD Integration](#cicd-integration)
- [Workflows](#workflows)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Installation

### Prerequisites

| Tool | Version | Purpose         |
| ---- | ------- | --------------- |
| Go   | 1.18+   | Compiler        |
| Git  | Any     | Version control |

### Build from source

```bash
git clone https://github.com/marcuwynu23/linea.git
cd linea

# Build linea
go build -o bin/linea .

# Build lineash (script interpreter)
go build -o bin/lineash ./lineash
```

Add `bin/` to your `PATH` or copy the binaries to a directory in your path:

```bash
# Linux / macOS
sudo cp bin/linea bin/lineash /usr/local/bin/

# Windows (PowerShell as Administrator)
Copy-Item bin\linea.exe, bin\lineash.exe C:\Windows\System32\
```

### Install via `go install`

```bash
go install github.com/marcuwynu23/linea@latest
```

### Pre-built binaries

Download from the [Releases page](https://github.com/marcuwynu23/linea/releases). Available for:

- Linux: amd64, 386, arm64
- Windows: amd64, 386, arm64
- macOS: amd64, arm64

### Verify installation

```bash
linea --help
lineash --help
```

---

## Quick Start

### 1. Create a workflow file

```yaml
# hello.yml
command: echo
args:
  - "Hello, Linea!"
```

### 2. Preview without executing

```bash
linea test hello.yml
```

Output:

```
Dry run - would execute:
echo Hello, Linea!
```

### 3. Execute

```bash
linea run hello.yml
```

### 4. Add variables

```yaml
# greet.yml
command: echo
args:
  - "Hello, $name! Running on $platform."
variables:
  platform: "Linea CLI"
```

```bash
linea run greet.yml -s name="Developer"
```

Output:

```
Hello, Developer! Running on Linea CLI.
```

### 5. Create a Linea App

```bash
linea app create my-project
cd my-project
# Add workflows to .linea/workflows/ and scripts to scripts/
```

---

## Command Reference

### `linea run`

Execute a command defined in a YAML file.

```bash
linea run [flags] <yaml-file>
```

| Flag                      | Default | Description                             |
| ------------------------- | ------- | --------------------------------------- |
| `-v, --verbose`           | `false` | Show the full command before execution  |
| `-s, --set <var>=<value>` | —       | Set or override a variable (repeatable) |

#### Examples

**Basic execution:**

```bash
linea run config.yml
```

**Verbose mode (shows command before running):**

```bash
linea run -v config.yml
```

**Override variables:**

```bash
linea run config.yml -s name="John" -s env=production
```

**Multi-command file (executes all commands sequentially):**

```bash
linea run -v workflow.yml
```

---

### `linea test`

Dry-run a command — prints what would be executed without actually running it.

```bash
linea test [flags] <yaml-file>
```

| Flag                      | Default | Description                             |
| ------------------------- | ------- | --------------------------------------- |
| `-s, --set <var>=<value>` | —       | Provide variable values for the dry-run |

#### Examples

```bash
linea test config.yml
linea test config.yml -s name="John"
linea test workflow.yml  # tests all commands in a multi-command file
```

Output:

```
Dry run - would execute:
docker ps -a
```

---

### `linea help`

Display information about a command defined in a YAML file.

```bash
linea help <yaml-file>
```

#### Examples

```bash
linea help config.yml
```

Output:

```
Command: docker
Subcommand: ps
Arguments: [-a]

Full command: docker ps -a
```

For multi-command files, help is shown for each command sequentially.

---

### `linea init`

Scaffold a new workflow YAML file with template structure, documentation comments, and variable examples.

```bash
linea init <file-name>
```

#### Examples

```bash
linea init workflow.yml
linea init examples/my-workflow.yml
linea init deploy-commands.yml
```

This creates a YAML file with:

- Template structure (command, subcommand, args, variables fields)
- Documentation comments explaining each field
- Example usage instructions
- Variable examples with both `{}` and `$` syntax
- Multi-command example (commented out)
- Command-line override examples

**Note:** The command will fail if the file already exists to prevent accidental overwrites.

---

### `linea app create`

Create a structured Linea App directory with workflow and script folders.

```bash
linea app create <app-name>
```

Creates:

```
<app-name>/
├── .linea/workflows/   # Place workflow YAML files here
├── scripts/            # Place Lineash scripts (.lnsh) here
└── README.md
```

#### Examples

```bash
# Create a new app
linea app create my-app

# Navigate to it
cd my-app

# Add a workflow
cat > .linea/workflows/deploy.yml <<EOF
command: echo
args:
  - "Deploying $app to $env"
EOF

# Write a script
cat > scripts/deploy.lnsh <<EOF
echo "Starting deployment..."
deploy -s app=my-app -s env=production
echo "Done!"
EOF

# Run it
lineash scripts/deploy.lnsh
```

---

### `lineash`

Execute a Lineash script — a bash-like interpreter that treats workflows in `.linea/workflows/` as native commands.

```bash
lineash <script.lnsh> [args...]
```

#### Examples

```bash
lineash scripts/deploy.lnsh
lineash scripts/deploy.lnsh my-app production
```

**Features:**

- **Variables**: `VAR="value"` assignment and `$VAR` substitution
- **Positional parameters**: `$1`, `$2`, ... from command-line arguments
- **Arithmetic**: `$((expression))` for calculations
- **Conditionals**: `if condition ... else ... end` with operators `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Loops**: `for VAR in list ... end` and `while condition ... end`
- **Workflow commands**: Any workflow in `.linea/workflows/` is callable as a command
- **System commands**: Unknown commands are forwarded to the system shell
- **No shebang required**: Scripts work without `#!/bin/lineash`
- **Backward compatible**: Traditional `if/then/else/fi`, `for/do/done` syntax also supported

#### Lineash script example

```bash
# No shebang required!
APP_NAME=$1
ENV=$2

echo "Deploying $APP_NAME to $ENV..."

if $ENV == "production"
    echo "Production deployment — running checks..."
    health-check -s app=$APP_NAME
else
    echo "Non-production deployment"
end

for step in build test deploy
    echo "Running $step..."
    run-step -s step=$step
end

echo "Deployment complete!"
```

---

## Configuration

### YAML file structure

#### Single command

```yaml
command: <executable>
subcommand: <optional-subcommand>
args:
  - <argument1>
  - <argument2>
variables:
  <key1>: <value1>
  <key2>: <value2>
```

#### Multi-command (sequential)

Separate command blocks with `---`:

```yaml
command: echo
args:
  - "First command"
variables:
  message: "hello"
---
command: docker
subcommand: ps
args:
  - -a
---
command: echo
args:
  - "Done"
```

All commands execute in order. Use `linea run -v` to see progress.

### Field reference

| Field        | Required | Type                | Description                                                       |
| ------------ | -------- | ------------------- | ----------------------------------------------------------------- |
| `command`    | Yes      | `string`            | Executable to run (e.g., `echo`, `docker`, `git`)                 |
| `subcommand` | No       | `string`            | Subcommand passed to the executable (e.g., `ps`, `build`, `push`) |
| `args`       | No       | `[]string`          | List of arguments; supports variable substitution                 |
| `variables`  | No       | `map[string]string` | Key-value pairs for `{var}` and `$var` placeholders               |

### Variable syntax

| Syntax       | Overridable via `-s/--set` | Behavior                                              |
| ------------ | -------------------------- | ----------------------------------------------------- |
| `{variable}` | No                         | Protected default — always resolves to the YAML value |
| `$variable`  | Yes                        | Resolves to YAML value unless overridden at runtime   |

```yaml
command: echo
args:
  - "Protected: {name}, Overridable: $name"
variables:
  name: "default"
```

```bash
linea run config.yml -s name="custom"
# Output: Protected: default, Overridable: custom
```

### Variable sources (precedence)

1. **CLI flags** (`-s/--set name=value`) — highest priority
2. **YAML variables section** — medium priority
3. **Hardcoded in args** — lowest (no substitution needed)

### Nested variable references

Variables can reference other variables:

```yaml
variables:
  base_path: "/home/user"
  config_path: "{base_path}/config"
  message: "Hello, $name!"
```

### Configuration precedence

```
CLI -s/--set  >  YAML variables section  >  Literal text in args
```

---

## Templates

### `linea init` template structure

Running `linea init workflow.yml` creates a file with:

```yaml
# Workflow Configuration
# This file defines commands for the Linea CLI tool.

command: echo
args:
  - "Hello, {message}!"
variables:
  message: "world"
```

Plus commented-out sections showing:

- Multiple commands with `---`
- Variable override usage with `-s/--set`
- Subcommand usage

### Data model

The internal Go types that define the YAML schema:

```go
// Single command config
type CommandConfig struct {
    Command     string            `yaml:"command"`
    Subcommand  string            `yaml:"subcommand,omitempty"`
    Args        []string          `yaml:"args,omitempty"`
    Variables   map[string]string `yaml:"variables,omitempty"`
}

// Wrapper for multi-command files
type WorkflowConfig struct {
    Commands []CommandConfig
}
```

---

## Concepts

### How variable substitution works

1. Linea reads the YAML file and resolves the `variables` map
2. For each arg, it substitutes `{variable}` and `$variable` placeholders:
   - `{variable}` is replaced with the YAML-defined value (never overridden)
   - `$variable` is replaced with the YAML value, unless `-s/--set` provides an override
3. Nested variables are resolved recursively (variables referencing other variables)
4. If any referenced variable is undefined, execution stops with a clear error

### Cross-platform path normalization

Linea automatically normalizes file paths based on the detected OS:

| OS            | Input               | Normalized          |
| ------------- | ------------------- | ------------------- |
| Windows       | `C:/Users/file.txt` | `C:\Users\file.txt` |
| Linux / macOS | `C:\Users\file.txt` | `C:/Users/file.txt` |

**Flags are preserved** — `/?`, `-v`, `--help` are never normalized.

### Windows shell built-ins

On Windows, if a command (like `echo`, `dir`, `type`) is not found in `PATH`, Linea automatically executes it through `cmd.exe /c` to access shell built-ins.

### Multi-command execution

When a YAML file contains multiple commands separated by `---`:

- Commands execute sequentially in order
- `linea run -v` shows progress: `[1/3]`, `[2/3]`, etc.
- If any command fails, subsequent commands still execute (no short-circuit)
- `linea test` dry-runs all commands
- `linea help` shows details for all commands

### Linea App concept

A Linea App is a structured directory that organizes related workflows and scripts:

```
my-app/
├── .linea/workflows/   # YAML workflow files become callable commands
├── scripts/            # Lineash scripts (.lnsh files)
└── README.md           # App documentation
```

Workflows placed in `.linea/workflows/` are automatically discovered by `lineash` and can be invoked as commands inside scripts.

---

## CI/CD Integration

### GitHub Actions

**Test workflow:**

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - run: go build -o bin/linea .
      - run: go test ./tests/...
```

**Deploy workflow with Linea:**

```yaml
name: Deploy
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - run: go build -o bin/linea .
      - run: |
          ./bin/linea run deploy.yml \
            -s env=production \
            -s tag=${{ github.sha }}
```

### GitLab CI

```yaml
image: golang:1.23

stages:
  - test
  - deploy

test:
  stage: test
  script:
    - go build -o bin/linea .
    - go test ./tests/...

deploy:
  stage: deploy
  script:
    - go build -o bin/linea .
    - ./bin/linea run deploy.yml -s env=$CI_ENVIRONMENT_NAME
  only:
    - main
```

### Local development

```bash
# Same YAML file, same command, different environment
linea run deploy.yml -s env=development
linea run deploy.yml -s env=staging
linea run deploy.yml -s env=production
```

---

## Workflows

### Monorepo workflow organization

```
monorepo/
├── .linea/
│   └── workflows/
│       ├── build.yml
│       ├── test.yml
│       └── deploy.yml
├── scripts/
│   ├── ci.lnsh
│   └── dev.lnsh
├── services/
│   ├── api/
│   └── web/
└── linea.yml
```

### Using Linea App for project scaffolding

```bash
# Create a new project with standard workflow structure
linea app create my-service
cd my-service

# Define common tasks
cat > .linea/workflows/build.yml <<EOF
command: docker
subcommand: build
args:
  - -t
  - "$image:$tag"
  - "."
variables:
  image: "my-service"
  tag: "latest"
EOF

# Create a development script
cat > scripts/dev.lnsh <<EOF
echo "Starting development environment..."
build -s tag=dev
docker run -p 8080:8080 my-service:dev
EOF
```

### Using variables across environments

```yaml
# deploy.yml
command: kubectl
subcommand: apply
args:
  - -f
  - "k8s/$env/deployment.yaml"
variables:
  env: "development"
```

```bash
# Same file, different environments
linea run deploy.yml                          # deploys to development
linea run deploy.yml -s env=staging           # deploys to staging
linea run deploy.yml -s env=production        # deploys to production
```

---

## Troubleshooting

### Errors and solutions

| Problem                                      | Cause                                        | Fix                                                                                 |
| -------------------------------------------- | -------------------------------------------- | ----------------------------------------------------------------------------------- |
| `executable file not found in %PATH%`        | Command not in PATH or is a Windows built-in | Linea auto-handles built-ins; for others, use full path or add to PATH              |
| `undefined variables: <name>`                | A variable is referenced but not defined     | Add it to the YAML `variables` section or pass it with `-s/--set`                   |
| `failed to parse YAML file`                  | Invalid YAML syntax                          | Check indentation, quotes, and required fields; validate with an online YAML linter |
| Paths not working on Windows                 | Forward slashes not converted                | Linea auto-normalizes; ensure you're using forward slashes in YAML                  |
| Command fails on one OS but works on another | OS-specific command or flag                  | Use variables to adapt; check if the command exists cross-platform                  |
| `unknown subcommand`                         | Wrong Linea subcommand                       | Use `linea --help` to see valid subcommands: `run`, `test`, `help`, `init`, `app`   |
| Multi-command file stops early               | Error in one command doesn't halt execution  | Check each command independently with `linea test`                                  |

### Debugging tips

```bash
# Preview before running
linea test config.yml

# See the exact command being built
linea run -v config.yml

# Check variable resolution
linea help config.yml

# Test with specific variables
linea test config.yml -s name="test"
```

### Getting help

- `linea --help` — built-in CLI help
- `README.md` — quick-start and overview
- `USER-GUIDE.md` — this document, comprehensive reference
- `examples/` directory — ready-to-run example workflows
- [GitHub Issues](https://github.com/marcuwynu23/linea/issues) — bug reports and feature requests
- [GitHub Discussions](https://github.com/marcuwynu23/linea/discussions) — questions and community help

---

## FAQ

**Q: What operating systems does Linea support?**  
A: Windows, Linux, and macOS. YAML files are portable across all three without modification.

**Q: What's the difference between `{variable}` and `$variable`?**  
A: `{variable}` is a protected default — it always uses the value from the YAML file and cannot be overridden. `$variable` can be overridden at runtime with `-s/--set`.

**Q: Can I run multiple commands as a workflow?**  
A: Yes. Separate command blocks with `---` in a single YAML file. All commands execute sequentially.

**Q: Do I need to install anything besides Linea?**  
A: No. Linea is a single static binary with no runtime dependencies.

**Q: What's the difference between `linea` and `lineash`?**  
A: `linea` runs individual commands or workflows defined in YAML files. `lineash` runs `.lnsh` script files that support variables, conditionals, loops, and can invoke workflows as native commands.

**Q: Can I use Linea in CI/CD pipelines?**  
A: Yes. Build the binary in your pipeline and use it exactly as you would locally. See the [CI/CD Integration](#cicd-integration) section.

**Q: How do I override a variable at runtime?**  
A: Use the `-s` or `--set` flag: `linea run config.yml -s name="John"`. Only `$variable` syntax can be overridden.

**Q: Does Linea work with PowerShell?**  
A: Yes. On Windows, Linea works with both Command Prompt and PowerShell. Paths and shell built-ins are handled automatically.

**Q: How do I create reusable command templates?**  
A: Use the `init` subcommand to scaffold new files: `linea init workflow.yml`. Use variables and YAML comments to create self-documenting templates.

**Q: Can I version-control my workflows?**  
A: Yes. YAML files are plain text and work great with Git — easy to diff, merge, and review.

**Q: What happens if I run a multi-command file and one command fails?**  
A: Subsequent commands still execute. There is no short-circuit behavior. Check each command independently with `linea test`.

**Q: How do I run a specific command from a multi-command file?**  
A: Currently Linea executes all commands in a multi-command file. For selective execution, use separate YAML files.

**Q: Is there a way to see all available variables for a workflow?**  
A: Yes. Use `linea help <file>` to see the full command with all variable values resolved.

**Q: Can Linea handle sensitive data like passwords?**  
A: Linea is a command execution tool — it passes values as-is to the shell. For secrets, use your CI/CD system's secret management or environment variables.

**Q: How do I contribute to Linea?**  
A: See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on setting up the dev environment, coding standards, and submitting pull requests.

---

<p align="center">
  <a href="README.md">README</a> •
  <a href="CONTRIBUTING.md">Contributing</a> •
  <a href="https://github.com/marcuwynu23/linea/issues">Issues</a>
</p>
