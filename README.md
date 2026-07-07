<div align="center">

# Linea

<p>
  <img src="https://img.shields.io/github/v/release/marcuwynu23/linea?include_prereleases&style=flat-square" alt="Release"/>
  <img src="https://img.shields.io/github/go-mod/go-version/marcuwynu23/linea?style=flat-square" alt="Go Version"/>
  <img src="https://img.shields.io/github/stars/marcuwynu23/linea?style=flat-square" alt="GitHub Stars"/>
  <img src="https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square" alt="License"/>
  <img src="https://img.shields.io/github/actions/workflow/status/marcuwynu23/linea/.github/workflows/test.yml?branch=main&style=flat-square" alt="CI"/>
  <img src="https://codecov.io/gh/marcuwynu23/linea/branch/main/graph/badge.svg?style=flat-square" alt="Codecov"/>
</p>

**Cross-Platform YAML-Driven Command Execution Tool.**  
Define, share, and run complex command-line workflows from simple YAML files.

➡️ **[Read the full user guide →](USER-GUIDE.md)**

</div>

## Table of Contents

- [What Is Linea?](#what-is-linea)
- [Use Cases](#use-cases)
- [Benefits for Developers](#benefits-for-developers)
- [Advantages Over Other Tools](#advantages-over-other-tools)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Reference](#cli-reference)
- [Configuration](#configuration)
- [Example Output](#example-output)
- [CI/CD Integration](#cicd-integration)
- [Development](#development)
- [Architecture](#architecture)
- [Contributing](CONTRIBUTING.md)
- [User Guide](USER-GUIDE.md)
- [License](#license)

---

## What Is Linea?

**Linea** is a cross-platform CLI tool that organizes commands into workflows using YAML. It makes defining, sharing, and managing complex command-line tasks simple and consistent across Windows, Linux, and macOS.

### What It Does

- **Define** — Declare commands, subcommands, and arguments in portable YAML files
- **Execute** — Run workflows directly with `linea run <file>`
- **Test** — Dry-run commands with `linea test` to preview before execution
- **Parameterize** — Use `{variable}` and `$variable` placeholders with built-in validation
- **Override** — Pass variables at runtime via `-s/--set` without modifying files
- **Structure** — Create Linea Apps with `linea app create` for organized workflow directories
- **Script** — Write bash-like `.lnsh` scripts that invoke workflows as first-class commands

### Why Use It?

| Problem | How Linea Solves It |
|---|---|
| Shell scripts break across OSes | **Single YAML works on Windows, Linux, and macOS** — paths normalize automatically |
| Scripts are hard to review in PRs | **YAML diffs cleanly** — intent is obvious |
| Errors surface only at runtime | **Pre-execution validation** catches missing variables and bad syntax |
| Commands are duplicated across projects | **Template-based workflows** with variables enable reuse |
| Onboarding new devs is slow | **Self-documenting YAML** makes command intent clear without comments |
| CI pipelines need the same commands as local dev | **Same YAML file works everywhere** — local, CI, across teams |

### The Philosophy

1. **Declarative over imperative.** Say *what* to run, not *how* to run it.
2. **Your process stays yours.** No lock-in — your YAML files are plain text under version control.
3. **Portable by default.** One definition runs identically on every major OS.

## Use Cases

| Scenario | How Linea Helps |
|---|---|
| **DevOps & Deployment** | Define deployment commands in version-controlled YAML; parameterize per environment |
| **Development Workflows** | Standardize build, test, and code-gen commands across projects |
| **System Administration** | Cross-platform scripts that work on Windows, Linux, and macOS without modification |
| **CI/CD Pipelines** | Reuse the same YAML workflows locally and in GitHub Actions, GitLab CI, etc. |
| **Team Onboarding** | New members see available workflows and their intent at a glance |
| **Documentation & Training** | Self-documenting commands replace tribal knowledge and README drift |

## Benefits for Developers

- **Write once, run anywhere** — YAML files are fully cross-platform
- **Pre-flight validation** — `linea test` catches errors before they happen
- **Zero-config variables** — `{var}` for protected defaults, `$var` for runtime overrides
- **Structured apps** — `linea app create` scaffolds workflow directories with scripts
- **Lineash scripting** — Bash-like interpreter that treats workflows as native commands
- **Git-friendly** — YAML diffs are readable; merges rarely conflict
- **Self-documenting** — The YAML structure IS the documentation
- **No runtime deps** — Single static binary; no Python, Node, or Java required
- **Dry-run mode** — Preview every command without executing it
- **Verbose mode** — See the full constructed command on execution

## Advantages Over Other Tools

| Aspect | Linea | Shell Scripts | Ansible | Makefile | Handwritten |
|---|---|---|---|---|---|
| **Setup time** | ~10 seconds | Immediate | Hours | Minutes | None |
| **Cross-platform** | Native | Fragile per-OS scripts | Via SSH/agents | Limited | N/A |
| **Validation** | Built-in pre-exec | Runtime only | Playbook syntax check | None | None |
| **Diff/Review** | Clean YAML diffs | Opaque diffs | YAML but complex | Tab-heavy | N/A |
| **Variables** | `{}` / `$` syntax, CLI override | Env vars / args | Vars / extra-vars | `$(...)` | N/A |
| **Dry-run** | `linea test` | Manual echo | `--check` | `-n` | N/A |
| **CLI complexity** | Minimal (5 subcommands) | Full language | ~50 modules | Targets/rules | N/A |
| **Binary size** | ~10 MB | None | ~500 MB | None | N/A |
| **License** | Apache 2.0 | Varies | GPL-3.0 | GPL-3.0 | N/A |
| **Learning curve** | Minutes | Moderate | Steep | Moderate | N/A |
| **Scripting** | Lineash (.lnsh) | Native shell | Playbooks | Recipes | N/A |
| **CI integration** | Drop-in YAML | Script per CI | AWX/Tower | Pre-installed | Manual |

## Installation

### From source (Go 1.18+)

```bash
git clone https://github.com/marcuwynu23/linea.git
cd linea
go build -o bin/linea .

# Optional: build lineash script interpreter
go build -o bin/lineash ./lineash
```

### On Windows

```powershell
go build -o bin\linea.exe .
go build -o bin\lineash.exe .\lineash
```

### Via `go install`

```bash
go install github.com/marcuwynu23/linea@latest
```

### Binary downloads

Pre-built binaries for Linux (amd64, 386, arm64), Windows (amd64, 386, arm64), and macOS (amd64, arm64) are available on the [Releases page](https://github.com/marcuwynu23/linea/releases).

### Verify

```bash
linea --help
lineash --help
```

## Quick Start

```bash
# Create a simple workflow
cat > hello.yml <<EOF
command: echo
args:
  - "Hello, Linea!"
EOF

# Dry-run (preview without executing)
linea test hello.yml

# Execute
linea run hello.yml

# Use variables
cat > greet.yml <<EOF
command: echo
args:
  - "Hello, $name! Welcome to $platform."
variables:
  platform: "Linea CLI"
EOF

linea run greet.yml -s name="Developer"
```

## CLI Reference

### `linea run <file>`

Execute the command defined in a YAML file.

| Flag | Default | Description |
|---|---|---|
| `-v, --verbose` | `false` | Show the full command before executing |
| `-s, --set <var>=<value>` | — | Provide or override variable values (repeatable) |

```bash
linea run config.yml
linea run -v config.yml
linea run config.yml -s name="John" -s env=production
```

### `linea test <file>`

Dry-run a command — prints what would be executed without running it.

| Flag | Default | Description |
|---|---|---|
| `-s, --set <var>=<value>` | — | Provide variable values for the dry-run |

```bash
linea test config.yml
```

### `linea help <file>`

Display the full constructed command and all variable values for a YAML file.

```bash
linea help config.yml
```

### `linea init <file>`

Scaffold a new workflow YAML file with template structure, documentation comments, and variable examples.

```bash
linea init workflow.yml
```

### `linea app create <name>`

Create a structured Linea App directory with workflow and script folders.

| Flag | Default | Description |
|---|---|---|
| `-t, --template` | `default` | App template to use |

```bash
linea app create my-app
```

Creates:
```
my-app/
├── .linea/workflows/   # Workflow YAML files
├── scripts/            # Lineash scripts (.lnsh)
└── README.md
```

### `lineash <script> [args...]`

Execute a Lineash script — a bash-like interpreter that treats `.linea/workflows/` workflows as native commands.

```bash
lineash scripts/deploy.lnsh my-app production
```

## Configuration

### Single-command YAML

```yaml
command: <executable>
subcommand: <optional-subcommand>
args:
  - <argument>
variables:
  <key>: <value>
```

| Field | Required | Description |
|---|---|---|
| `command` | Yes | The executable to run |
| `subcommand` | No | Subcommand passed to the executable |
| `args` | No | List of arguments (supports variable substitution) |
| `variables` | No | Key-value pairs for `{var}` / `$var` substitution |

### Multi-command YAML

Separate commands with `---`:

```yaml
command: echo
args:
  - "Step 1"
---
command: echo
args:
  - "Step 2"
```

### Variable syntax

| Syntax | Overridable via `-s/--set` | Use case |
|---|---|---|
| `{variable}` | No | Protected defaults — always uses YAML value |
| `$variable` | Yes | Runtime overrides from CLI or CI |

### Configuration precedence

```
CLI flag (-s/--set) > YAML variables section > defaults (none)
```

## Example Output

```bash
$ linea test examples/docker-ps.yml
Dry run - would execute:
docker ps -a

$ linea run examples/greet.yml -s name="John"
Hello, John! Welcome to Linea CLI.

$ linea run -v examples/multi.yml
Found 3 commands in YAML file

[1/3] Executing: echo First command: Hello from first command
First command: Hello from first command

[2/3] Executing: echo Second command: Hello from second command
Second command: Hello from second command

[3/3] Executing: echo Third command: Hello from third command
Third command: Hello from third command
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Deploy
on: [push]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - run: go build -o bin/linea .
      - run: ./bin/linea run deploy.yml -s env=production
```

### GitLab CI

```yaml
deploy:
  image: golang:1.23
  script:
    - go build -o bin/linea .
    - ./bin/linea run deploy.yml -s env=production
```

## Development

### Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.18+ | Compiler |
| Git | Any | Version control |

### Commands

```bash
go build -o bin/linea .          # Build linea binary
go build -o bin/lineash ./lineash # Build lineash binary
go test ./tests/...              # Run all tests
go test -cover ./tests/...       # Run with coverage
```

### Project structure

```
linea/
├── main.go              # CLI entry point
├── cmd/                 # Subcommand implementations
│   ├── run.go
│   ├── test.go
│   ├── help.go
│   ├── init.go
│   ├── app.go
│   └── lineash.go
├── internal/            # Core logic (not exported)
│   ├── parser.go        # YAML parsing
│   ├── executor.go      # Command execution
│   ├── lineash.go       # Script interpreter
│   ├── types.go         # Type definitions
│   └── utils.go         # Utilities (path normalization, etc.)
├── lineash/             # Lineash script interpreter entry point
│   └── main.go
├── examples/            # Example YAML workflows
├── tests/               # Test suite
├── docs/                # Cloudflare Pages documentation site
├── .github/workflows/   # CI/CD workflows
└── bin/                 # Build output (gitignored)
```

## Architecture

- **`main.go`** dispatches subcommands (`run`, `test`, `help`, `init`, `app`) to handlers in `cmd/`
- **`cmd/`** packages parse CLI flags, read YAML files, and delegate to `internal/`
- **`internal/parser`** deserializes YAML into typed configs using `gopkg.in/yaml.v3`
- **`internal/executor`** resolves variables (`{var}` / `$var`), normalizes paths per OS, and invokes the OS shell
- **`internal/lineash`** implements a bash-like interpreter that resolves workflow names from `.linea/workflows/`
- **Path normalization** converts `/` to `\` on Windows (and vice versa) while preserving flags like `/?`

## License

Apache 2.0 — see [LICENSE](LICENSE).

---

<p align="center">
  <a href="CONTRIBUTING.md">Contributing Guidelines</a> •
  <a href="USER-GUIDE.md">User Guide</a> •
  <a href="https://github.com/marcuwynu23/linea/issues">Issues</a> •
  <a href="https://github.com/marcuwynu23/linea/discussions">Discussions</a>
</p>
