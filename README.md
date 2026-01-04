# termcharts

Beautiful ASCII/Unicode terminal charting library and CLI tool for Go.

## Overview

`termcharts` is a lightweight, dependency-minimal Go library for rendering beautiful charts in the terminal. Perfect for DevOps dashboards, CLI tools, data pipelines, and any application that needs quick data visualization without leaving the terminal.

## Features

- Multiple chart types: bar charts, line charts, sparklines, heatmaps, scatter plots, and gauges
- Progressive enhancement: ASCII → Unicode → Colors (auto-detect or override)
- Responsive: auto-detects terminal width and scales accordingly
- Both library and CLI tool included
- Zero-configuration defaults with extensive customization options
- Pipe-friendly CLI for shell integration

## Installation

### As a Library

```bash
go get github.com/neilpeterson/termcharts/pkg/termcharts
```

### As a CLI Tool

#### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/neilpeterson/termcharts/releases).

**macOS:**
```bash
# Intel Macs
curl -L https://github.com/neilpeterson/termcharts/releases/latest/download/termcharts_darwin_amd64.tar.gz | tar xz
sudo mv termcharts /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/neilpeterson/termcharts/releases/latest/download/termcharts_darwin_arm64.tar.gz | tar xz
sudo mv termcharts /usr/local/bin/
```

**Linux:**
```bash
# AMD64
curl -L https://github.com/neilpeterson/termcharts/releases/latest/download/termcharts_linux_amd64.tar.gz | tar xz
sudo mv termcharts /usr/local/bin/

# ARM64
curl -L https://github.com/neilpeterson/termcharts/releases/latest/download/termcharts_linux_arm64.tar.gz | tar xz
sudo mv termcharts /usr/local/bin/
```

**Windows:**
Download the `.zip` file from the [releases page](https://github.com/neilpeterson/termcharts/releases) and extract manually.

#### Install with Go

```bash
go install github.com/neilpeterson/termcharts/cmd/termcharts@latest
```

#### Build from Source

```bash
git clone https://github.com/neilpeterson/termcharts.git
cd termcharts
make build
# Binary will be in bin/termcharts
```

## Quick Start

### Library Usage

```go
package main

import (
    "fmt"
    "github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
    // Quick one-liner
    fmt.Println(termcharts.Bar([]float64{10, 25, 15, 30}))

    // With options
    chart := termcharts.NewBarChart(
        termcharts.WithData([]float64{10, 25, 15, 30}),
        termcharts.WithLabels([]string{"Q1", "Q2", "Q3", "Q4"}),
        termcharts.WithTitle("Quarterly Sales"),
        termcharts.WithWidth(60),
        termcharts.WithColor(true),
    )
    fmt.Println(chart.Render())
}
```

Output:
```
Quarterly Sales
Q1  ██████████████████████
Q2  ███████████████████████████████████████████████████
Q3  █████████████████████████████████
Q4  █████████████████████████████████████████████████████████████████
```

### CLI Usage

```bash
# From arguments
termcharts bar 10 25 15 30 --labels "Q1,Q2,Q3,Q4"

# With title and values displayed
termcharts bar 120 98 145 --labels "North,South,East" --title "Regional Sales" --show-values

# Vertical bar chart
termcharts bar 10 25 15 30 --vertical --labels "Q1,Q2,Q3,Q4"

# From file (JSON or CSV)
termcharts bar data.json
```

## Chart Types

### Sparklines ✓

Compact, inline charts that visualize data trends in a single line. Perfect for dashboards and monitoring.

```go
// Library
fmt.Println(termcharts.Spark([]float64{1, 5, 2, 8, 3, 7, 4, 6}))
// Output: ▁▅▂█▃▇▄▆
```

```bash
# CLI
termcharts spark 10 20 30 25 15 --color --width 50
echo "1 5 2 8" | termcharts spark
```

**[→ Full Sparkline Documentation](docs/sparkline.md)**

### Bar Charts ✓

Horizontal and vertical bar charts for comparing values. Supports labels, titles, and value display.

```go
// Library
chart := termcharts.NewBarChart(
    termcharts.WithData([]float64{10, 25, 15, 30}),
    termcharts.WithLabels([]string{"Q1", "Q2", "Q3", "Q4"}),
    termcharts.WithTitle("Quarterly Sales"),
)
fmt.Println(chart.Render())
```

Output:
```
Quarterly Sales
Q1  ██████████████████████
Q2  ███████████████████████████████████████████████████
Q3  █████████████████████████████████
Q4  █████████████████████████████████████████████████████████████████
```

**Vertical bars:**
```bash
termcharts bar 10 25 15 30 --vertical --labels "Q1,Q2,Q3,Q4"
```

Output:
```
            ███
            ███
            ███
    ███     ███
    ███     ███
    ███     ███
    ███ ███ ███
    ███ ███ ███
███ ███ ███ ███
███ ███ ███ ███
Q1  Q2  Q3  Q4
```

**With values:**
```bash
termcharts bar 120 98 145 --labels "North,South,East" --show-values
```

Output:
```
North  ████████████████████████████████████████████████████████ 120.0
South  ██████████████████████████████████████████████ 98.0
East   ███████████████████████████████████████████████████████████████████ 145.0
```

**[→ Full Bar Chart Documentation](docs/bar-chart.md)**

### Pie Charts ✓

Pie charts display proportional data with color-coded segments and percentages.

```go
// Library
pie := termcharts.NewPieChart(
    termcharts.WithData([]float64{30, 25, 20, 15, 10}),
    termcharts.WithLabels([]string{"Chrome", "Firefox", "Safari", "Edge", "Other"}),
    termcharts.WithTitle("Browser Market Share"),
)
fmt.Println(pie.Render())
```

Output:
```
Browser Market Share

              █
       ██████████████
    ████████████████████
  ██████████████████████████
 ████████████████████████████
 ████████████████████████████
████████████████████████████████
 ██████████████████████████████
 ██████████████████████████████
  ████████████████████████████
   ██████████████████████████
    ████████████████████████
      ██████████████████████
        ████████████████████
           ██████████████
              ████

  █ Chrome   ( 30.0%)
  █ Firefox  ( 25.0%)
  █ Safari   ( 20.0%)
  █ Edge     ( 15.0%)
  █ Other    ( 10.0%)
```

```bash
# CLI
termcharts pie 30 25 20 15 10 --labels "Chrome,Firefox,Safari,Edge,Other"
termcharts pie 50 30 20 --title "Survey Results" --show-values
```

**[→ Full Pie Chart Documentation](docs/pie-chart.md)**

### Coming Soon

- **Line Charts** - ASCII and high-resolution Braille line charts
- **Heatmaps** - 2D data visualization with color gradients
- **Scatter Plots** - XY coordinate plotting
- **Gauges** - Progress bars and percentage indicators

## Roadmap

### v0.1.0 - Core Library + Sparklines ✓
- [x] Project scaffolding and build system
- [x] Core types and interfaces
- [x] Terminal utilities (size, color, Unicode detection)
- [x] Sparkline implementation with comprehensive tests (95% coverage)
- [x] Sparkline CLI command (`termcharts spark`)
- [x] Unit tests for all core components (95% coverage)
- [x] GitHub Actions CI/CD (multi-platform, multi-version testing)
- [x] Complete API documentation

### v0.2.0 - Bar Charts ✓
- [x] Horizontal bar charts
- [x] Vertical bar charts
- [x] Labels and axis rendering
- [x] Value display option
- [x] CLI `bar` subcommand
- [x] Comprehensive tests and documentation
- [ ] Grouped/stacked variants (deferred to v0.5.0)

### v0.3.0 - Pie Charts ✓
- [x] Pie chart with Unicode/ASCII modes
- [x] Percentage and value display
- [x] Legend and color coding
- [x] CLI `pie` subcommand
- [x] Comprehensive tests and documentation

### v0.4.0 - Line Charts
- [ ] ASCII line charts
- [ ] Braille high-resolution charts
- [ ] Multi-series support
- [ ] CLI `line` subcommand

See [docs/status.md](docs/status.md) for detailed status and milestones.

## Documentation

- **[API Reference](docs/api-reference.md)** - Complete API documentation with all options and examples
- **[Sparkline Guide](docs/sparkline.md)** - Complete sparkline documentation with examples, API reference, and CLI usage
- **[Bar Chart Guide](docs/bar-chart.md)** - Complete bar chart documentation with examples, API reference, and CLI usage
- **[Pie Chart Guide](docs/pie-chart.md)** - Complete pie chart documentation with examples, API reference, and CLI usage
- **[Project Status](docs/status.md)** - Current development status and roadmap
- **[GoDoc](https://pkg.go.dev/github.com/neilpeterson/termcharts)** - Generated API documentation (coming soon)

## Development

```bash
# Run tests
make test

# Run tests with coverage
make cover

# Run linter
make lint

# Build binary
make build

# Install to $GOPATH/bin
make install

# Cross-compile for all platforms
make release

# Clean build artifacts
make clean
```

## Design Principles

- **Sensible defaults**: Charts look good with zero configuration
- **Progressive enhancement**: Works everywhere, looks best where supported
- **Functional options**: Clean, composable API using the options pattern
- **Pipe-friendly**: CLI integrates seamlessly with Unix pipelines
- **Zero allocation hot paths**: Efficient rendering for real-time updates

## Requirements

- Go 1.21 or higher
- Terminal with Unicode support (recommended, ASCII fallback available)

## Testing

termcharts has comprehensive test coverage (95%+):

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detector
go test -race ./...
```

The project uses GitHub Actions for CI/CD with automated testing across:
- Multiple platforms: Linux, macOS, Windows
- Multiple Go versions: 1.21, 1.22, 1.23

## License

MIT License - See LICENSE file for details

## Contributing

Contributions welcome! Please see CONTRIBUTING.md for guidelines.

## Inspiration

This project draws inspiration from:
- [termgraph](https://github.com/mkaz/termgraph)
- [asciigraph](https://github.com/guptarohit/asciigraph)
- [sparkline](https://github.com/jolmg/cliplot)
- [pterm](https://github.com/pterm/pterm)
