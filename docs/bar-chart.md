# Bar Chart Documentation

Bar charts are one of the core visualization types in termcharts. They display data as rectangular bars with lengths proportional to the values they represent. Bar charts can be rendered horizontally or vertically and support labels, titles, and various styling options.

## Features

- **Horizontal and Vertical Orientation**: Display bars left-to-right or bottom-to-top
- **Labels and Titles**: Add meaningful labels to each bar and a title to the chart
- **Value Display**: Optionally show numeric values alongside bars
- **Unicode and ASCII Modes**: Automatic detection or manual selection
- **Colored Output**: Color support with auto-detection
- **Flexible Sizing**: Customizable width and height

## Code Structure

### Location
- **Library**: `pkg/termcharts/bar.go`
- **Tests**: `pkg/termcharts/bar_test.go`
- **CLI**: `cmd/termcharts/bar.go`
- **Examples**: `examples/bar-chart/main.go`

### Types

```go
type BarChart struct {
    opts *Options
}

type BarMode int
const (
    BarModeGrouped BarMode = iota
    BarModeStacked
)
```

### Key Functions

- `NewBarChart(opts ...Option) *BarChart` - Create a new bar chart
- `Render() string` - Generate the chart output
- `Bar(data []float64) string` - Convenience function for quick bar charts
- `BarWithLabels(data []float64, labels []string) string` - Bar chart with labels
- `BarVertical(data []float64) string` - Vertical bar chart

## Library API

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
    // Simple horizontal bar chart
    data := []float64{10, 25, 15, 30}
    chart := termcharts.NewBarChart(
        termcharts.WithData(data),
    )
    fmt.Println(chart.Render())
}
```

### With Labels and Title

```go
data := []float64{45, 67, 52, 78}
labels := []string{"Q1", "Q2", "Q3", "Q4"}
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithLabels(labels),
    termcharts.WithTitle("Quarterly Sales"),
)
fmt.Println(chart.Render())
```

### Vertical Bar Chart

```go
data := []float64{12, 19, 25, 18, 22, 30}
labels := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithLabels(labels),
    termcharts.WithDirection(termcharts.Vertical),
    termcharts.WithHeight(12),
)
fmt.Println(chart.Render())
```

### With Values Displayed

```go
data := []float64{120.5, 98.3, 145.7, 132.1}
labels := []string{"North", "South", "East", "West"}
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithLabels(labels),
    termcharts.WithTitle("Regional Sales"),
    termcharts.WithShowValues(true),
)
fmt.Println(chart.Render())
```

### ASCII Mode

```go
data := []float64{8, 15, 12, 20}
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithStyle(termcharts.StyleASCII),
)
fmt.Println(chart.Render())
```

### Custom Width

```go
data := []float64{10, 25, 15, 30}
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithWidth(60),
)
fmt.Println(chart.Render())
```

### Convenience Functions

```go
// Quick bar chart
fmt.Println(termcharts.Bar([]float64{5, 10, 15, 20}))

// With labels
fmt.Println(termcharts.BarWithLabels(
    []float64{25, 40, 35, 50},
    []string{"Alpha", "Beta", "Gamma", "Delta"},
))

// Vertical
fmt.Println(termcharts.BarVertical([]float64{10, 20, 15, 25}))
```

## CLI Usage

The `termcharts bar` command provides a convenient way to create bar charts from the command line.

### Basic Commands

```bash
# Simple bar chart from arguments
termcharts bar 10 25 15 30

# With labels
termcharts bar 10 25 15 30 --labels "Q1,Q2,Q3,Q4"

# Vertical orientation
termcharts bar 10 25 15 30 --vertical

# With title and values
termcharts bar 10 25 15 30 --title "Sales Report" --show-values
```

### Reading from Files

```bash
# From a file
termcharts bar data.txt

# File with labels (format: value,label per line)
cat > sales.txt << EOF
45,Q1
67,Q2
52,Q3
78,Q4
EOF
termcharts bar sales.txt --labels "Q1,Q2,Q3,Q4"
```

### Reading from Stdin

```bash
# From pipeline
echo "10 20 30 40" | termcharts bar

# From command output
seq 5 | termcharts bar --vertical
```

### Styling Options

```bash
# ASCII mode for compatibility
termcharts bar 10 20 30 --ascii

# Colored output
termcharts bar 10 20 30 --color

# Disable colors
termcharts bar 10 20 30 --no-color

# Custom width
termcharts bar 10 20 30 --width 60

# Custom height (vertical mode)
termcharts bar 10 20 30 --vertical --height 20
```

### Complete Example

```bash
termcharts bar 120 98 145 132 \
    --labels "North,South,East,West" \
    --title "Regional Sales" \
    --show-values \
    --width 70 \
    --color
```

## Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithData()` | []float64 | required | Data values to visualize |
| `WithLabels()` | []string | none | Labels for each bar |
| `WithTitle()` | string | none | Chart title |
| `WithDirection()` | Direction | Horizontal | Orientation (Horizontal/Vertical) |
| `WithWidth()` | int | 80 | Chart width in columns |
| `WithHeight()` | int | 24 | Chart height in rows (vertical mode) |
| `WithShowValues()` | bool | false | Display numeric values |
| `WithShowAxes()` | bool | true | Display axes and labels |
| `WithStyle()` | RenderStyle | StyleAuto | ASCII, Unicode, or Auto |
| `WithColor()` | bool | auto | Enable/disable colors |
| `WithTheme()` | *Theme | DefaultTheme | Color theme |

## CLI Flags Reference

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--width` | `-w` | int | 80 | Chart width in characters |
| `--height` | | int | 15 | Chart height in rows (vertical mode) |
| `--vertical` | `-v` | bool | false | Render vertical bar chart |
| `--labels` | `-l` | string | "" | Comma-separated labels |
| `--title` | `-t` | string | "" | Chart title |
| `--show-values` | | bool | false | Display numeric values |
| `--color` | `-c` | bool | false | Enable colored output |
| `--no-color` | | bool | false | Disable colored output |
| `--ascii` | | bool | false | Use ASCII characters only |

## Implementation Details

### Horizontal Bar Chart

- Each bar is rendered on its own line
- Bar length is proportional to the value relative to the maximum
- Labels are left-aligned with consistent padding
- Values are displayed to the right of bars when enabled
- Uses Unicode block characters (█) or ASCII (#) depending on mode

### Vertical Bar Chart

- Bars are rendered from bottom to top
- Each bar occupies 3 characters width with 1 character spacing
- Labels are centered below each bar
- Height is adjustable via the `--height` flag
- Bars fill from bottom based on value proportion

### Scaling

- All values are scaled relative to the maximum value in the dataset
- Zero values are handled correctly (no bar drawn)
- Negative values are currently treated as zero (future enhancement)

### Character Sets

**Unicode Mode:**
- Full block: █
- Provides clean, solid bars

**ASCII Mode:**
- Hash character: #
- Maximum compatibility with all terminals

## Best Practices

1. **Choose the Right Orientation**
   - Horizontal: Better for long labels, comparing values
   - Vertical: Better for time series, many data points

2. **Label Wisely**
   - Keep labels short and descriptive
   - Use abbreviations for space-constrained displays

3. **Width Considerations**
   - Default 80 columns works for most terminals
   - Increase width for more detail or many data points
   - Decrease for embedding in narrow contexts

4. **Value Display**
   - Use `--show-values` when exact numbers matter
   - Omit values for cleaner visual comparison

5. **Colors**
   - Colors enhance readability but may not work in all terminals
   - Use `--ascii` and `--no-color` for maximum compatibility

## Examples

### Sales Dashboard

```go
sales := []float64{45000, 67000, 52000, 78000}
quarters := []string{"Q1", "Q2", "Q3", "Q4"}

chart := termcharts.NewBarChart(
    termcharts.WithData(sales),
    termcharts.WithLabels(quarters),
    termcharts.WithTitle("2024 Sales by Quarter ($)"),
    termcharts.WithShowValues(true),
    termcharts.WithWidth(70),
)
fmt.Println(chart.Render())
```

### Server Metrics

```bash
#!/bin/bash
# Monitor server CPU usage
top -l 1 | grep "CPU usage" | awk '{print $3}' | sed 's/%//' | \
termcharts bar --title "CPU Usage" --show-values
```

### Build Times Comparison

```go
buildTimes := []float64{45.2, 38.7, 52.1, 41.3, 36.8}
commits := []string{"abc123", "def456", "ghi789", "jkl012", "mno345"}

chart := termcharts.NewBarChart(
    termcharts.WithData(buildTimes),
    termcharts.WithLabels(commits),
    termcharts.WithTitle("Build Times (seconds)"),
    termcharts.WithShowValues(true),
)
fmt.Println(chart.Render())
```

## Future Enhancements

The following features are planned for future releases:

- **Grouped Bar Charts**: Display multiple series side-by-side
- **Stacked Bar Charts**: Stack multiple series on top of each other
- **Negative Value Support**: Properly handle and display negative values
- **Horizontal Reference Lines**: Add grid lines or reference markers
- **Custom Bar Characters**: Allow user-defined characters for bars
- **Gradient Fills**: Use multiple characters to create gradient effects
- **Logarithmic Scale**: Support for logarithmic scaling

## Testing

The bar chart implementation includes comprehensive tests:

```bash
# Run all bar chart tests
go test ./pkg/termcharts -v -run TestBar

# Test specific functionality
go test ./pkg/termcharts -v -run TestBarChart_Render_WithLabels
```

## Performance

Bar chart rendering is highly efficient:
- Time complexity: O(n) where n is the number of data points
- Memory usage: Minimal, uses string builders
- No allocations in hot paths for single-series charts

## Compatibility

- **Terminals**: Works in all modern terminals
- **Unicode Support**: Auto-detects and falls back to ASCII
- **Colors**: Auto-detects 256-color and truecolor support
- **Operating Systems**: Cross-platform (macOS, Linux, Windows)

## See Also

- [Sparkline Documentation](sparkline.md)
- [API Reference](api-reference.md)
- [Examples Directory](../examples/)
