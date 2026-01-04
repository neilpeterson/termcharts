package termcharts

import (
	"fmt"
	"math"
	"strings"

	"github.com/neilpeterson/termcharts/internal"
)

// PieChart represents a pie chart visualization.
// Pie charts display data as proportional slices of a circle,
// rendered using Unicode or ASCII characters in the terminal.
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

// Unicode characters for pie chart rendering
const (
	pieBlockFull  = '█'
	pieBlockASCII = '#'
)

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

	// Render circular pie visualization
	pieVis := p.renderCircularPie(slices, colorEnabled, theme)
	result.WriteString(pieVis)

	// Render legend
	legend := p.renderLegend(slices, colorEnabled, theme)
	result.WriteString(legend)

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

// renderCircularPie renders a circular pie chart using ASCII/Unicode art.
func (p *PieChart) renderCircularPie(slices []Slice, colorEnabled bool, theme *Theme) string {
	var result strings.Builder

	// Determine character set
	useUnicode := p.shouldUseUnicode()

	// Calculate radius based on available space
	// Terminal chars are roughly 2x taller than wide, so we adjust
	radius := 8
	if p.opts.Height > 0 && p.opts.Height < 20 {
		radius = p.opts.Height / 2
		if radius < 4 {
			radius = 4
		}
	}

	// Character aspect ratio compensation (chars are ~2x taller than wide)
	aspectRatio := 2.0

	// Calculate cumulative angles for each slice
	angles := make([]float64, len(slices)+1)
	angles[0] = -math.Pi / 2 // Start at top (12 o'clock)
	for i, slice := range slices {
		angles[i+1] = angles[i] + (slice.Percentage/100)*2*math.Pi
	}

	// Render the circle row by row
	for y := -radius; y <= radius; y++ {
		result.WriteString("  ") // Left padding
		for x := -radius * int(aspectRatio); x <= radius*int(aspectRatio); x++ {
			// Calculate actual position accounting for aspect ratio
			actualX := float64(x) / aspectRatio
			actualY := float64(y)

			// Calculate distance from center
			distance := math.Sqrt(actualX*actualX + actualY*actualY)

			// Check if point is inside the circle
			if distance <= float64(radius) {
				// Calculate angle of this point
				angle := math.Atan2(actualY, actualX)

				// Find which slice this angle belongs to
				sliceIndex := 0
				for i := 0; i < len(slices); i++ {
					if angle >= angles[i] && angle < angles[i+1] {
						sliceIndex = i
						break
					}
					// Handle wrap-around at the top
					if angles[i+1] > math.Pi && angle < angles[0] {
						// Adjust angle for comparison
						adjustedAngle := angle + 2*math.Pi
						if adjustedAngle >= angles[i] && adjustedAngle < angles[i+1] {
							sliceIndex = i
							break
						}
					}
				}

				// Get character and color for this slice
				char := string(pieBlockFull)
				if !useUnicode {
					char = string(pieBlockASCII)
				}

				color := theme.GetSeriesColor(sliceIndex)
				if colorEnabled {
					result.WriteString(Colorize(char, color, true))
				} else {
					// In non-color mode, use different characters for each slice
					chars := []rune{'█', '▓', '▒', '░', '▪', '▫'}
					if !useUnicode {
						chars = []rune{'#', '*', '+', 'o', 'x', '.'}
					}
					result.WriteString(string(chars[sliceIndex%len(chars)]))
				}
			} else {
				result.WriteString(" ")
			}
		}
		result.WriteString("\n")
	}

	result.WriteString("\n")
	return result.String()
}

// renderLegend renders the pie chart legend with labels, values, and percentages.
func (p *PieChart) renderLegend(slices []Slice, colorEnabled bool, theme *Theme) string {
	var result strings.Builder

	// Determine character set
	useUnicode := p.shouldUseUnicode()

	// Characters for non-color mode legend
	legendChars := []rune{'█', '▓', '▒', '░', '▪', '▫'}
	if !useUnicode {
		legendChars = []rune{'#', '*', '+', 'o', 'x', '.'}
	}

	// Find max label width for alignment
	maxLabelWidth := 0
	for _, slice := range slices {
		if len(slice.Label) > maxLabelWidth {
			maxLabelWidth = len(slice.Label)
		}
	}

	// Render each legend entry
	for i, slice := range slices {
		// Color indicator
		color := theme.GetSeriesColor(i)
		indicator := string(pieBlockFull)
		if !useUnicode {
			indicator = string(pieBlockASCII)
		}

		if colorEnabled {
			indicator = Colorize(indicator, color, true)
		} else {
			indicator = string(legendChars[i%len(legendChars)])
		}

		result.WriteString("  ")
		result.WriteString(indicator)
		result.WriteString(" ")

		// Label
		labelText := fmt.Sprintf("%-*s", maxLabelWidth, slice.Label)
		if colorEnabled {
			labelText = Colorize(labelText, theme.Text, true)
		}
		result.WriteString(labelText)

		// Value and percentage
		if p.opts.ShowValues {
			valueText := fmt.Sprintf("  %6.1f", slice.Value)
			if colorEnabled {
				valueText = Colorize(valueText, theme.Muted, true)
			}
			result.WriteString(valueText)
		}

		percentText := fmt.Sprintf("  (%5.1f%%)", slice.Percentage)
		if colorEnabled {
			percentText = Colorize(percentText, theme.Muted, true)
		}
		result.WriteString(percentText)

		result.WriteString("\n")
	}

	return result.String()
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
