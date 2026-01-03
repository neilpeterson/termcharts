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
//
// The library auto-detects terminal capabilities and adjusts rendering accordingly.
package termcharts

import "errors"

// Chart represents a terminal-based data visualization.
// All chart types implement this interface.
type Chart interface {
	// Render generates the chart as a string ready for terminal output.
	// The output includes newlines for multi-line charts.
	Render() string
}

// Series represents a labeled data series for multi-series charts.
type Series struct {
	// Label is the display name for this data series.
	Label string
	// Data contains the numeric values to visualize.
	Data []float64
	// Color is an optional color for this series (empty means auto-assign).
	Color string
}

// Direction specifies the orientation of a chart.
type Direction int

const (
	// Horizontal renders charts from left to right.
	Horizontal Direction = iota
	// Vertical renders charts from bottom to top.
	Vertical
)

// String returns the string representation of the Direction.
func (d Direction) String() string {
	switch d {
	case Horizontal:
		return "horizontal"
	case Vertical:
		return "vertical"
	default:
		return "unknown"
	}
}

// Common errors returned by the library.
var (
	// ErrEmptyData indicates no data was provided for visualization.
	ErrEmptyData = errors.New("data cannot be empty")
	// ErrInvalidData indicates the data contains invalid values (NaN, Inf, etc.).
	ErrInvalidData = errors.New("data contains invalid values")
	// ErrInvalidDimensions indicates chart dimensions are too small to render.
	ErrInvalidDimensions = errors.New("chart dimensions too small")
)
