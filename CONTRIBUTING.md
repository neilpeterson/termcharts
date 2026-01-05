# Contributing to termcharts

Thank you for your interest in contributing to termcharts! This document provides guidelines and information for contributors.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/termcharts.git
   cd termcharts
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/neilpeterson/termcharts.git
   ```
4. **Install dependencies**:
   ```bash
   go mod download
   ```

## Development Workflow

### Creating a Branch

Always create a feature branch for your work:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-heatmap-chart`
- `fix/sparkline-unicode-rendering`
- `docs/update-api-reference`

### Making Changes

1. **Write code** following the existing style and conventions
2. **Add tests** for new functionality
3. **Run tests** to ensure nothing is broken:
   ```bash
   make test
   ```
4. **Run the linter**:
   ```bash
   make lint
   ```
5. **Build the project**:
   ```bash
   make build
   ```

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Use meaningful variable and function names
- Add godoc comments on all exported types and functions
- Keep functions focused and reasonably sized
- Use the functional options pattern for configuration

### Testing

- Write table-driven tests where appropriate
- Aim for high test coverage on new code
- Test edge cases (empty data, invalid values, etc.)
- Run the full test suite before submitting:
  ```bash
  go test -race ./...
  ```

### Commit Messages

Write clear, concise commit messages:

```
Add horizontal grouped bar chart support

- Implement renderHorizontalGrouped function
- Add WithBarMode option for grouped/stacked selection
- Update CLI with --grouped flag
```

- Use the imperative mood ("Add feature" not "Added feature")
- Keep the first line under 50 characters
- Add a blank line before the body if needed
- Explain what and why, not how

## Submitting Changes

### Pull Requests

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub

3. **Fill out the PR template** with:
   - Summary of changes
   - Related issues (if any)
   - Test plan
   - Screenshots (for visual changes)

4. **Wait for review** - maintainers will review your PR and may request changes

### PR Guidelines

- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update documentation if needed
- Ensure all CI checks pass
- Respond to review feedback promptly

## Project Structure

```
termcharts/
├── pkg/termcharts/     # Public library package
│   ├── chart.go        # Base interfaces and types
│   ├── options.go      # Functional options
│   ├── bar.go          # Bar chart implementation
│   ├── sparkline.go    # Sparkline implementation
│   ├── pie.go          # Pie chart implementation
│   ├── line.go         # Line chart implementation
│   └── style.go        # Colors and themes
├── internal/           # Internal utilities
├── cmd/termcharts/     # CLI application
├── examples/           # Example programs
├── docs/               # Documentation
└── Makefile            # Build commands
```

## Adding a New Chart Type

1. Create `pkg/termcharts/newchart.go` with the implementation
2. Create `pkg/termcharts/newchart_test.go` with tests
3. Add CLI command in `cmd/termcharts/newchart.go`
4. Create example in `examples/newchart/main.go`
5. Add documentation in `docs/newchart.md`
6. Update README.md with examples
7. Update docs/status.md with progress

## Reporting Issues

### Bug Reports

Include:
- termcharts version
- Go version
- Operating system
- Terminal emulator
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if applicable

### Feature Requests

Include:
- Use case description
- Proposed API (if applicable)
- Examples of similar features in other tools

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Assume good intentions

## Questions?

- Open an issue for questions about contributing
- Check existing issues and PRs before creating new ones

Thank you for contributing to termcharts!
