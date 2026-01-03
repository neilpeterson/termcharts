# Go Documentation Standards

Apply these documentation standards to ALL Go code written in this project.

## File Header

Every Go file MUST start with a header comment containing the last modified date:

```go
// Last modified: 2026-01-02

package termcharts
```

**Rules:**
- Place at the very top of the file, before package declaration
- Use ISO 8601 date format (YYYY-MM-DD)
- Update the date whenever the file is modified
- Leave a blank line between header and package comment/declaration

## Package Documentation (One File Per Package)

```go
// Package termcharts provides ASCII and Unicode terminal charting capabilities.
//
// termcharts supports multiple chart types including bar charts, line charts,
// sparklines, and more. Charts can be rendered using pure ASCII for maximum
// compatibility or Unicode block/Braille characters for higher fidelity.
//
// Basic usage:
//
//	data := []float64{10, 25, 15, 30}
//	chart := termcharts.NewBarChart(
//	    termcharts.WithData(data),
//	    termcharts.WithWidth(60),
//	)
//	fmt.Println(chart.Render())
package termcharts
```

## Exported Types

```go
// Chart represents a terminal-based data visualization.
// All chart types implement this interface.
type Chart interface {
    // Render generates the chart as a string ready for terminal output.
    Render() string
}

// BarChart displays data as horizontal or vertical bars.
// It supports labels, colors, and automatic scaling.
type BarChart struct {
    data   []float64
    labels []string
    opts   *Options
}
```

## Exported Functions/Methods

Start with the function name, use complete sentences:

```go
// NewBarChart creates a new bar chart with the given options.
// If no options are provided, sensible defaults are used.
func NewBarChart(opts ...Option) *BarChart {
    // implementation
}

// Render generates the bar chart as a multi-line string.
// It automatically scales data to fit the configured width.
func (b *BarChart) Render() string {
    // implementation
}

// WithWidth sets the maximum chart width in terminal columns.
// Default is 80 or auto-detected terminal width.
func WithWidth(width int) Option {
    return func(o *Options) {
        o.Width = width
    }
}
```

## Inline Comments

Use sparingly, only for non-obvious logic:

```go
func (b *BarChart) Render() string {
    // Normalize data to 0-1 range for scaling
    max := maxValue(b.data)
    if max == 0 {
        max = 1 // Avoid division by zero
    }

    // Calculate bar length: scale to width, subtract label space
    barLen := int((val / max) * float64(b.opts.Width-20))
}
```

## Constants/Variables

```go
// Default chart dimensions
const (
    DefaultWidth  = 80  // Maximum width in terminal columns
    DefaultHeight = 24  // Maximum height in terminal rows
)

// ErrInvalidData indicates the provided data cannot be visualized.
var ErrInvalidData = errors.New("invalid or empty data")
```

## Key Rules

1. **Every file MUST have a header with last modified date**
2. **All exported symbols MUST have documentation comments**
3. **Start comments with the name of the symbol being documented**
4. **Use complete sentences**
5. **Include usage examples for complex APIs**
6. **Keep inline comments minimal - code should be self-explanatory**
7. **Comments explain WHY, code explains WHAT**
8. **First sentence is crucial - it appears in package indexes**
9. **Follow godoc conventions** (rendered by `go doc`)

## Internal/Unexported Code

Short comments or none if obvious:

```go
// maxValue returns the maximum value in the slice
func maxValue(data []float64) float64 {
    // implementation
}

// No comment needed for obvious internal helpers
func clamp(val, min, max int) int {
    if val < min {
        return min
    }
    if val > max {
        return max
    }
    return val
}
```

---

**Apply these standards automatically when writing or modifying Go code.**

**IMPORTANT:** When modifying any Go file, always update the "Last modified" date in the file header to the current date.
