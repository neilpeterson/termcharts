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

```bash
go install github.com/neilpeterson/termcharts/cmd/termcharts@latest
```

Or build from source:

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
        termcharts.WithWidth(60),
        termcharts.WithColor(true),
    )
    fmt.Println(chart.Render())

    // Sparkline
    fmt.Println(termcharts.Spark([]float64{1, 5, 2, 8, 3, 7, 4, 6}))
    // Output: ▁▅▂█▃▇▄▆
}
```

### CLI Usage

```bash
# From arguments
termcharts bar 10 25 15 30 --labels "Q1,Q2,Q3,Q4"

# From file
termcharts bar data.json

# From stdin (pipe-friendly)
cat data.csv | termcharts line --x date --y value

# Sparkline
echo "1 5 2 8 3 7" | termcharts spark
# Output: ▁▅▂█▃▇▄▆

# Real-world usage: CPU monitoring
ps aux | awk '{sum+=$3} END {print sum}' | termcharts spark --color
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

### Coming Soon

- **Bar Charts** - Horizontal and vertical bar charts with labels
- **Line Charts** - ASCII and high-resolution Braille line charts
- **Heatmaps** - 2D data visualization with color gradients
- **Scatter Plots** - XY coordinate plotting
- **Gauges** - Progress bars and percentage indicators

## Roadmap

### v0.1.0 - Core Library + Sparklines ✓
- [x] Project scaffolding
- [x] Core types and interfaces
- [x] Terminal utilities
- [x] Sparkline implementation with comprehensive tests
- [x] Sparkline CLI command (`termcharts spark`)

### v0.2.0 - Bar Charts
- [ ] Horizontal bar charts
- [ ] Vertical bar charts
- [ ] Grouped/stacked variants
- [ ] CLI `bar` subcommand

### v0.3.0 - Line Charts
- [ ] ASCII line charts
- [ ] Braille high-resolution charts
- [ ] Multi-series support
- [ ] CLI `line` subcommand

See [docs/status.md](docs/status.md) for detailed status and milestones.

## Documentation

- **[Sparkline Guide](docs/sparkline.md)** - Complete sparkline documentation with examples, API reference, and CLI usage
- **[Project Status](docs/status.md)** - Current development status and roadmap
- **[API Documentation](https://pkg.go.dev/github.com/neilpeterson/termcharts)** - GoDoc reference (coming soon)

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

- Go 1.19 or higher
- Terminal with Unicode support (recommended, ASCII fallback available)

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
