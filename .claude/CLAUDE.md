# termcharts

ASCII/Unicode terminal charting library and CLI tool, written in Go.

## Workflow Rules

**At the start of every session:**
1. Read `docs/status.md` to understand current project state
2. Confirm with user what task/feature to work on
3. Update status.md to mark task as "IN PROGRESS"

**After completing any task:**
1. Update `docs/status.md`: mark complete, add notes, update "Recent Changes"
2. If new tasks discovered, add them to the backlog
3. Summarize what was done and suggest next steps

**When creating new features:**
1. **ALWAYS create a new feature branch FIRST** - Never make changes on main or existing branches
   - Use naming pattern: `feature/feature-name` (e.g., `feature/bar-charts`)
   - Command: `git checkout -b feature/feature-name`
   - This is a BLOCKING requirement - do not proceed without creating a branch
2. Add feature to status.md before starting implementation
3. Break into subtasks if complex (>1 hour of work)
4. Update status as each subtask completes
5. **After completing a feature**: Update `README.md` with usage examples and documentation
6. **When completing a chart type**: Create feature-specific documentation in `docs/` folder

**Documentation Maintenance:**
- `README.md` is the main entry point - keep it current with:
  - Installation instructions
  - Quick start examples for each chart type
  - CLI usage examples
  - Link to detailed docs in `docs/` folder
- Create detailed docs in `docs/` for each major feature (e.g., `docs/sparkline.md`)
- Include code structure, API examples, and CLI usage in feature docs

**Never skip status updates** - this is how we maintain continuity across sessions.

## Project Overview

Two deliverables:
1. **Library** (`termcharts`): Importable Go package for rendering charts to terminal
2. **CLI** (`termcharts`): Standalone binary using the library, with configurable data sources

Target users: DevOps engineers, data engineers, CLI enthusiasts who want beautiful terminal visualizations.

## Tech Stack

- Go 1.21+
- Minimal dependencies for core library
- Optional: `github.com/fatih/color` or `github.com/muesli/termenv` for colors
- CLI: `github.com/spf13/cobra` for commands
- Testing: standard `testing` package + `github.com/stretchr/testify`

## Architecture

```
termcharts/
├── pkg/
│   └── termcharts/           # Public library package
│       ├── chart.go          # Base Chart interface and common types
│       ├── options.go        # Functional options pattern
│       ├── bar.go            # Horizontal/vertical bar charts
│       ├── sparkline.go      # Inline sparklines
│       └── style.go          # Colors, themes, rendering modes
│       # Future: line.go, scatter.go, heatmap.go, gauge.go
├── internal/                 # Internal utilities (not part of public API)
│   ├── terminal.go           # Terminal size detection
│   └── scale.go              # Data scaling/normalization
├── cmd/
│   └── termcharts/
│       ├── main.go           # Entry point
│       ├── root.go           # Root command
│       ├── bar.go            # `termcharts bar` subcommand
│       └── spark.go          # `termcharts spark` subcommand
│       # Future: line.go, sources.go
├── examples/
│   ├── bar-chart/            # Bar chart demo
│   └── sparkline/            # Sparkline demo
├── docs/                     # Documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Design Principles

- **Sensible defaults**: Charts look good with zero configuration
- **Progressive enhancement**: ASCII → Unicode → Colors (auto-detect or override)
- **Functional options**: Configure via `WithWidth()`, `WithColor()`, etc.
- **Pipe-friendly**: CLI works with stdin/stdout for shell pipelines
- **Responsive**: Auto-detect terminal width, scale accordingly
- **Zero allocation hot paths**: Efficient rendering for real-time updates

## Chart Types (Priority Order)

1. Horizontal bar chart
2. Sparklines
3. Line chart (ASCII curves, then Braille high-res)
4. Vertical bar chart
5. Progress bars / gauges
6. Heatmap
7. Scatter plot

## Key Unicode Characters

```
Blocks: ░ ▒ ▓ █ ▄ ▀ ▌ ▐ ▁▂▃▄▅▆▇█
Braille: ⠀⠁⠂⠃⠄...⣿ (256 patterns, 2x4 dots per cell)
Box drawing: ─ │ ┌ ┐ └ ┘ ├ ┤ ┬ ┴ ┼ ╭ ╮ ╯ ╰
Arrows: ← → ↑ ↓ ↗ ↘ ▲ ▼
Misc: ● ○ ◐ ◑ ■ □ ▪ ▫
```

## Commands

```bash
# Build
go build -o termcharts ./cmd/termcharts

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter
golangci-lint run

# Install locally
go install ./cmd/termcharts

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o termcharts-linux ./cmd/termcharts
GOOS=darwin GOARCH=arm64 go build -o termcharts-darwin ./cmd/termcharts
GOOS=windows GOARCH=amd64 go build -o termcharts.exe ./cmd/termcharts
```

## Code Style

- Follow standard Go conventions (gofmt, golint)
- Use functional options for configuration
- Keep public API in `pkg/termcharts`, internals in `internal/`
- Table-driven tests
- Godoc comments on all exported types/functions

## API Design Goals

Library usage should be simple:
```go
import "github.com/neilpeterson/termcharts/pkg/termcharts"

// Quick one-liner
fmt.Println(termcharts.Bar([]float64{10, 25, 15, 30}))

// With options
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithLabels([]string{"Q1", "Q2", "Q3", "Q4"}),
    termcharts.WithWidth(60),
    termcharts.WithColor(true),
)
fmt.Println(chart.Render())

// Sparkline
fmt.Println(termcharts.Spark([]float64{1, 5, 2, 8, 3, 7, 4, 6}))
// Output: ▁▅▂█▃▇▄▆
```

CLI usage should be intuitive:
```bash
# From arguments
termcharts bar 10 25 15 30 --labels "Q1,Q2,Q3,Q4"

# From file
termcharts bar data.json

# From stdin
cat data.csv | termcharts line --x date --y value

# From API with live refresh
termcharts bar --url "https://api.example.com/metrics" --watch 5s

# Sparkline inline
echo "1 5 2 8 3 7" | termcharts spark
```

## Testing Approach

- Table-driven tests for each chart type
- Golden file tests: compare rendered output against expected `.golden` files
- Test terminal width handling with mocked sizes
- Fuzz testing for edge cases (empty data, huge values, unicode handling)

## Makefile Targets

```makefile
build      # Build binary
test       # Run tests
cover      # Run tests with coverage report
lint       # Run golangci-lint
install    # Install to $GOPATH/bin
deps       # Download and tidy dependencies
release    # Cross-compile for all platforms
clean      # Remove build artifacts
help       # Display available targets
```

## References

- Inspiration: `termgraph`, `asciigraph`, `sparkline`, `pterm`
- Unicode blocks: https://en.wikipedia.org/wiki/Block_Elements
- Braille patterns: https://en.wikipedia.org/wiki/Braille_Patterns
- Go terminal libs: `github.com/muesli/termenv`, `github.com/charmbracelet/lipgloss`