# Line Charts

Line charts visualize data trends over time or sequential data points. termcharts supports multiple rendering modes for line charts, from basic ASCII to high-resolution Braille patterns.

## Features

- **ASCII Mode**: Uses basic ASCII characters (`/`, `\`, `-`, `|`, `*`) for maximum terminal compatibility
- **Unicode Mode**: Uses box-drawing characters (`─`, `│`, `╱`, `╲`, `•`) for cleaner lines
- **Braille Mode**: High-resolution rendering using Unicode Braille patterns (2x4 dots per character)
- **Multi-series Support**: Plot multiple data series on the same chart with automatic color assignment
- **Customizable**: Configurable width, height, colors, titles, axes, and labels
- **Legend**: Automatic legend generation for multi-series charts

## Library Usage

### Basic Line Chart

```go
import "github.com/neilpeterson/termcharts/pkg/termcharts"

data := []float64{1, 5, 2, 8, 3, 7, 4, 6}
line := termcharts.NewLineChart(
    termcharts.WithData(data),
    termcharts.WithWidth(60),
    termcharts.WithHeight(10),
)
fmt.Println(line.Render())
```

### With Title and Labels

```go
data := []float64{150, 230, 180, 290}
labels := []string{"Q1", "Q2", "Q3", "Q4"}

line := termcharts.NewLineChart(
    termcharts.WithData(data),
    termcharts.WithLabels(labels),
    termcharts.WithTitle("Quarterly Revenue"),
    termcharts.WithShowAxes(true),
    termcharts.WithColor(true),
)
fmt.Println(line.Render())
```

### Braille High-Resolution

```go
data := []float64{1, 5, 2, 8, 3, 7, 4, 6}
line := termcharts.NewLineChart(
    termcharts.WithData(data),
    termcharts.WithStyle(termcharts.StyleBraille),
    termcharts.WithColor(true),
)
fmt.Println(line.Render())
```

### Multi-Series Chart

```go
series := []termcharts.Series{
    {Label: "Revenue", Data: []float64{100, 150, 130, 180}},
    {Label: "Costs", Data: []float64{80, 90, 100, 110}},
}
line := termcharts.NewLineChart(
    termcharts.WithSeries(series),
    termcharts.WithTitle("Financial Comparison"),
    termcharts.WithColor(true),
)
fmt.Println(line.Render())
```

### Convenience Functions

```go
// Quick line chart with defaults
fmt.Println(termcharts.Line([]float64{1, 5, 2, 8, 3, 7}))

// High-resolution Braille line chart
fmt.Println(termcharts.LineBraille([]float64{1, 5, 2, 8, 3, 7}))

// Multi-series line chart
series := []termcharts.Series{
    {Label: "A", Data: []float64{1, 2, 3}},
    {Label: "B", Data: []float64{3, 2, 1}},
}
fmt.Println(termcharts.LineMultiSeries(series))
```

## CLI Usage

### Basic Usage

```bash
# From command-line arguments
termcharts line 1 5 2 8 3 7 4 6

# From file
termcharts line data.txt

# From stdin
echo "1 5 2 8 3 7" | termcharts line
```

### Rendering Modes

```bash
# ASCII mode (maximum compatibility)
termcharts line 1 5 2 8 3 7 --ascii

# High-resolution Braille
termcharts line 1 5 2 8 3 7 --braille
```

### Customization

```bash
# With title
termcharts line 10 25 15 30 --title "Sales Trend"

# With X-axis labels
termcharts line 10 25 15 30 --labels "Jan,Feb,Mar,Apr"

# Custom dimensions
termcharts line 1 5 2 8 3 7 --width 80 --height 15

# With color
termcharts line 1 5 2 8 3 7 --color

# Hide axes
termcharts line 1 5 2 8 3 7 --axes=false

# Different themes
termcharts line 1 5 2 8 3 7 --color --theme dark
```

## Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `WithData` | `[]float64` | - | Primary data series |
| `WithSeries` | `[]Series` | - | Multiple data series |
| `WithWidth` | `int` | 80 | Chart width in characters |
| `WithHeight` | `int` | 24 | Chart height in rows |
| `WithTitle` | `string` | "" | Chart title |
| `WithLabels` | `[]string` | - | X-axis labels |
| `WithStyle` | `RenderStyle` | Auto | ASCII, Unicode, or Braille |
| `WithColor` | `bool` | auto | Enable ANSI colors |
| `WithShowAxes` | `bool` | true | Show axes and labels |
| `WithTheme` | `*Theme` | Default | Color theme |

## Render Styles

### ASCII (`StyleASCII`)
```
    *
   / \
  /   \
 /     *
*
```

Uses characters: `/`, `\`, `-`, `|`, `*`

### Unicode (`StyleUnicode`)
```
    •
   ╱ ╲
  ╱   ╲
 ╱     •
•
```

Uses box-drawing characters for smoother lines.

### Braille (`StyleBraille`)
```
⠀⠀⠀⠀⢠⠳⡀⠀⠀⠀
⠀⠀⠀⢠⠃⠀⠱⡀⠀⠀
⠀⠀⢠⠃⠀⠀⠀⠱⡀⠀
⢠⠔⠁⠀⠀⠀⠀⠀⠈⠢
```

Each character cell contains a 2x4 dot matrix, providing 4x higher vertical resolution.

## Themes

Available color themes:
- `default` - Standard terminal colors (blue primary)
- `dark` - Optimized for dark backgrounds (cyan primary)
- `light` - Optimized for light backgrounds (blue primary)
- `mono` - Grayscale only

## Example Output

### ASCII Line Chart
```
    8.0     *
    6.3    / \    *
    4.7   /   \  / \
    3.0  /     \/   \  *
    1.3 *            \/
        ----------------
        A  B  C  D  E  F
```

### Braille Line Chart
```
    8.0 ⠀⠀⠀⠀⢠⠳⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
    6.0 ⠀⠀⠀⢠⠃⠀⠱⡀⠀⠀⡔⢄⠀⠀⠀⠀
    4.0 ⠀⠀⢠⠃⠀⠀⠀⠱⡀⢀⠎⠀⠱⡀⠀⢀
    2.0 ⠀⢠⠃⠀⠀⠀⠀⠀⠙⠁⠀⠀⠀⠈⠔⠁
        ──────────────────
```

## Running the Example

```bash
go run examples/line-chart/main.go
```
