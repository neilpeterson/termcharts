# API Reference

Complete API reference for the termcharts library.

## Table of Contents

- [Core Interfaces](#core-interfaces)
- [Options Pattern](#options-pattern)
- [Render Styles](#render-styles)
- [Themes and Colors](#themes-and-colors)
- [Data Types](#data-types)
- [Error Handling](#error-handling)

## Core Interfaces

### Chart Interface

All chart types implement the `Chart` interface:

```go
type Chart interface {
    Render() string
}
```

**Methods:**
- `Render() string` - Generates the chart as a string ready for terminal output

**Example:**

```go
chart := termcharts.NewSparkline(
    termcharts.WithData([]float64{1, 5, 2, 8, 3}),
)
output := chart.Render()
fmt.Println(output)
```

## Options Pattern

termcharts uses the functional options pattern for clean, composable configuration.

### NewOptions

```go
func NewOptions(opts ...Option) *Options
```

Creates a new Options struct with sensible defaults.

**Default Values:**
- `Width`: 80
- `Height`: 24
- `Style`: StyleAuto
- `Direction`: Horizontal
- `ShowValues`: false
- `ShowAxes`: true

### Available Options

#### WithData

```go
func WithData(data []float64) Option
```

Sets the primary data series for the chart.

**Example:**

```go
opts := termcharts.NewOptions(
    termcharts.WithData([]float64{10, 20, 30, 40}),
)
```

#### WithLabels

```go
func WithLabels(labels []string) Option
```

Sets labels for each data point. The number of labels should match the number of data points.

**Example:**

```go
chart := termcharts.NewBarChart(
    termcharts.WithData([]float64{10, 25, 15, 30}),
    termcharts.WithLabels([]string{"Q1", "Q2", "Q3", "Q4"}),
)
```

#### WithWidth

```go
func WithWidth(width int) Option
```

Sets the maximum chart width in terminal columns. Use 0 to auto-detect terminal width.

**Example:**

```go
chart := termcharts.NewSparkline(
    termcharts.WithData(data),
    termcharts.WithWidth(50),
)
```

#### WithHeight

```go
func WithHeight(height int) Option
```

Sets the maximum chart height in terminal rows. Use 0 to auto-detect terminal height.

#### WithTitle

```go
func WithTitle(title string) Option
```

Sets an optional title displayed above the chart.

#### WithColor

```go
func WithColor(enabled bool) Option
```

Enables or disables ANSI color output. If not set, color support is auto-detected.

**Example:**

```go
// Force colors on
chart := termcharts.NewSparkline(
    termcharts.WithData(data),
    termcharts.WithColor(true),
)

// Force colors off
chart := termcharts.NewSparkline(
    termcharts.WithData(data),
    termcharts.WithColor(false),
)
```

#### WithStyle

```go
func WithStyle(style RenderStyle) Option
```

Sets the rendering style. StyleAuto automatically selects the best style based on terminal capabilities.

**Available Styles:**
- `StyleAuto` - Auto-detect best style
- `StyleASCII` - Pure ASCII characters
- `StyleUnicode` - Unicode block characters
- `StyleBraille` - Unicode Braille patterns (line charts)

**Example:**

```go
// Force ASCII mode for maximum compatibility
chart := termcharts.NewSparkline(
    termcharts.WithData(data),
    termcharts.WithStyle(termcharts.StyleASCII),
)

// Use Unicode for better visuals
chart := termcharts.NewSparkline(
    termcharts.WithData(data),
    termcharts.WithStyle(termcharts.StyleUnicode),
)
```

#### WithDirection

```go
func WithDirection(dir Direction) Option
```

Sets the chart orientation (Horizontal or Vertical).

#### WithTheme

```go
func WithTheme(theme *Theme) Option
```

Sets the color theme for the chart.

**Example:**

```go
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithTheme(termcharts.DarkTheme),
)
```

## Render Styles

### RenderStyle Type

```go
type RenderStyle int

const (
    StyleAuto RenderStyle = iota
    StyleASCII
    StyleUnicode
    StyleBraille
)
```

**String Method:**

```go
func (s RenderStyle) String() string
```

Returns the string representation ("auto", "ascii", "unicode", or "braille").

## Themes and Colors

### Theme Type

```go
type Theme struct {
    Primary    string
    Secondary  string
    Accent     string
    Muted      string
    Background string
    Text       string
    Series     []string
}
```

**Methods:**

```go
func (t *Theme) GetSeriesColor(index int) string
```

Returns the color for a data series at the given index. Cycles through Series colors.

### Predefined Themes

- `DefaultTheme` - Standard terminal colors
- `DarkTheme` - Optimized for dark backgrounds
- `LightTheme` - Optimized for light backgrounds
- `MonochromeTheme` - Grayscale only

**Example:**

```go
chart := termcharts.NewBarChart(
    termcharts.WithData(data),
    termcharts.WithTheme(termcharts.DarkTheme),
)
```

### Colorize Function

```go
func Colorize(text, color string, colorEnabled bool) string
```

Wraps text with ANSI color codes. Returns text unchanged if colorEnabled is false.

**Supported Colors:**
- black, red, green, yellow, orange (alias: yellow)
- blue, magenta, purple (alias: magenta), cyan
- white, gray, grey, brown (alias: red)

## Data Types

### Series

```go
type Series struct {
    Label string
    Data  []float64
    Color string
}
```

Represents a labeled data series for multi-series charts.

### Direction

```go
type Direction int

const (
    Horizontal Direction = iota
    Vertical
)
```

**String Method:**

```go
func (d Direction) String() string
```

Returns "horizontal" or "vertical".

## Error Handling

### Common Errors

```go
var (
    ErrEmptyData         = errors.New("data cannot be empty")
    ErrInvalidData       = errors.New("data contains invalid values")
    ErrInvalidDimensions = errors.New("chart dimensions too small")
)
```

Charts automatically validate data and return errors for:
- Empty data sets
- Invalid values (NaN, Inf)
- Dimension constraints too small to render

## See Also

- **[Sparkline Guide](sparkline.md)** - Detailed sparkline documentation
- **[Project Status](status.md)** - Development roadmap
- **[README](../README.md)** - Quick start guide
