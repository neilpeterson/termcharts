package termcharts

import (
	"fmt"
	"math"
	"strings"

	"github.com/neilpeterson/termcharts/internal"
)

// PieChart represents a pie chart visualization.
// Pie charts display data as proportional slices of a circle,
// rendered using different characters for each slice.
type PieChart struct {
	opts *Options
}

// Slice represents a single slice of the pie chart.
type Slice struct {
	Label      string
	Value      float64
	Percentage float64
	Color      string
}

// Slice characters - uniform character for pie, different for legend when no color
var pieChar = '*'
var pieCharASCII = '*'

// Legend characters - different symbols for each slice (used when colors disabled)
var legendChars = []rune{'●', '○', '◆', '◇', '■', '□', '▲', '△', '★', '☆'}
var legendCharsASCII = []rune{'*', 'o', '#', 'x', '+', '@', '=', '~', '%', '&'}

// NewPieChart creates a new pie chart with the given options.
// At minimum, data must be provided via WithData option.
//
// Example:
//
//	pie := termcharts.NewPieChart(
//	    termcharts.WithData([]float64{30, 25, 20, 15, 10}),
//	    termcharts.WithLabels([]string{"Chrome", "Firefox", "Safari", "Edge", "Other"}),
//	)
//	fmt.Println(pie.Render())
func NewPieChart(opts ...Option) *PieChart {
	options := NewOptions(opts...)
	return &PieChart{
		opts: options,
	}
}

// Render generates the pie chart as a multi-line string.
func (p *PieChart) Render() string {
	// Validate data
	if len(p.opts.Data) == 0 {
		return ""
	}

	// Check for invalid values
	if !internal.AllValid(p.opts.Data) {
		return ""
	}

	// Calculate total and slices
	slices := p.calculateSlices()
	if len(slices) == 0 {
		return ""
	}

	// Get rendering settings
	colorEnabled := p.isColorEnabled()
	theme := p.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	var result strings.Builder

	// Render title if provided
	if p.opts.Title != "" {
		titleText := p.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n\n")
	}

	// Render circular pie with legend on the right
	pieWithLegend := p.renderCircularPieWithLegend(slices, colorEnabled, theme)
	result.WriteString(pieWithLegend)

	return result.String()
}

// calculateSlices calculates the slice data including percentages.
func (p *PieChart) calculateSlices() []Slice {
	data := p.opts.Data
	labels := p.opts.Labels

	// Calculate total
	total := 0.0
	for _, v := range data {
		if v > 0 {
			total += v
		}
	}

	if total == 0 {
		return nil
	}

	// Create slices
	slices := make([]Slice, len(data))
	for i, v := range data {
		label := fmt.Sprintf("Item %d", i+1)
		if i < len(labels) && labels[i] != "" {
			label = labels[i]
		}

		percentage := 0.0
		if v > 0 {
			percentage = (v / total) * 100
		}

		slices[i] = Slice{
			Label:      label,
			Value:      v,
			Percentage: percentage,
		}
	}

	return slices
}

// renderCircularPieWithLegend renders a circular pie chart with legend on the right.
func (p *PieChart) renderCircularPieWithLegend(slices []Slice, colorEnabled bool, theme *Theme) string {
	// Determine character set
	useUnicode := p.shouldUseUnicode()

	// Pie uses uniform character, legend uses different chars when no color
	pChar := pieChar
	lChars := legendChars
	if !useUnicode {
		pChar = pieCharASCII
		lChars = legendCharsASCII
	}

	// Pie chart dimensions
	radius := 6
	aspectRatio := 2.0 // Terminal chars are ~2x taller than wide

	// Calculate cumulative angles for each slice (starting at top, going clockwise)
	angles := make([]float64, len(slices)+1)
	angles[0] = -math.Pi / 2 // Start at 12 o'clock
	for i, slice := range slices {
		angles[i+1] = angles[i] + (slice.Percentage/100)*2*math.Pi
	}

	// Build pie chart rows
	pieRows := make([]string, 0)
	for y := -radius; y <= radius; y++ {
		var row strings.Builder
		for x := -int(float64(radius) * aspectRatio); x <= int(float64(radius)*aspectRatio); x++ {
			// Calculate actual position accounting for aspect ratio
			actualX := float64(x) / aspectRatio
			actualY := float64(y)

			// Calculate distance from center
			distance := math.Sqrt(actualX*actualX + actualY*actualY)

			// Check if point is inside the circle
			if distance <= float64(radius)-0.5 {
				// Calculate angle of this point
				angle := math.Atan2(actualY, actualX)

				// Find which slice this angle belongs to
				sliceIndex := p.findSliceForAngle(angle, angles, len(slices))

				// Apply color if enabled - use uniform character with different colors
				if colorEnabled {
					color := theme.GetSeriesColor(sliceIndex)
					row.WriteString(Colorize(string(pChar), color, true))
				} else {
					// Without colors, use different characters to distinguish slices
					char := lChars[sliceIndex%len(lChars)]
					row.WriteString(string(char))
				}
			} else {
				row.WriteString(" ")
			}
		}
		pieRows = append(pieRows, row.String())
	}

	// Build legend entries
	legendEntries := make([]string, len(slices))
	for i, slice := range slices {
		var entry strings.Builder

		if colorEnabled {
			// With colors: use uniform char with slice color
			color := theme.GetSeriesColor(i)
			entry.WriteString(Colorize(string(pChar), color, true))
		} else {
			// Without colors: use different chars to match pie
			char := lChars[i%len(lChars)]
			entry.WriteString(string(char))
		}
		entry.WriteString(" ")

		// Format: symbol label percentage [value]
		if p.opts.ShowValues {
			entry.WriteString(fmt.Sprintf("%-8s %5.1f%% [%.1f]", slice.Label, slice.Percentage, slice.Value))
		} else {
			entry.WriteString(fmt.Sprintf("%-8s %5.1f%%", slice.Label, slice.Percentage))
		}

		legendEntries[i] = entry.String()
	}

	// Combine pie chart and legend side by side
	var result strings.Builder
	legendStartRow := (len(pieRows) - len(legendEntries)) / 2
	if legendStartRow < 0 {
		legendStartRow = 0
	}

	for i, pieRow := range pieRows {
		result.WriteString(pieRow)
		result.WriteString("   ") // Gap between pie and legend

		// Add legend entry if available for this row
		legendIdx := i - legendStartRow
		if legendIdx >= 0 && legendIdx < len(legendEntries) {
			result.WriteString(legendEntries[legendIdx])
		}

		result.WriteString("\n")
	}

	return result.String()
}

// findSliceForAngle finds which slice a given angle belongs to.
func (p *PieChart) findSliceForAngle(angle float64, angles []float64, numSlices int) int {
	for i := 0; i < numSlices; i++ {
		startAngle := angles[i]
		endAngle := angles[i+1]

		// Handle the wrap-around case
		if endAngle > math.Pi {
			// Slice crosses the -π/π boundary
			if angle >= startAngle || angle < endAngle-2*math.Pi {
				return i
			}
		} else if angle >= startAngle && angle < endAngle {
			return i
		}
	}

	// Default to last slice
	return numSlices - 1
}

// shouldUseUnicode determines whether to use Unicode characters based on style.
func (p *PieChart) shouldUseUnicode() bool {
	if p.opts.Style == StyleASCII {
		return false
	} else if p.opts.Style == StyleUnicode {
		return true
	}
	// StyleAuto - detect Unicode support
	return internal.SupportsUnicode()
}

// isColorEnabled determines whether colors should be used.
func (p *PieChart) isColorEnabled() bool {
	if p.opts.ColorEnabled != nil {
		return *p.opts.ColorEnabled
	}
	return internal.SupportsColor()
}

// Pie is a convenience function that creates and renders a pie chart.
//
// Example:
//
//	fmt.Println(termcharts.Pie([]float64{30, 25, 20, 15, 10}))
func Pie(data []float64) string {
	pie := NewPieChart(WithData(data))
	return pie.Render()
}

// PieWithLabels is a convenience function that creates a pie chart with labels.
//
// Example:
//
//	fmt.Println(termcharts.PieWithLabels(
//	    []float64{30, 25, 20, 15, 10},
//	    []string{"Chrome", "Firefox", "Safari", "Edge", "Other"},
//	))
func PieWithLabels(data []float64, labels []string) string {
	pie := NewPieChart(
		WithData(data),
		WithLabels(labels),
	)
	return pie.Render()
}

// PieWithValues is a convenience function that creates a pie chart with values displayed.
//
// Example:
//
//	fmt.Println(termcharts.PieWithValues(
//	    []float64{30, 25, 20, 15, 10},
//	    []string{"Chrome", "Firefox", "Safari", "Edge", "Other"},
//	))
func PieWithValues(data []float64, labels []string) string {
	pie := NewPieChart(
		WithData(data),
		WithLabels(labels),
		WithShowValues(true),
	)
	return pie.Render()
}
