// Last modified: 2026-01-02

package termcharts

// Options holds configuration for chart rendering.
// Options are set using functional options via With* functions.
type Options struct {
	// Width is the maximum chart width in terminal columns (0 = auto-detect).
	Width int
	// Height is the maximum chart height in terminal rows (0 = auto-detect).
	Height int
	// Title is an optional chart title displayed above the chart.
	Title string
	// Data contains the primary data series to visualize.
	Data []float64
	// Labels contains optional labels for each data point.
	Labels []string
	// Series contains multiple data series for multi-series charts.
	Series []Series
	// ColorEnabled controls whether to use ANSI colors (auto-detected if nil).
	ColorEnabled *bool
	// Style specifies the rendering mode (ASCII, Unicode, or Braille).
	Style RenderStyle
	// Direction specifies chart orientation (horizontal or vertical).
	Direction Direction
	// ShowValues controls whether to display numeric values on the chart.
	ShowValues bool
	// ShowAxes controls whether to display axes and labels.
	ShowAxes bool
	// Theme specifies the color theme to use.
	Theme *Theme
}

// Option is a function that configures chart Options using the functional options pattern.
type Option func(*Options)

// NewOptions creates a new Options struct with sensible defaults.
func NewOptions(opts ...Option) *Options {
	// Default options
	o := &Options{
		Width:      80,  // Standard terminal width
		Height:     24,  // Standard terminal height
		Style:      StyleAuto,
		Direction:  Horizontal,
		ShowValues: false,
		ShowAxes:   true,
	}

	// Apply all provided options
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithData sets the primary data series for the chart.
func WithData(data []float64) Option {
	return func(o *Options) {
		o.Data = data
	}
}

// WithLabels sets labels for each data point.
// The number of labels should match the number of data points.
func WithLabels(labels []string) Option {
	return func(o *Options) {
		o.Labels = labels
	}
}

// WithSeries sets multiple data series for multi-series charts.
func WithSeries(series []Series) Option {
	return func(o *Options) {
		o.Series = series
	}
}

// WithWidth sets the maximum chart width in terminal columns.
// Use 0 to auto-detect terminal width.
func WithWidth(width int) Option {
	return func(o *Options) {
		o.Width = width
	}
}

// WithHeight sets the maximum chart height in terminal rows.
// Use 0 to auto-detect terminal height.
func WithHeight(height int) Option {
	return func(o *Options) {
		o.Height = height
	}
}

// WithTitle sets an optional title displayed above the chart.
func WithTitle(title string) Option {
	return func(o *Options) {
		o.Title = title
	}
}

// WithColor enables or disables ANSI color output.
// If not set, color support is auto-detected.
func WithColor(enabled bool) Option {
	return func(o *Options) {
		o.ColorEnabled = &enabled
	}
}

// WithStyle sets the rendering style (ASCII, Unicode, or Braille).
// StyleAuto automatically selects the best style based on terminal capabilities.
func WithStyle(style RenderStyle) Option {
	return func(o *Options) {
		o.Style = style
	}
}

// WithDirection sets the chart orientation (horizontal or vertical).
func WithDirection(dir Direction) Option {
	return func(o *Options) {
		o.Direction = dir
	}
}

// WithShowValues controls whether numeric values are displayed on the chart.
func WithShowValues(show bool) Option {
	return func(o *Options) {
		o.ShowValues = show
	}
}

// WithShowAxes controls whether axes and labels are displayed.
func WithShowAxes(show bool) Option {
	return func(o *Options) {
		o.ShowAxes = show
	}
}

// WithTheme sets the color theme for the chart.
func WithTheme(theme *Theme) Option {
	return func(o *Options) {
		o.Theme = theme
	}
}
