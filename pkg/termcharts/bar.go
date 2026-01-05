package termcharts

import (
	"fmt"
	"math"
	"strings"

	"github.com/neilpeterson/termcharts/internal"
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

	// If multi-series, render grouped or stacked
	if len(b.opts.Series) > 0 {
		if b.opts.Direction == Horizontal {
			return b.renderHorizontalMultiSeries()
		}
		return b.renderVerticalMultiSeries()
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
	if !internal.AllValid(data) {
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
	if !internal.AllValid(data) {
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
	return internal.SupportsUnicode()
}

// isColorEnabled determines whether colors should be used.
func (b *BarChart) isColorEnabled() bool {
	if b.opts.ColorEnabled != nil {
		return *b.opts.ColorEnabled
	}
	return internal.SupportsColor()
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

// BarGrouped is a convenience function that creates a grouped bar chart.
//
// Example:
//
//	fmt.Println(termcharts.BarGrouped([]termcharts.Series{
//	    {Label: "2023", Data: []float64{10, 20, 15}},
//	    {Label: "2024", Data: []float64{12, 25, 18}},
//	}))
func BarGrouped(series []Series) string {
	bar := NewBarChart(
		WithSeries(series),
		WithBarMode(BarModeGrouped),
	)
	return bar.Render()
}

// BarStacked is a convenience function that creates a stacked bar chart.
//
// Example:
//
//	fmt.Println(termcharts.BarStacked([]termcharts.Series{
//	    {Label: "Product A", Data: []float64{10, 20, 15}},
//	    {Label: "Product B", Data: []float64{5, 10, 8}},
//	}))
func BarStacked(series []Series) string {
	bar := NewBarChart(
		WithSeries(series),
		WithBarMode(BarModeStacked),
	)
	return bar.Render()
}

// renderHorizontalMultiSeries renders a horizontal bar chart with multiple series.
//
//nolint:gocyclo // Complex by nature; splitting would harm readability
func (b *BarChart) renderHorizontalMultiSeries() string {
	series := b.opts.Series
	labels := b.opts.Labels

	// Validate all series data
	for _, s := range series {
		if !internal.AllValid(s.Data) {
			return ""
		}
	}

	// Determine character set and color settings
	useUnicode := b.shouldUseUnicode()
	colorEnabled := b.isColorEnabled()
	theme := b.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Find the number of categories (max data points across series)
	numCategories := 0
	for _, s := range series {
		if len(s.Data) > numCategories {
			numCategories = len(s.Data)
		}
	}

	// Calculate max value based on bar mode
	maxVal := b.calculateMaxValue(series)
	if maxVal == 0 {
		maxVal = 1
	}

	// Calculate label width
	maxLabelWidth := 0
	if b.opts.ShowAxes && len(labels) > 0 {
		maxLabelWidth = maxStringLength(labels) + 1
	}

	// Calculate bar width
	barWidth := b.opts.Width - maxLabelWidth - 2
	if barWidth < 1 {
		barWidth = 20
	}

	var result strings.Builder

	// Render title
	if b.opts.Title != "" {
		titleText := b.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n")
	}

	// Render based on mode
	if b.opts.BarMode == BarModeStacked {
		b.renderHorizontalStacked(&result, series, labels, numCategories, maxVal, barWidth, maxLabelWidth, useUnicode, colorEnabled, theme)
	} else {
		b.renderHorizontalGrouped(&result, series, labels, numCategories, maxVal, barWidth, maxLabelWidth, useUnicode, colorEnabled, theme)
	}

	// Render legend if enabled
	if b.opts.ShowLegend {
		result.WriteString("\n")
		for i, s := range series {
			color := theme.GetSeriesColor(i)
			if s.Color != "" {
				color = s.Color
			}
			legendChar := "█"
			if !useUnicode {
				legendChar = "#"
			}
			if colorEnabled {
				legendChar = Colorize(legendChar, color, true)
			}
			result.WriteString(fmt.Sprintf("%s %s  ", legendChar, s.Label))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderHorizontalGrouped renders horizontal grouped bars.
func (b *BarChart) renderHorizontalGrouped(result *strings.Builder, series []Series, labels []string, numCategories int, maxVal float64, barWidth, maxLabelWidth int, useUnicode, colorEnabled bool, theme *Theme) {
	for cat := 0; cat < numCategories; cat++ {
		// Render label for this category
		if b.opts.ShowAxes {
			label := ""
			if cat < len(labels) {
				label = labels[cat]
			}
			labelText := fmt.Sprintf("%-*s ", maxLabelWidth, label)
			if colorEnabled {
				labelText = Colorize(labelText, theme.Muted, true)
			}
			result.WriteString(labelText)
		}

		// Render bars for each series side by side
		for i, s := range series {
			val := 0.0
			if cat < len(s.Data) {
				val = s.Data[cat]
			}

			barLen := int(float64(barWidth/len(series)) * (val / maxVal))
			if barLen < 0 {
				barLen = 0
			}

			color := theme.GetSeriesColor(i)
			if s.Color != "" {
				color = s.Color
			}

			bar := b.renderBar(barLen, barWidth/len(series), useUnicode, colorEnabled, color)
			result.WriteString(bar)
		}
		result.WriteString("\n")
	}
}

// renderHorizontalStacked renders horizontal stacked bars.
func (b *BarChart) renderHorizontalStacked(result *strings.Builder, series []Series, labels []string, numCategories int, maxVal float64, barWidth, maxLabelWidth int, useUnicode, colorEnabled bool, theme *Theme) {
	for cat := 0; cat < numCategories; cat++ {
		// Render label for this category
		if b.opts.ShowAxes {
			label := ""
			if cat < len(labels) {
				label = labels[cat]
			}
			labelText := fmt.Sprintf("%-*s ", maxLabelWidth, label)
			if colorEnabled {
				labelText = Colorize(labelText, theme.Muted, true)
			}
			result.WriteString(labelText)
		}

		// Render stacked bars (each series stacked horizontally)
		for i, s := range series {
			val := 0.0
			if cat < len(s.Data) {
				val = s.Data[cat]
			}

			barLen := int(float64(barWidth) * (val / maxVal))
			if barLen < 0 {
				barLen = 0
			}

			color := theme.GetSeriesColor(i)
			if s.Color != "" {
				color = s.Color
			}

			bar := b.renderBar(barLen, barWidth, useUnicode, colorEnabled, color)
			result.WriteString(bar)
		}
		result.WriteString("\n")
	}
}

// renderVerticalMultiSeries renders a vertical bar chart with multiple series.
//
//nolint:gocyclo // Complex by nature; splitting would harm readability
func (b *BarChart) renderVerticalMultiSeries() string {
	series := b.opts.Series
	labels := b.opts.Labels

	// Validate all series data
	for _, s := range series {
		if !internal.AllValid(s.Data) {
			return ""
		}
	}

	// Determine character set and color settings
	useUnicode := b.shouldUseUnicode()
	colorEnabled := b.isColorEnabled()
	theme := b.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Find the number of categories
	numCategories := 0
	for _, s := range series {
		if len(s.Data) > numCategories {
			numCategories = len(s.Data)
		}
	}

	// Calculate max value based on bar mode
	maxVal := b.calculateMaxValue(series)
	if maxVal == 0 {
		maxVal = 1
	}

	// Calculate bar height
	barHeight := b.opts.Height
	if b.opts.Title != "" {
		barHeight--
	}
	if b.opts.ShowAxes && len(labels) > 0 {
		barHeight--
	}
	if b.opts.ShowLegend {
		barHeight -= 2
	}
	if barHeight < 3 {
		barHeight = 10
	}

	var result strings.Builder

	// Render title
	if b.opts.Title != "" {
		titleText := b.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n")
	}

	// Render based on mode
	if b.opts.BarMode == BarModeStacked {
		b.renderVerticalStacked(&result, series, labels, numCategories, maxVal, barHeight, useUnicode, colorEnabled, theme)
	} else {
		b.renderVerticalGrouped(&result, series, labels, numCategories, maxVal, barHeight, useUnicode, colorEnabled, theme)
	}

	// Render legend if enabled
	if b.opts.ShowLegend {
		result.WriteString("\n")
		for i, s := range series {
			color := theme.GetSeriesColor(i)
			if s.Color != "" {
				color = s.Color
			}
			legendChar := "█"
			if !useUnicode {
				legendChar = "#"
			}
			if colorEnabled {
				legendChar = Colorize(legendChar, color, true)
			}
			result.WriteString(fmt.Sprintf("%s %s  ", legendChar, s.Label))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderVerticalGrouped renders vertical grouped bars.
func (b *BarChart) renderVerticalGrouped(result *strings.Builder, series []Series, labels []string, numCategories int, maxVal float64, barHeight int, useUnicode, colorEnabled bool, theme *Theme) {
	barWidth := 3                  // Width of each bar
	groupSpacing := 2              // Space between groups
	barSpacing := 0                // Space between bars in a group
	groupWidth := len(series)*barWidth + (len(series)-1)*barSpacing

	// Render bars from top to bottom
	for row := barHeight; row > 0; row-- {
		for cat := 0; cat < numCategories; cat++ {
			for i, s := range series {
				val := 0.0
				if cat < len(s.Data) {
					val = s.Data[cat]
				}

				barRows := int(float64(barHeight) * (val / maxVal))
				color := theme.GetSeriesColor(i)
				if s.Color != "" {
					color = s.Color
				}

				if row <= barRows {
					char := b.renderVerticalBar(useUnicode, colorEnabled, color)
					result.WriteString(strings.Repeat(char, barWidth))
				} else {
					result.WriteString(strings.Repeat(" ", barWidth))
				}

				// Add spacing between bars in group
				if i < len(series)-1 {
					result.WriteString(strings.Repeat(" ", barSpacing))
				}
			}

			// Add spacing between groups
			if cat < numCategories-1 {
				result.WriteString(strings.Repeat(" ", groupSpacing))
			}
		}
		result.WriteString("\n")
	}

	// Render labels
	if b.opts.ShowAxes && len(labels) > 0 {
		for cat := 0; cat < numCategories; cat++ {
			label := ""
			if cat < len(labels) {
				label = labels[cat]
				if len(label) > groupWidth {
					label = label[:groupWidth]
				} else {
					label = fmt.Sprintf("%-*s", groupWidth, label)
				}
			} else {
				label = strings.Repeat(" ", groupWidth)
			}

			labelText := label
			if colorEnabled {
				labelText = Colorize(labelText, theme.Muted, true)
			}
			result.WriteString(labelText)

			if cat < numCategories-1 {
				result.WriteString(strings.Repeat(" ", groupSpacing))
			}
		}
		result.WriteString("\n")
	}
}

// renderVerticalStacked renders vertical stacked bars.
func (b *BarChart) renderVerticalStacked(result *strings.Builder, series []Series, labels []string, numCategories int, maxVal float64, barHeight int, useUnicode, colorEnabled bool, theme *Theme) {
	barWidth := 3  // Width of each bar
	spacing := 1   // Space between bars

	// Pre-calculate the stacked heights for each category
	stackedHeights := make([][]int, numCategories)
	for cat := 0; cat < numCategories; cat++ {
		stackedHeights[cat] = make([]int, len(series))
		cumulative := 0.0
		for i, s := range series {
			val := 0.0
			if cat < len(s.Data) {
				val = s.Data[cat]
			}
			cumulative += val
			stackedHeights[cat][i] = int(float64(barHeight) * (cumulative / maxVal))
		}
	}

	// Render bars from top to bottom
	for row := barHeight; row > 0; row-- {
		for cat := 0; cat < numCategories; cat++ {
			// Find which series this row belongs to (from top to bottom)
			seriesIdx := -1
			for i := len(series) - 1; i >= 0; i-- {
				if row <= stackedHeights[cat][i] {
					prevHeight := 0
					if i > 0 {
						prevHeight = stackedHeights[cat][i-1]
					}
					if row > prevHeight {
						seriesIdx = i
						break
					}
				}
			}

			if seriesIdx >= 0 {
				color := theme.GetSeriesColor(seriesIdx)
				if series[seriesIdx].Color != "" {
					color = series[seriesIdx].Color
				}
				char := b.renderVerticalBar(useUnicode, colorEnabled, color)
				result.WriteString(strings.Repeat(char, barWidth))
			} else {
				result.WriteString(strings.Repeat(" ", barWidth))
			}

			if cat < numCategories-1 {
				result.WriteString(strings.Repeat(" ", spacing))
			}
		}
		result.WriteString("\n")
	}

	// Render labels
	if b.opts.ShowAxes && len(labels) > 0 {
		for cat := 0; cat < numCategories; cat++ {
			label := ""
			if cat < len(labels) {
				label = labels[cat]
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

			if cat < numCategories-1 {
				result.WriteString(strings.Repeat(" ", spacing))
			}
		}
		result.WriteString("\n")
	}
}

// calculateMaxValue calculates the maximum value based on bar mode.
func (b *BarChart) calculateMaxValue(series []Series) float64 {
	if b.opts.BarMode == BarModeStacked {
		// For stacked, find max sum across categories
		numCategories := 0
		for _, s := range series {
			if len(s.Data) > numCategories {
				numCategories = len(s.Data)
			}
		}

		maxSum := 0.0
		for cat := 0; cat < numCategories; cat++ {
			sum := 0.0
			for _, s := range series {
				if cat < len(s.Data) {
					sum += s.Data[cat]
				}
			}
			if sum > maxSum {
				maxSum = sum
			}
		}
		return maxSum
	}

	// For grouped, find max individual value
	maxVal := 0.0
	for _, s := range series {
		for _, v := range s.Data {
			if !math.IsNaN(v) && !math.IsInf(v, 0) && v > maxVal {
				maxVal = v
			}
		}
	}
	return maxVal
}
