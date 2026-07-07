# Contributing to Linea

Thank you for your interest in contributing to Linea! This guide covers everything you need to know to contribute effectively.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Conventions](#commit-conventions)
- [PR Process](#pr-process)
- [Release Process](#release-process)
- [Questions](#questions)

---

## Code of Conduct

This project is governed by the [Contributor Covenant](CODE_OF_CONDUCT.md). By participating, you agree to uphold its standards. Report unacceptable behavior to the project maintainers.

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.18+ | Compiler and toolchain |
| Git | Any | Version control |
| golangci-lint | Latest | Linting (optional, CI runs it) |

## Project Structure

```
linea/
├── main.go              # CLI entry point — dispatches subcommands
├── cmd/                 # Subcommand implementations
│   ├── run.go           #   `linea run`
│   ├── test.go          #   `linea test`
│   ├── help.go          #   `linea help`
│   ├── init.go          #   `linea init`
│   ├── app.go           #   `linea app`
│   └── lineash.go       #   lineash entry handling
├── internal/            # Core logic (not exported outside module)
│   ├── parser.go        #   YAML deserialization
│   ├── executor.go      #   Command execution + variable substitution
│   ├── lineash.go       #   Lineash script interpreter
│   ├── types.go         #   Shared type definitions
│   └── utils.go         #   Path normalization, OS detection, helpers
├── lineash/             # Lineash binary entry point
│   └── main.go
├── examples/            # Example YAML workflow files
├── tests/               # Test suite (flat package)
├── docs/                # Cloudflare Pages documentation site
├── .github/workflows/   # CI/CD pipeline definitions
└── bin/                 # Build output (gitignored)
```

## Development Workflow

### 1. Fork and clone

```bash
git clone https://github.com/<your-username>/linea.git
cd linea
```

### 2. Create a feature branch

```bash
git checkout -b feat/my-feature
```

Use branch prefixes:
- `feat/` — new features
- `fix/` — bug fixes
- `docs/` — documentation
- `refactor/` — code refactoring
- `test/` — test additions or changes

### 3. Make changes

Write code following the [coding standards](#coding-standards) below.

### 4. Test

```bash
go test ./tests/...
```

### 5. Commit

Follow the [commit conventions](#commit-conventions).

### 6. Push and open a PR

```bash
git push origin feat/my-feature
```

Then open a pull request on GitHub. See [PR Process](#pr-process).

## Coding Standards

### Go style

- Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format code with `gofmt` (or `go fmt ./...`)
- Run `golangci-lint run` before committing (CI enforces this)

### Naming

| Construct | Convention | Example |
|---|---|---|
| Packages | lowercase, single word | `parser`, `executor` |
| Exported functions | PascalCase | `ParseYAML` |
| Unexported functions | camelCase | `resolveVars` |
| Variables | camelCase | `configPath` |
| Constants | PascalCase | `DefaultTimeout` |
| Types | PascalCase | `CommandConfig` |

### Error handling

```go
// Always handle errors explicitly
config, err := ParseYAML(filePath)
if err != nil {
    return fmt.Errorf("parse %s: %w", filePath, err)
}

// Provide context in error messages
// Wrap with %w to enable errors.Is / errors.As
```

### Comments

- Document all exported identifiers with complete sentences
- Start comments with the identifier name:
  ```go
  // ParseYAML reads and parses a YAML file into a CommandConfig.
  func ParseYAML(filePath string) (*CommandConfig, error) {
  ```

## Testing

### Expectations

- Maintain >80% code coverage
- Test both success and failure paths
- Use table-driven tests for multiple cases
- Test cross-platform behavior where relevant (path normalization, OS detection)

### Running tests

```bash
# All tests
go test ./tests/...

# With coverage
go test -cover ./tests/...

# Specific test
go test -run TestSubstituteVariables ./tests/...

# Short mode (skips slow tests)
go test -short ./tests/...
```

### Test pattern

```go
func TestSubstituteVariables(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        vars     map[string]string
        expected string
    }{
        {
            name:     "curly brace syntax",
            input:    "{name}",
            vars:     map[string]string{"name": "John"},
            expected: "John",
        },
        {
            name:     "dollar syntax",
            input:    "$name",
            vars:     map[string]string{"name": "John"},
            expected: "John",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := SubstituteVariables(tt.input, tt.vars)
            if got != tt.expected {
                t.Errorf("got %q, want %q", got, tt.expected)
            }
        })
    }
}
```

## Commit Conventions

Linea uses [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages. This aligns with how the tool itself structures workflows.

### Format

```
<type>(<scope>): <short description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Usage | Example |
|---|---|---|
| `feat` | New feature | `feat(parser): add support for YAML anchors` |
| `fix` | Bug fix | `fix(executor): handle empty args slice` |
| `docs` | Documentation | `docs: add CI/CD integration guide` |
| `refactor` | Code change with no behavior change | `refactor(parser): extract validate function` |
| `test` | Test additions or changes | `test: add edge cases for path normalization` |
| `style` | Formatting, linting | `style: gofmt all files` |
| `chore` | Build, CI, dependencies | `chore: bump yaml.v3 to v3.0.1` |

### Scope

Scope should be the package or area affected — one of: `parser`, `executor`, `lineash`, `cmd`, `test`, `docs`, `ci`.

### Examples

```
feat(parser): add validation for duplicate variables
fix(executor): normalize paths on Windows for args with spaces
docs: add troubleshooting section for PATH issues
test: add table-driven tests for SubstituteVariables
refactor(executor): extract command builder into separate function
```

### Breaking changes

Add `BREAKING CHANGE` in the footer:

```
feat(api): change variable syntax from %var% to $var

BREAKING CHANGE: %var% syntax is no longer supported. All existing
workflows must be updated to use $var or {var} syntax.
```

## PR Process

### Before submitting

- [ ] Code follows Go style guidelines (gofmt, golangci-lint)
- [ ] All tests pass (`go test ./tests/...`)
- [ ] New features include tests
- [ ] Documentation is updated (README, USER-GUIDE, or inline)
- [ ] No breaking changes without clear justification
- [ ] Error handling is thorough
- [ ] Cross-platform compatibility is considered (Windows, Linux, macOS)

### Review criteria

- **Correctness** — Does the change work as intended?
- **Maintainability** — Is the code easy to understand and modify?
- **Test coverage** — Are new behaviors tested?
- **Performance** — No unnecessary allocations or O(n²) patterns
- **Backward compatibility** — Existing workflows should continue working

### Merge

- PRs require at least one maintainer approval
- All CI checks must pass
- Squash merge is preferred to keep history clean

## Release Process

1. Ensure all tests pass on `main`
2. Create a version tag: `git tag v1.2.3`
3. Push the tag: `git push origin v1.2.3`
4. CI builds binaries for all platforms and creates a GitHub Release
5. Release notes are auto-generated from commit history

Versioning follows [SemVer](https://semver.org/):
- **MAJOR** — breaking changes
- **MINOR** — new features (backward compatible)
- **PATCH** — bug fixes (backward compatible)

## Questions

- **Issues**: [github.com/marcuwynu23/linea/issues](https://github.com/marcuwynu23/linea/issues)
- **Discussions**: [github.com/marcuwynu23/linea/discussions](https://github.com/marcuwynu23/linea/discussions)
- **Existing docs**: [README.md](README.md) | [USER-GUIDE.md](USER-GUIDE.md)

---

*Thank you for contributing to Linea!*
