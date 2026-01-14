# Contributing to cp-config

Thank you for your interest in contributing to cp-config! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Getting Started

```bash
# Clone the repository
git clone https://github.com/tr4d3r/cp-config.git
cd cp-config

# Install dependencies
go mod download

# Build the project
go build ./...

# Run tests
go test ./...
```

### Running Locally

```bash
# Run directly
go run ./cmd/cp-config --help

# Or build and run
go build -o cp-config ./cmd/cp-config
./cp-config --help
```

## Project Structure

```
cp-config/
├── cmd/cp-config/           # CLI entry point
│   ├── main.go
│   └── cmd/                 # Cobra commands
│       ├── root.go
│       ├── list.go
│       ├── install.go
│       ├── diff.go
│       └── profile.go
├── internal/                # Internal packages
│   ├── config/             # Configuration types
│   ├── parser/             # File parsing
│   ├── github/             # Target directory operations
│   └── repo/               # Repository utilities
├── agents/                  # Sample agent configs
├── skills/                  # Sample skill configs
└── profiles/                # Sample profiles
```

## Code Style

### General Guidelines

- Follow standard Go conventions and idioms
- Use `gofmt` to format code
- Keep functions focused and small
- Write descriptive variable and function names
- Add comments for exported functions and complex logic

### Formatting

```bash
# Format all code
go fmt ./...

# Run linter (if using golangci-lint)
golangci-lint run
```

### Error Handling

- Always handle errors explicitly
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Return errors to callers; let them decide how to handle

### Testing

- Write tests for new functionality
- Place tests in `*_test.go` files alongside the code
- Use table-driven tests where appropriate
- Aim for meaningful coverage, not 100%

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...
```

## Making Changes

### Branch Naming

- `feature/*` - New features
- `fix/*` - Bug fixes
- `docs/*` - Documentation updates
- `refactor/*` - Code refactoring

### Commit Messages

Write clear, concise commit messages:

```
Short summary (50 chars or less)

More detailed explanation if needed. Wrap at 72 characters.
Explain the problem this commit solves and why this approach
was chosen.

- Bullet points are fine
- Use present tense ("Add feature" not "Added feature")
```

### Pull Request Process

1. **Fork and branch**: Create a feature branch from `develop`
2. **Make changes**: Implement your changes with tests
3. **Test**: Ensure all tests pass with `go test ./...`
4. **Format**: Run `go fmt ./...`
5. **Commit**: Write clear commit messages
6. **Push**: Push to your fork
7. **PR**: Open a pull request to `develop`

### PR Checklist

- [ ] Tests pass locally
- [ ] Code is formatted with `go fmt`
- [ ] New functionality has tests
- [ ] Documentation updated if needed
- [ ] Commit messages are clear

## Adding New Commands

1. Create a new file in `cmd/cp-config/cmd/`
2. Define the command using Cobra:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Short description",
    Long:  `Longer description with examples.`,
    Run:   runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)
}

func runMyCommand(cmd *cobra.Command, args []string) {
    // Implementation
}
```

3. Add tests if the command has complex logic

## Adding Configuration Types

1. Add struct definitions in `internal/config/`
2. Add parsing logic in `internal/parser/`
3. Update relevant commands to use the new types
4. Add tests for parsing and validation

## Reporting Issues

When reporting issues, please include:

- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant error messages

## Questions?

Feel free to open an issue for questions or discussions about contributing.
