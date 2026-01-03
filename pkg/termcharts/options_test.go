package termcharts

import (
	"testing"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()

	if opts.Width != 80 {
		t.Errorf("default Width = %v, want %v", opts.Width, 80)
	}

	if opts.Height != 24 {
		t.Errorf("default Height = %v, want %v", opts.Height, 24)
	}

	if opts.Style != StyleAuto {
		t.Errorf("default Style = %v, want %v", opts.Style, StyleAuto)
	}

	if opts.Direction != Horizontal {
		t.Errorf("default Direction = %v, want %v", opts.Direction, Horizontal)
	}

	if opts.ShowValues != false {
		t.Errorf("default ShowValues = %v, want %v", opts.ShowValues, false)
	}

	if opts.ShowAxes != true {
		t.Errorf("default ShowAxes = %v, want %v", opts.ShowAxes, true)
	}
}

func TestWithData(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	opts := NewOptions(WithData(data))

	if len(opts.Data) != 5 {
		t.Errorf("len(Data) = %v, want %v", len(opts.Data), 5)
	}

	for i, v := range data {
		if opts.Data[i] != v {
			t.Errorf("Data[%d] = %v, want %v", i, opts.Data[i], v)
		}
	}
}

func TestWithLabels(t *testing.T) {
	labels := []string{"A", "B", "C"}
	opts := NewOptions(WithLabels(labels))

	if len(opts.Labels) != 3 {
		t.Errorf("len(Labels) = %v, want %v", len(opts.Labels), 3)
	}

	for i, v := range labels {
		if opts.Labels[i] != v {
			t.Errorf("Labels[%d] = %v, want %v", i, opts.Labels[i], v)
		}
	}
}

func TestWithSeries(t *testing.T) {
	series := []Series{
		{Label: "Series 1", Data: []float64{1, 2, 3}},
		{Label: "Series 2", Data: []float64{4, 5, 6}},
	}
	opts := NewOptions(WithSeries(series))

	if len(opts.Series) != 2 {
		t.Errorf("len(Series) = %v, want %v", len(opts.Series), 2)
	}

	if opts.Series[0].Label != "Series 1" {
		t.Errorf("Series[0].Label = %v, want %v", opts.Series[0].Label, "Series 1")
	}
}

func TestWithWidth(t *testing.T) {
	opts := NewOptions(WithWidth(120))

	if opts.Width != 120 {
		t.Errorf("Width = %v, want %v", opts.Width, 120)
	}
}

func TestWithHeight(t *testing.T) {
	opts := NewOptions(WithHeight(30))

	if opts.Height != 30 {
		t.Errorf("Height = %v, want %v", opts.Height, 30)
	}
}

func TestWithTitle(t *testing.T) {
	opts := NewOptions(WithTitle("Test Chart"))

	if opts.Title != "Test Chart" {
		t.Errorf("Title = %v, want %v", opts.Title, "Test Chart")
	}
}

func TestWithColor(t *testing.T) {
	tests := []struct {
		name     string
		enabled  bool
		expected bool
	}{
		{
			name:     "color enabled",
			enabled:  true,
			expected: true,
		},
		{
			name:     "color disabled",
			enabled:  false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions(WithColor(tt.enabled))

			if opts.ColorEnabled == nil {
				t.Fatal("ColorEnabled is nil")
			}

			if *opts.ColorEnabled != tt.expected {
				t.Errorf("ColorEnabled = %v, want %v", *opts.ColorEnabled, tt.expected)
			}
		})
	}
}

func TestWithStyle(t *testing.T) {
	tests := []struct {
		name     string
		style    RenderStyle
		expected RenderStyle
	}{
		{
			name:     "StyleASCII",
			style:    StyleASCII,
			expected: StyleASCII,
		},
		{
			name:     "StyleUnicode",
			style:    StyleUnicode,
			expected: StyleUnicode,
		},
		{
			name:     "StyleBraille",
			style:    StyleBraille,
			expected: StyleBraille,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions(WithStyle(tt.style))

			if opts.Style != tt.expected {
				t.Errorf("Style = %v, want %v", opts.Style, tt.expected)
			}
		})
	}
}

func TestWithDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction Direction
		expected  Direction
	}{
		{
			name:      "Horizontal",
			direction: Horizontal,
			expected:  Horizontal,
		},
		{
			name:      "Vertical",
			direction: Vertical,
			expected:  Vertical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions(WithDirection(tt.direction))

			if opts.Direction != tt.expected {
				t.Errorf("Direction = %v, want %v", opts.Direction, tt.expected)
			}
		})
	}
}

func TestWithShowValues(t *testing.T) {
	tests := []struct {
		name     string
		show     bool
		expected bool
	}{
		{
			name:     "show values true",
			show:     true,
			expected: true,
		},
		{
			name:     "show values false",
			show:     false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions(WithShowValues(tt.show))

			if opts.ShowValues != tt.expected {
				t.Errorf("ShowValues = %v, want %v", opts.ShowValues, tt.expected)
			}
		})
	}
}

func TestWithShowAxes(t *testing.T) {
	tests := []struct {
		name     string
		show     bool
		expected bool
	}{
		{
			name:     "show axes true",
			show:     true,
			expected: true,
		},
		{
			name:     "show axes false",
			show:     false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions(WithShowAxes(tt.show))

			if opts.ShowAxes != tt.expected {
				t.Errorf("ShowAxes = %v, want %v", opts.ShowAxes, tt.expected)
			}
		})
	}
}

func TestWithTheme(t *testing.T) {
	theme := &Theme{
		Primary:   "red",
		Secondary: "blue",
	}
	opts := NewOptions(WithTheme(theme))

	if opts.Theme == nil {
		t.Fatal("Theme is nil")
	}

	if opts.Theme.Primary != "red" {
		t.Errorf("Theme.Primary = %v, want %v", opts.Theme.Primary, "red")
	}

	if opts.Theme.Secondary != "blue" {
		t.Errorf("Theme.Secondary = %v, want %v", opts.Theme.Secondary, "blue")
	}
}

func TestMultipleOptions(t *testing.T) {
	opts := NewOptions(
		WithData([]float64{1, 2, 3}),
		WithWidth(100),
		WithHeight(20),
		WithTitle("Multi-Option Test"),
		WithStyle(StyleUnicode),
		WithDirection(Vertical),
		WithShowValues(true),
	)

	if len(opts.Data) != 3 {
		t.Errorf("len(Data) = %v, want %v", len(opts.Data), 3)
	}

	if opts.Width != 100 {
		t.Errorf("Width = %v, want %v", opts.Width, 100)
	}

	if opts.Height != 20 {
		t.Errorf("Height = %v, want %v", opts.Height, 20)
	}

	if opts.Title != "Multi-Option Test" {
		t.Errorf("Title = %v, want %v", opts.Title, "Multi-Option Test")
	}

	if opts.Style != StyleUnicode {
		t.Errorf("Style = %v, want %v", opts.Style, StyleUnicode)
	}

	if opts.Direction != Vertical {
		t.Errorf("Direction = %v, want %v", opts.Direction, Vertical)
	}

	if opts.ShowValues != true {
		t.Errorf("ShowValues = %v, want %v", opts.ShowValues, true)
	}
}
