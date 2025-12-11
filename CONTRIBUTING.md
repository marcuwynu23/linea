# Contributing to Linea

Thank you for your interest in contributing to Linea! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different viewpoints and experiences

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:

1. **Clear title** describing the bug
2. **Description** of the issue
3. **Steps to reproduce** the bug
4. **Expected behavior** vs **actual behavior**
5. **Environment details:**
   - OS and version
   - Go version
   - Linea version
6. **Screenshots/logs** if applicable

### Suggesting Features

Feature suggestions are welcome! Please include:

1. **Clear description** of the feature
2. **Use case** - why is this feature needed?
3. **Proposed implementation** (if you have ideas)
4. **Examples** of how it would be used

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**
4. **Write/update tests** for your changes
5. **Ensure all tests pass:**
   ```bash
   go test ./tests/...
   ```
6. **Update documentation** if needed
7. **Commit your changes:**
   ```bash
   git commit -m "Add: description of your changes"
   ```
8. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```
9. **Open a Pull Request** with a clear description

## Development Setup

### Prerequisites

- Go 1.18 or higher
- Git
- A code editor (VS Code, GoLand, etc.)

### Getting Started

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd linea
   ```

2. **Build the project:**
   ```bash
   go build -o bin/linea
   ```

3. **Run tests:**
   ```bash
   go test ./tests/...
   ```

4. **Test your changes:**
   ```bash
   ./bin/linea test examples/echo-simple.yml
   ```

## Coding Guidelines

### Code Style

- Follow Go conventions and style guide
- Use `gofmt` to format code
- Keep functions focused and small
- Write clear, self-documenting code
- Add comments for complex logic

### Project Structure

```
linea/
  main.go              # CLI entry point
  cmd/                 # Subcommands
  internal/            # Core logic (not exported)
  examples/            # Example YAML files
  tests/               # Test files
  bin/                 # Compiled executables
```

### Testing

- Write tests for new features
- Ensure existing tests still pass
- Test on multiple platforms if possible
- Test edge cases and error conditions

**Test Structure:**
- Unit tests in `tests/` directory
- Test files: `*_test.go`
- Use descriptive test names

**Example:**
```go
func TestSubstituteVariables(t *testing.T) {
    // Test implementation
}
```

### Commit Messages

Use clear, descriptive commit messages:

**Format:**
```
<type>: <description>

[optional body]
```

**Types:**
- `Add:` - New feature
- `Fix:` - Bug fix
- `Update:` - Update existing feature
- `Refactor:` - Code refactoring
- `Docs:` - Documentation changes
- `Test:` - Test additions/changes

**Examples:**
```
Add: support for environment variables
Fix: path normalization on Windows
Update: improve error messages
Docs: add examples for nested variables
```

## Areas for Contribution

### High Priority

- Additional test coverage
- Documentation improvements
- Bug fixes
- Performance optimizations

### Feature Ideas

- Environment variable support
- Multi-command execution
- Command templates
- Plugin system
- Interactive mode
- Logging and history

### Documentation

- More examples
- Tutorial guides
- Video tutorials
- Blog posts
- Translations

## Review Process

1. **Automated checks** must pass (tests, linting)
2. **Code review** by maintainers
3. **Discussion** of any requested changes
4. **Approval** and merge

## Questions?

- Open an issue for questions
- Check existing issues and discussions
- Review documentation files

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md (if created)
- Credited in release notes
- Appreciated by the community!

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to Linea! 🚀


