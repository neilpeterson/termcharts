// Last modified: 2026-01-02

package termcharts

import (
	"strings"

	"github.com/neilpeterson/termcharts/internal"
)

// Sparkline represents a compact, inline chart showing data trends.
// Sparklines use Unicode block characters (▁▂▃▄▅▆▇█) to visualize
// data in a single line, perfect for dashboards and monitoring.
type Sparkline struct {
	opts *Options
}

// Unicode block characters for sparkline rendering (8 levels).
var sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// ASCII characters for sparkline rendering when Unicode is not supported.
var sparkCharsASCII = []rune{'_', '.', '-', '=', '+', '*', '#', '@'}

// NewSparkline creates a new sparkline chart with the given options.
// At minimum, data must be provided via WithData option.
//
// Example:
//
//	spark := termcharts.NewSparkline(
//	    termcharts.WithData([]float64{1, 5, 2, 8, 3, 7}),
//	)
//	fmt.Println(spark.Render())
func NewSparkline(opts ...Option) *Sparkline {
	options := NewOptions(opts...)
	return &Sparkline{
		opts: options,
	}
}

// Render generates the sparkline as a single-line string.
// Each data point is represented by a single character, with height
// proportional to the value relative to the min/max in the dataset.
func (s *Sparkline) Render() string {
	// Validate data
	if len(s.opts.Data) == 0 {
		return ""
	}

	// Check for invalid values
	if !internal.AllValid(s.opts.Data) {
		return ""
	}

	// Determine character set based on style
	chars := sparkChars
	if s.opts.Style == StyleASCII {
		chars = sparkCharsASCII
	} else if s.opts.Style == StyleAuto {
		// Auto-detect Unicode support
		if !internal.SupportsUnicode() {
			chars = sparkCharsASCII
		}
	}

	// Normalize data to 0-1 range
	normalized, _, _ := internal.Normalize(s.opts.Data)

	var result strings.Builder

	// Apply width limit if specified
	data := normalized
	if s.opts.Width > 0 && len(normalized) > s.opts.Width {
		// Sample data to fit width
		data = sampleData(normalized, s.opts.Width)
	}

	// Map each value to a character
	for _, val := range data {
		// Map 0-1 to character index (0-7)
		level := int(val * float64(len(chars)-1))
		if level < 0 {
			level = 0
		}
		if level >= len(chars) {
			level = len(chars) - 1
		}

		char := chars[level]

		// Apply color if enabled
		if s.opts.ColorEnabled != nil && *s.opts.ColorEnabled {
			color := s.getColorForLevel(level, len(chars))
			result.WriteString(Colorize(string(char), color, true))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// getColorForLevel returns a color based on the value level.
// Lower values are blue/green, higher values are yellow/red.
func (s *Sparkline) getColorForLevel(level, maxLevel int) string {
	theme := s.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Map level to color based on intensity
	ratio := float64(level) / float64(maxLevel-1)
	if ratio < 0.33 {
		return theme.Muted // Low values
	} else if ratio < 0.66 {
		return theme.Primary // Medium values
	} else {
		return theme.Accent // High values
	}
}

// sampleData reduces the data to the target width by sampling.
// Uses simple downsampling - takes every Nth value.
func sampleData(data []float64, targetWidth int) []float64 {
	if len(data) <= targetWidth {
		return data
	}

	result := make([]float64, targetWidth)
	step := float64(len(data)) / float64(targetWidth)

	for i := 0; i < targetWidth; i++ {
		index := int(float64(i) * step)
		if index >= len(data) {
			index = len(data) - 1
		}
		result[i] = data[index]
	}

	return result
}

// Spark is a convenience function that creates and renders a sparkline in one call.
// This is the simplest way to generate a sparkline from data.
//
// Example:
//
//	fmt.Println(termcharts.Spark([]float64{1, 5, 2, 8, 3, 7, 4, 6}))
//	// Output: ▁▅▂█▃▇▄▆
func Spark(data []float64) string {
	spark := NewSparkline(WithData(data))
	return spark.Render()
}

// SparkASCII is a convenience function that creates an ASCII-only sparkline.
// Use this when you need maximum compatibility with terminals that don't support Unicode.
//
// Example:
//
//	fmt.Println(termcharts.SparkASCII([]float64{1, 5, 2, 8, 3, 7, 4, 6}))
//	// Output: _*-@+#.=
func SparkASCII(data []float64) string {
	spark := NewSparkline(
		WithData(data),
		WithStyle(StyleASCII),
	)
	return spark.Render()
}

// SparkColor creates a colored sparkline with auto-detected color support.
//
// Example:
//
//	fmt.Println(termcharts.SparkColor([]float64{1, 5, 2, 8, 3, 7, 4, 6}))
func SparkColor(data []float64) string {
	colorEnabled := internal.SupportsColor()
	spark := NewSparkline(
		WithData(data),
		WithColor(colorEnabled),
	)
	return spark.Render()
}
