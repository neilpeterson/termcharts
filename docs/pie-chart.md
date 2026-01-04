# Pie Chart Documentation

Pie charts display data as proportional slices of a whole, making it easy to visualize the relative size of different categories. The termcharts implementation renders pie charts using a proportional bar visualization with a color-coded legend.

## Features

- **Proportional Visualization**: Data displayed as a segmented bar showing relative proportions
- **Color-Coded Legend**: Each slice has a unique color with label and percentage
- **Value Display**: Optionally show numeric values alongside percentages
- **Unicode and ASCII Modes**: Automatic detection or manual selection
- **Title Support**: Add descriptive titles to charts
- **Customizable Width**: Control chart width for different terminal sizes
- **Theme Support**: Multiple color themes available

## Code Structure

### Location
- **Library**: `pkg/termcharts/pie.go`
- **Tests**: `pkg/termcharts/pie_test.go`
- **CLI**: `cmd/termcharts/pie.go`
- **Examples**: `examples/pie-chart/main.go`

### Types

```go
type PieChart struct {
    opts *Options
}

type Slice struct {
    Label      string
    Value      float64
    Percentage float64
    Color      string
}
```

### Key Functions

- `NewPieChart(opts ...Option) *PieChart` - Create a new pie chart
- `Render() string` - Generate the chart output
- `Pie(data []float64) string` - Convenience function for quick pie charts
- `PieWithLabels(data []float64, labels []string) string` - Pie chart with labels
- `PieWithValues(data []float64, labels []string) string` - Pie chart with values displayed

## Library API

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
    // Simple pie chart
    data := []float64{30, 25, 20, 15, 10}
    fmt.Println(termcharts.Pie(data))
}
```

### With Labels

```go
data := []float64{30, 25, 20, 15, 10}
labels := []string{"Chrome", "Firefox", "Safari", "Edge", "Other"}

pie := termcharts.NewPieChart(
    termcharts.WithData(data),
    termcharts.WithLabels(labels),
)
fmt.Println(pie.Render())
```

Output:
```
        ■■■■●●●●●
     ◇◇■■■■■●●●●●●●●
   ◇◇◇◇◇■■■■●●●●●●●●●●
  ◇◇◇◇◇◇◇◇■■●●●●●●●●●●●     ● Chrome    30.0%
  ◇◇◇◇◇◇◇◇◇■●●●●●●●●●●●     ○ Firefox   25.0%
 ◇◇◇◇◇◇◇◇◇◇◇●●●●●●●●●●●●    ◆ Safari    20.0%
  ◆◆◆◆◆◆◆◆◆◆○○○○○○○●●●●     ◇ Edge      15.0%
  ◆◆◆◆◆◆◆◆◆○○○○○○○○○○○○     ■ Other     10.0%
   ◆◆◆◆◆◆◆◆○○○○○○○○○○○
     ◆◆◆◆◆○○○○○○○○○○
        ◆○○○○○○○○
```

### With Title and Values

```go
pie := termcharts.NewPieChart(
    termcharts.WithData([]float64{1250.50, 980.25, 750.00, 520.75}),
    termcharts.WithLabels([]string{"Product A", "Product B", "Product C", "Product D"}),
    termcharts.WithTitle("Revenue by Product ($)"),
    termcharts.WithShowValues(true),
)
fmt.Println(pie.Render())
```

Output:
```
Revenue by Product ($)

        ◇◇◇◇●●●●●
     ◇◇◇◇◇◇◇●●●●●●●●
   ◆◇◇◇◇◇◇◇◇●●●●●●●●●●
  ◆◆◆◆◆◇◇◇◇◇●●●●●●●●●●●     ● Product A  35.7% [1250.5]
  ◆◆◆◆◆◆◆◆◇◇●●●●●●●●●●●     ○ Product B  28.0% [980.2]
 ◆◆◆◆◆◆◆◆◆◆◆●●●●●●●●●●●●    ◆ Product C  21.4% [750.0]
  ◆◆◆◆◆◆◆◆○○○○○●●●●●●●●     ◇ Product D  14.9% [520.8]
  ◆◆◆◆◆◆○○○○○○○○○○●●●●●
   ◆◆◆○○○○○○○○○○○○○○●●
     ○○○○○○○○○○○○○○○
        ○○○○○○○○○
```

### With Color

```go
pie := termcharts.NewPieChart(
    termcharts.WithData([]float64{50, 30, 20}),
    termcharts.WithLabels([]string{"Desktop", "Mobile", "Tablet"}),
    termcharts.WithColor(true),
)
fmt.Println(pie.Render())
```

![Pie Chart with Color](images/pie-chart-color.png)

### ASCII Mode

```go
pie := termcharts.NewPieChart(
    termcharts.WithData([]float64{50, 30, 20}),
    termcharts.WithLabels([]string{"Yes", "No", "Maybe"}),
    termcharts.WithStyle(termcharts.StyleASCII),
)
fmt.Println(pie.Render())
```

### Convenience Functions

```go
// Quick pie chart
fmt.Println(termcharts.Pie([]float64{50, 30, 20}))

// With labels
fmt.Println(termcharts.PieWithLabels(
    []float64{50, 30, 20},
    []string{"A", "B", "C"},
))

// With values
fmt.Println(termcharts.PieWithValues(
    []float64{100, 75, 50},
    []string{"Large", "Medium", "Small"},
))
```

## CLI Usage

The `termcharts pie` command provides a convenient way to create pie charts from the command line.

### Basic Commands

```bash
# Simple pie chart from arguments
termcharts pie 30 25 20 15 10

# With labels
termcharts pie 30 25 20 15 10 --labels "Chrome,Firefox,Safari,Edge,Other"

# With title
termcharts pie 30 25 20 --title "Market Share"

# With values displayed
termcharts pie 30 25 20 --show-values --labels "A,B,C"
```

### Reading from Files

```bash
# From a file
termcharts pie data.txt

# With labels
termcharts pie data.txt --labels "Q1,Q2,Q3,Q4"
```

### Reading from Stdin

```bash
# From pipeline
echo "30 25 20 15 10" | termcharts pie

# From command output
seq 1 5 | termcharts pie --labels "One,Two,Three,Four,Five"
```

### Styling Options

```bash
# ASCII mode for compatibility
termcharts pie 50 30 20 --ascii

# Colored output
termcharts pie 50 30 20 --color

# Disable colors
termcharts pie 50 30 20 --no-color

# Custom width
termcharts pie 50 30 20 --width 60
```

### Complete Example

```bash
termcharts pie 35 28 22 15 \
    --labels "North,South,East,West" \
    --title "Sales by Region" \
    --show-values \
    --color
```

## Options Reference

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithData()` | []float64 | required | Data values to visualize |
| `WithLabels()` | []string | auto | Labels for each slice |
| `WithTitle()` | string | none | Chart title |
| `WithWidth()` | int | 80 | Chart width in columns |
| `WithShowValues()` | bool | false | Display numeric values |
| `WithStyle()` | RenderStyle | StyleAuto | ASCII, Unicode, or Auto |
| `WithColor()` | bool | auto | Enable/disable colors |
| `WithTheme()` | *Theme | DefaultTheme | Color theme |

## CLI Flags Reference

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--width` | `-w` | int | 80 | Chart width in characters |
| `--labels` | `-l` | string | "" | Comma-separated labels |
| `--title` | `-t` | string | "" | Chart title |
| `--show-values` | | bool | false | Display numeric values |
| `--color` | `-c` | bool | false | Enable colored output |
| `--no-color` | | bool | false | Disable colored output |
| `--ascii` | | bool | false | Use ASCII characters only |

## Implementation Details

### Rendering

The pie chart is rendered as two components:

1. **Proportional Bar**: A horizontal bar divided into segments, where each segment's width corresponds to its percentage of the total
2. **Legend**: A list showing each slice's color indicator, label, optional value, and percentage

### Percentages

- All values are converted to percentages based on the total sum
- Percentages always sum to 100% (with rounding)
- Negative values are treated as zero
- Zero values show as 0.0% but may not have visible bar segments

### Color Coding

When colors are enabled:
- Each slice gets a unique color from the theme's series colors
- Colors cycle through the series if there are more slices than colors
- The legend shows colored indicators matching the bar segments

### Character Sets

**With Colors Enabled:**
- Asterisk character: *
- Each slice colored differently from theme

**Without Colors (Unicode):**
- Different symbols per slice: ●, ○, ◆, ◇, ■, □, ▲, △, ★, ☆
- Symbols distinguish slices visually

**Without Colors (ASCII):**
- Different characters per slice: *, o, #, x, +, @, =, ~, %, &
- Maximum compatibility with all terminals

## Best Practices

1. **Keep Slices Limited**
   - Aim for 5-7 slices maximum for readability
   - Combine smaller values into an "Other" category

2. **Label Clearly**
   - Use short, descriptive labels
   - Labels are automatically aligned in the legend

3. **Use Colors Wisely**
   - Colors help distinguish slices
   - Consider accessibility when using color themes

4. **Show Values When Needed**
   - Use `--show-values` when exact numbers matter
   - Percentages are always shown

## Examples

### Market Share Analysis

```go
pie := termcharts.NewPieChart(
    termcharts.WithData([]float64{65, 20, 10, 5}),
    termcharts.WithLabels([]string{"Market Leader", "Challenger", "Niche", "Others"}),
    termcharts.WithTitle("Market Share 2024"),
    termcharts.WithShowValues(true),
    termcharts.WithColor(true),
)
fmt.Println(pie.Render())
```

### Budget Breakdown

```bash
termcharts pie 35 25 20 12 8 \
    --labels "Engineering,Marketing,Operations,Sales,Support" \
    --title "Department Budget Allocation" \
    --show-values \
    --color
```

### Survey Results

```go
results := []float64{45.2, 30.8, 15.5, 8.5}
options := []string{"Strongly Agree", "Agree", "Disagree", "Strongly Disagree"}

pie := termcharts.NewPieChart(
    termcharts.WithData(results),
    termcharts.WithLabels(options),
    termcharts.WithTitle("Customer Satisfaction Survey"),
    termcharts.WithWidth(70),
)
fmt.Println(pie.Render())
```

## Testing

The pie chart implementation includes comprehensive tests:

```bash
# Run all pie chart tests
go test ./pkg/termcharts -v -run TestPie

# Test specific functionality
go test ./pkg/termcharts -v -run TestPieChart_Render_WithLabels
```

## Performance

Pie chart rendering is efficient:
- Time complexity: O(n) where n is the number of slices
- Memory usage: Minimal, uses string builders
- No allocations in hot paths

## Compatibility

- **Terminals**: Works in all modern terminals
- **Unicode Support**: Auto-detects and falls back to ASCII
- **Colors**: Auto-detects color support
- **Operating Systems**: Cross-platform (macOS, Linux, Windows)

## See Also

- [Bar Chart Documentation](bar-chart.md)
- [Sparkline Documentation](sparkline.md)
- [API Reference](api-reference.md)
- [Examples Directory](../examples/)
