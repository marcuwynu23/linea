# Development Guidelines

This document outlines coding standards, best practices, and development guidelines for the Linea project.

## Code Style

### Go Conventions

- Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for code formatting
- Run `golint` or `golangci-lint` before committing
- Keep line length reasonable (aim for < 100 characters)

### Naming Conventions

- **Packages:** lowercase, single word
- **Functions:** PascalCase for exported, camelCase for unexported
- **Variables:** camelCase
- **Constants:** PascalCase or UPPER_CASE for exported
- **Types:** PascalCase

**Examples:**
```go
// Good
func ParseYAML(filePath string) (*Config, error)
var configPath string
const DefaultTimeout = 30

// Avoid
func parse_yaml(file_path string) // Not Go style
```

### Error Handling

- Always handle errors explicitly
- Return errors, don't ignore them
- Provide context in error messages
- Use `fmt.Errorf` with `%w` for error wrapping

**Example:**
```go
config, err := ParseYAML(filePath)
if err != nil {
    return nil, fmt.Errorf("failed to parse YAML file: %w", err)
}
```

### Comments

- Export all public functions, types, and variables
- Use complete sentences in comments
- Start with the name of the thing being described

**Example:**
```go
// ParseYAML reads and parses a YAML file into a CommandConfig.
func ParseYAML(filePath string) (*CommandConfig, error) {
    // Implementation
}
```

## Project Structure

### Directory Organization

```
linea/
  main.go              # Entry point, minimal logic
  cmd/                 # CLI subcommands
    run.go
    test.go
    help.go
  internal/            # Internal packages (not exported)
    parser.go          # YAML parsing
    executor.go        # Command execution
    utils.go           # Utilities
    types.go           # Type definitions
  examples/            # Example YAML files
  tests/               # Test files
  bin/                 # Build output
  .github/             # GitHub configuration
```

### Package Organization

- **`cmd/`**: CLI command implementations
- **`internal/`**: Core logic, not exported outside package
- Keep packages focused and cohesive
- Avoid circular dependencies

## Testing Guidelines

### Test Structure

- Test files: `*_test.go`
- Test functions: `TestFunctionName`
- Use table-driven tests when appropriate
- Test both success and failure cases

**Example:**
```go
func TestSubstituteVariables(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        vars     map[string]string
        expected string
    }{
        {
            name:     "simple substitution",
            input:    "{name}",
            vars:     map[string]string{"name": "John"},
            expected: "John",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := SubstituteVariables(tt.input, tt.vars)
            if result != tt.expected {
                t.Errorf("Expected %q, got %q", tt.expected, result)
            }
        })
    }
}
```

### Test Coverage

- Aim for > 80% code coverage
- Test edge cases and error conditions
- Test cross-platform behavior when relevant

### Running Tests

```bash
# Run all tests
go test ./tests/...

# Run with coverage
go test -cover ./tests/...

# Run specific test
go test -run TestName ./tests/...
```

## Code Review Checklist

Before submitting a PR, ensure:

- [ ] Code follows Go style guidelines
- [ ] All tests pass
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] No breaking changes (or documented)
- [ ] Error handling is proper
- [ ] Code is commented appropriately
- [ ] No hardcoded values (use constants)
- [ ] Cross-platform compatibility considered

## Git Workflow

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation
- `refactor/description` - Refactoring

### Commit Messages

Follow the format:
```
<type>: <short description>

<optional longer description>
```

**Types:**
- `Add:` - New feature
- `Fix:` - Bug fix
- `Update:` - Update existing
- `Refactor:` - Code refactoring
- `Docs:` - Documentation
- `Test:` - Tests
- `Style:` - Code style
- `Chore:` - Maintenance

## Error Messages

### Guidelines

- Be clear and actionable
- Include context
- Suggest solutions when possible
- Use consistent formatting

**Example:**
```go
return fmt.Errorf("undefined variables: %s (use -s or --set to provide values)", 
    strings.Join(missing, ", "))
```

## Documentation

### Code Documentation

- Document all exported functions
- Explain complex algorithms
- Include usage examples in comments
- Keep documentation up-to-date

### User Documentation

- Update README.md for user-facing changes
- Add examples to DOCUMENTATION.md
- Update FEATURES.md for new features
- Keep examples working

## Performance

### Considerations

- Avoid premature optimization
- Profile before optimizing
- Consider memory allocations
- Use appropriate data structures

### Best Practices

- Reuse buffers when possible
- Avoid unnecessary allocations
- Use `strings.Builder` for string concatenation
- Consider using sync.Pool for frequently allocated objects

## Security

### Guidelines

- Never execute user input directly
- Validate all inputs
- Sanitize file paths
- Handle errors securely (don't leak sensitive info)

## Cross-Platform Support

### Testing

- Test on Windows, Linux, and macOS when possible
- Use `runtime.GOOS` for OS-specific code
- Test path handling on different platforms
- Verify shell command execution

### Best Practices

- Use `filepath` package for paths
- Use `os.PathSeparator` for separators
- Handle Windows vs Unix differences
- Test shell built-ins on Windows

## Dependencies

### Guidelines

- Minimize external dependencies
- Prefer standard library
- Document why external deps are needed
- Keep dependencies up-to-date

### Current Dependencies

- `gopkg.in/yaml.v3` - YAML parsing

## Release Process

1. Update version numbers
2. Update CHANGELOG.md (if exists)
3. Run all tests
4. Build for all platforms
5. Create release tag
6. Publish release notes

## Questions?

- Check existing code for patterns
- Review similar implementations
- Ask in issues or discussions
- Follow existing conventions

## Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Testing](https://go.dev/doc/tutorial/testing)
- [Go Best Practices](https://go.dev/doc/effective_go)

