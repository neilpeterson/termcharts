// Last modified: 2026-01-03

package termcharts

import (
	"fmt"
	"math"
	"strings"

	"github.com/neilpeterson/termcharts/internal/util"
)

// BarChart represents a bar chart visualization.
// Bar charts can be rendered horizontally or vertically and support
// single or multiple data series with grouped or stacked modes.
type BarChart struct {
	opts *Options
}

// BarMode specifies how multiple series are displayed in a bar chart.
type BarMode int

const (
	// BarModeGrouped displays bars for each series side-by-side.
	BarModeGrouped BarMode = iota
	// BarModeStacked displays bars for each series stacked on top of each other.
	BarModeStacked
)

const unknownString = "unknown"

// String returns the string representation of the BarMode.
func (b BarMode) String() string {
	switch b {
	case BarModeGrouped:
		return "grouped"
	case BarModeStacked:
		return "stacked"
	default:
		return unknownString
	}
}

// ASCII characters for bar rendering when Unicode is not supported.
const barCharASCII = '#'

// NewBarChart creates a new bar chart with the given options.
// At minimum, data must be provided via WithData option.
//
// Example:
//
//	bar := termcharts.NewBarChart(
//	    termcharts.WithData([]float64{10, 25, 15, 30}),
//	    termcharts.WithLabels([]string{"Q1", "Q2", "Q3", "Q4"}),
//	)
//	fmt.Println(bar.Render())
func NewBarChart(opts ...Option) *BarChart {
	options := NewOptions(opts...)
	return &BarChart{
		opts: options,
	}
}

// Render generates the bar chart as a multi-line string.
func (b *BarChart) Render() string {
	// Validate data
	if len(b.opts.Data) == 0 && len(b.opts.Series) == 0 {
		return ""
	}

	// Render based on direction
	if b.opts.Direction == Horizontal {
		return b.renderHorizontal()
	}
	return b.renderVertical()
}

// renderHorizontal renders a horizontal bar chart.
//
//nolint:gocyclo // Complex by nature; splitting would harm readability
func (b *BarChart) renderHorizontal() string {
	data := b.opts.Data
	labels := b.opts.Labels

	// Check for invalid values
	if !util.AllValid(data) {
		return ""
	}

	// Determine character set based on style
	useUnicode := b.shouldUseUnicode()

	// Get color settings
	colorEnabled := b.isColorEnabled()
	theme := b.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Find max value for scaling
	maxVal := findMax(data)
	if maxVal == 0 {
		maxVal = 1 // Avoid division by zero
	}

	// Calculate bar width (leave room for labels and values)
	maxLabelWidth := 0
	if b.opts.ShowAxes && len(labels) > 0 {
		maxLabelWidth = maxStringLength(labels) + 1
	}

	valueWidth := 0
	if b.opts.ShowValues {
		valueWidth = len(fmt.Sprintf(" %.1f", maxVal)) + 1
	}

	// Calculate available width for bars
	barWidth := b.opts.Width - maxLabelWidth - valueWidth - 2
	if barWidth < 1 {
		barWidth = 20 // Minimum bar width
	}

	var result strings.Builder

	// Render title if provided
	if b.opts.Title != "" {
		titleText := b.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n")
	}

	// Render each bar
	for i, val := range data {
		// Render label
		if b.opts.ShowAxes {
			label := ""
			if i < len(labels) {
				label = labels[i]
			}
			labelText := fmt.Sprintf("%-*s ", maxLabelWidth, label)
			if colorEnabled {
				labelText = Colorize(labelText, theme.Muted, true)
			}
			result.WriteString(labelText)
		}

		// Calculate bar length
		barLen := int(float64(barWidth) * (val / maxVal))
		if barLen < 0 {
			barLen = 0
		}

		// Render bar
		bar := b.renderBar(barLen, barWidth, useUnicode, colorEnabled, theme.Primary)
		result.WriteString(bar)

		// Render value
		if b.opts.ShowValues {
			valueText := fmt.Sprintf(" %.1f", val)
			if colorEnabled {
				valueText = Colorize(valueText, theme.Muted, true)
			}
			result.WriteString(valueText)
		}

		result.WriteString("\n")
	}

	return result.String()
}

// renderBar renders a single horizontal bar with the given length.
func (b *BarChart) renderBar(length, maxWidth int, useUnicode bool, colorEnabled bool, color string) string {
	var bar strings.Builder

	if useUnicode {
		// Render full blocks
		fullBlocks := length
		for i := 0; i < fullBlocks && i < maxWidth; i++ {
			char := string('█')
			if colorEnabled {
				char = Colorize(char, color, true)
			}
			bar.WriteString(char)
		}
	} else {
		// ASCII mode - use '#' characters
		for i := 0; i < length && i < maxWidth; i++ {
			char := string(barCharASCII)
			if colorEnabled {
				char = Colorize(char, color, true)
			}
			bar.WriteString(char)
		}
	}

	return bar.String()
}

// renderVertical renders a vertical bar chart.
//
//nolint:gocyclo // Complex by nature; splitting would harm readability
func (b *BarChart) renderVertical() string {
	data := b.opts.Data
	labels := b.opts.Labels

	// Check for invalid values
	if !util.AllValid(data) {
		return ""
	}

	// Determine character set based on style
	useUnicode := b.shouldUseUnicode()

	// Get color settings
	colorEnabled := b.isColorEnabled()
	theme := b.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Find max value for scaling
	maxVal := findMax(data)
	if maxVal == 0 {
		maxVal = 1
	}

	// Calculate bar height
	barHeight := b.opts.Height
	if b.opts.Title != "" {
		barHeight-- // Leave room for title
	}
	if b.opts.ShowAxes && len(labels) > 0 {
		barHeight-- // Leave room for labels
	}
	if barHeight < 3 {
		barHeight = 10 // Minimum height
	}

	var result strings.Builder

	// Render title if provided
	if b.opts.Title != "" {
		titleText := b.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n")
	}

	// Calculate bar widths
	barWidth := 3 // Width of each bar column
	spacing := 1  // Space between bars

	// Render bars from top to bottom
	for row := barHeight; row > 0; row-- {
		for i, val := range data {
			// Calculate how many rows this bar should fill
			barRows := int(float64(barHeight) * (val / maxVal))

			// Determine if this row should have a bar
			if row <= barRows {
				// Render bar
				char := b.renderVerticalBar(useUnicode, colorEnabled, theme.Primary)
				result.WriteString(strings.Repeat(char, barWidth))
			} else {
				// Render empty space
				result.WriteString(strings.Repeat(" ", barWidth))
			}

			// Add spacing between bars (except after last bar)
			if i < len(data)-1 {
				result.WriteString(strings.Repeat(" ", spacing))
			}
		}
		result.WriteString("\n")
	}

	// Render labels if enabled
	if b.opts.ShowAxes && len(labels) > 0 {
		for i := range data {
			label := ""
			if i < len(labels) {
				label = labels[i]
				// Truncate or pad to bar width
				if len(label) > barWidth {
					label = label[:barWidth]
				} else {
					label = fmt.Sprintf("%-*s", barWidth, label)
				}
			} else {
				label = strings.Repeat(" ", barWidth)
			}

			labelText := label
			if colorEnabled {
				labelText = Colorize(labelText, theme.Muted, true)
			}
			result.WriteString(labelText)

			// Add spacing between labels
			if i < len(data)-1 {
				result.WriteString(strings.Repeat(" ", spacing))
			}
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderVerticalBar renders a single character for a vertical bar.
func (b *BarChart) renderVerticalBar(useUnicode bool, colorEnabled bool, color string) string {
	char := string(barCharASCII)
	if useUnicode {
		char = string('█')
	}

	if colorEnabled {
		return Colorize(char, color, true)
	}
	return char
}

// shouldUseUnicode determines whether to use Unicode characters based on style.
func (b *BarChart) shouldUseUnicode() bool {
	if b.opts.Style == StyleASCII {
		return false
	} else if b.opts.Style == StyleUnicode {
		return true
	}
	// StyleAuto - detect Unicode support
	return util.SupportsUnicode()
}

// isColorEnabled determines whether colors should be used.
func (b *BarChart) isColorEnabled() bool {
	if b.opts.ColorEnabled != nil {
		return *b.opts.ColorEnabled
	}
	return util.SupportsColor()
}

// findMax finds the maximum value in a slice of floats.
func findMax(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for _, v := range data[1:] {
		if !math.IsNaN(v) && !math.IsInf(v, 0) && v > max {
			max = v
		}
	}
	return max
}

// maxStringLength returns the length of the longest string in a slice.
func maxStringLength(strings []string) int {
	max := 0
	for _, s := range strings {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

// Bar is a convenience function that creates and renders a horizontal bar chart.
//
// Example:
//
//	fmt.Println(termcharts.Bar([]float64{10, 25, 15, 30}))
func Bar(data []float64) string {
	bar := NewBarChart(WithData(data))
	return bar.Render()
}

// BarWithLabels is a convenience function that creates a horizontal bar chart with labels.
//
// Example:
//
//	fmt.Println(termcharts.BarWithLabels(
//	    []float64{10, 25, 15, 30},
//	    []string{"Q1", "Q2", "Q3", "Q4"},
//	))
func BarWithLabels(data []float64, labels []string) string {
	bar := NewBarChart(
		WithData(data),
		WithLabels(labels),
	)
	return bar.Render()
}

// BarVertical is a convenience function that creates a vertical bar chart.
//
// Example:
//
//	fmt.Println(termcharts.BarVertical([]float64{10, 25, 15, 30}))
func BarVertical(data []float64) string {
	bar := NewBarChart(
		WithData(data),
		WithDirection(Vertical),
	)
	return bar.Render()
}
