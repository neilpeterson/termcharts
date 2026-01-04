package termcharts

import (
	"math"
	"strings"
	"testing"
)

func TestNewLineChart(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	line := NewLineChart(WithData(data))

	if line == nil {
		t.Fatal("NewLineChart returned nil")
	}
	if line.opts == nil {
		t.Fatal("Options not initialized")
	}
	if len(line.opts.Data) != len(data) {
		t.Errorf("Expected data length %d, got %d", len(data), len(line.opts.Data))
	}
}

func TestLineChart_Render_BasicData(t *testing.T) {
	tests := []struct {
		name string
		data []float64
	}{
		{
			name: "simple ascending",
			data: []float64{1, 2, 3, 4, 5},
		},
		{
			name: "simple descending",
			data: []float64{5, 4, 3, 2, 1},
		},
		{
			name: "mixed values",
			data: []float64{1, 5, 2, 8, 3, 7},
		},
		{
			name: "negative values",
			data: []float64{-5, -2, 0, 3, 5},
		},
		{
			name: "fractional values",
			data: []float64{0.1, 0.5, 0.2, 0.8, 0.3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := NewLineChart(
				WithData(tt.data),
				WithHeight(10),
				WithWidth(40),
			)
			result := line.Render()

			// Should not be empty
			if result == "" {
				t.Error("Render returned empty string")
			}

			// Should contain multiple lines
			lines := strings.Split(result, "\n")
			if len(lines) < 2 {
				t.Error("Expected multi-line output")
			}
		})
	}
}

func TestLineChart_Render_EmptyData(t *testing.T) {
	line := NewLineChart(WithData([]float64{}))
	result := line.Render()

	if result != "" {
		t.Errorf("Expected empty string for empty data, got: %s", result)
	}
}

func TestLineChart_Render_SingleValue(t *testing.T) {
	line := NewLineChart(
		WithData([]float64{42}),
		WithHeight(5),
		WithWidth(20),
	)
	result := line.Render()

	// Should still render something
	if result == "" {
		t.Error("Expected non-empty output for single value")
	}
}

func TestLineChart_Render_AllSameValues(t *testing.T) {
	line := NewLineChart(
		WithData([]float64{5, 5, 5, 5, 5}),
		WithHeight(10),
		WithWidth(40),
	)
	result := line.Render()

	// Should render a horizontal line
	if result == "" {
		t.Error("Expected non-empty output")
	}
}

func TestLineChart_Render_ASCIIMode(t *testing.T) {
	data := []float64{1, 5, 2, 8, 3, 7}
	line := NewLineChart(
		WithData(data),
		WithStyle(StyleASCII),
		WithHeight(8),
		WithWidth(40),
		WithShowAxes(false),
	)
	result := line.Render()

	// Check that only ASCII characters are used (no Unicode box-drawing)
	for _, r := range result {
		if r > 127 && r != '\n' {
			t.Errorf("Expected ASCII only, got Unicode character: %c (U+%04X)", r, r)
		}
	}
}

func TestLineChart_Render_UnicodeMode(t *testing.T) {
	data := []float64{1, 5, 2, 8, 3, 7}
	line := NewLineChart(
		WithData(data),
		WithStyle(StyleUnicode),
		WithHeight(8),
		WithWidth(40),
	)
	result := line.Render()

	// Should contain Unicode characters
	hasUnicode := false
	for _, r := range result {
		if r > 127 && r != '\n' {
			hasUnicode = true
			break
		}
	}
	if !hasUnicode {
		t.Error("Expected Unicode characters in output")
	}
}

func TestLineChart_Render_BrailleMode(t *testing.T) {
	data := []float64{1, 5, 2, 8, 3, 7}
	line := NewLineChart(
		WithData(data),
		WithStyle(StyleBraille),
		WithHeight(8),
		WithWidth(40),
	)
	result := line.Render()

	// Should contain Braille patterns (U+2800 to U+28FF)
	hasBraille := false
	for _, r := range result {
		if r >= 0x2800 && r <= 0x28FF {
			hasBraille = true
			break
		}
	}
	if !hasBraille {
		t.Error("Expected Braille patterns in output")
	}
}

func TestLineChart_Render_InvalidData(t *testing.T) {
	tests := []struct {
		name string
		data []float64
	}{
		{
			name: "contains NaN",
			data: []float64{1, 2, math.NaN(), 4, 5},
		},
		{
			name: "contains positive infinity",
			data: []float64{1, 2, math.Inf(1), 4, 5},
		},
		{
			name: "contains negative infinity",
			data: []float64{1, 2, math.Inf(-1), 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := NewLineChart(WithData(tt.data))
			result := line.Render()

			// Should return empty string for invalid data
			if result != "" {
				t.Errorf("Expected empty string for invalid data, got non-empty output")
			}
		})
	}
}

func TestLineChart_Render_WithTitle(t *testing.T) {
	title := "Test Chart"
	line := NewLineChart(
		WithData([]float64{1, 2, 3, 4, 5}),
		WithTitle(title),
		WithHeight(8),
		WithWidth(40),
	)
	result := line.Render()

	if !strings.Contains(result, title) {
		t.Errorf("Expected title '%s' in output", title)
	}
}

func TestLineChart_Render_WithLabels(t *testing.T) {
	labels := []string{"A", "B", "C", "D", "E"}
	line := NewLineChart(
		WithData([]float64{1, 2, 3, 4, 5}),
		WithLabels(labels),
		WithShowAxes(true),
		WithHeight(8),
		WithWidth(40),
	)
	result := line.Render()

	// At least some labels should appear
	foundLabels := 0
	for _, label := range labels {
		if strings.Contains(result, label) {
			foundLabels++
		}
	}
	if foundLabels == 0 {
		t.Error("Expected at least some labels to appear in output")
	}
}

func TestLineChart_Render_WithColor(t *testing.T) {
	colorEnabled := true
	line := NewLineChart(
		WithData([]float64{1, 5, 2, 8, 3}),
		WithColor(colorEnabled),
		WithHeight(8),
		WithWidth(40),
	)
	result := line.Render()

	// Should contain ANSI color codes
	if !strings.Contains(result, "\033[") {
		t.Error("Expected ANSI color codes in output")
	}
}

func TestLineChart_Render_WithoutAxes(t *testing.T) {
	line := NewLineChart(
		WithData([]float64{1, 2, 3, 4, 5}),
		WithShowAxes(false),
		WithHeight(8),
		WithWidth(40),
	)
	result := line.Render()

	// Should not contain axis line
	if strings.Contains(result, "────────") {
		t.Error("Expected no axis line when axes disabled")
	}
}

func TestLineChart_Render_MultiSeries(t *testing.T) {
	series := []Series{
		{Label: "Series A", Data: []float64{1, 3, 2, 4, 3}},
		{Label: "Series B", Data: []float64{2, 1, 3, 2, 4}},
	}
	line := NewLineChart(
		WithSeries(series),
		WithHeight(10),
		WithWidth(50),
	)
	result := line.Render()

	// Should render something
	if result == "" {
		t.Error("Expected non-empty output for multi-series")
	}

	// Should contain legend with series labels
	if !strings.Contains(result, "Series A") {
		t.Error("Expected 'Series A' in legend")
	}
	if !strings.Contains(result, "Series B") {
		t.Error("Expected 'Series B' in legend")
	}
}

func TestLineChart_Render_MultiSeriesWithColor(t *testing.T) {
	series := []Series{
		{Label: "Sales", Data: []float64{10, 20, 15, 25}, Color: "red"},
		{Label: "Costs", Data: []float64{8, 12, 10, 18}, Color: "blue"},
	}
	colorEnabled := true
	line := NewLineChart(
		WithSeries(series),
		WithColor(colorEnabled),
		WithHeight(10),
		WithWidth(50),
	)
	result := line.Render()

	// Should contain ANSI color codes
	if !strings.Contains(result, "\033[") {
		t.Error("Expected ANSI color codes in output")
	}
}

func TestLineChart_ConvenienceFunction(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	result := Line(data)

	if result == "" {
		t.Error("Line() returned empty string")
	}
}

func TestLineBraille_ConvenienceFunction(t *testing.T) {
	data := []float64{1, 5, 2, 8, 3, 7}
	result := LineBraille(data)

	if result == "" {
		t.Error("LineBraille() returned empty string")
	}

	// Should contain Braille patterns
	hasBraille := false
	for _, r := range result {
		if r >= 0x2800 && r <= 0x28FF {
			hasBraille = true
			break
		}
	}
	if !hasBraille {
		t.Error("Expected Braille patterns in output")
	}
}

func TestLineMultiSeries_ConvenienceFunction(t *testing.T) {
	series := []Series{
		{Label: "A", Data: []float64{1, 2, 3}},
		{Label: "B", Data: []float64{3, 2, 1}},
	}
	result := LineMultiSeries(series)

	if result == "" {
		t.Error("LineMultiSeries() returned empty string")
	}
}

func TestLineChart_Render_Dimensions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small", 20, 5},
		{"medium", 60, 15},
		{"large", 100, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := NewLineChart(
				WithData([]float64{1, 5, 2, 8, 3, 7, 4, 6}),
				WithWidth(tt.width),
				WithHeight(tt.height),
				WithShowAxes(false),
			)
			result := line.Render()

			lines := strings.Split(strings.TrimSuffix(result, "\n"), "\n")

			// Height should roughly match (may vary slightly due to rendering)
			if len(lines) < tt.height-2 || len(lines) > tt.height+2 {
				t.Errorf("Expected approximately %d lines, got %d", tt.height, len(lines))
			}
		})
	}
}

func TestLineChart_Render_Theme(t *testing.T) {
	tests := []struct {
		name  string
		theme *Theme
	}{
		{"default", DefaultTheme},
		{"dark", DarkTheme},
		{"light", LightTheme},
		{"mono", MonochromeTheme},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorEnabled := true
			line := NewLineChart(
				WithData([]float64{1, 2, 3, 4, 5}),
				WithTheme(tt.theme),
				WithColor(colorEnabled),
				WithHeight(8),
				WithWidth(40),
			)
			result := line.Render()

			if result == "" {
				t.Error("Expected non-empty output")
			}
		})
	}
}

func TestLineChart_findGlobalMinMax(t *testing.T) {
	line := NewLineChart()

	tests := []struct {
		name        string
		series      []Series
		expectedMin float64
		expectedMax float64
	}{
		{
			name: "single series",
			series: []Series{
				{Data: []float64{1, 5, 3}},
			},
			expectedMin: 1,
			expectedMax: 5,
		},
		{
			name: "multiple series",
			series: []Series{
				{Data: []float64{2, 4}},
				{Data: []float64{1, 6}},
			},
			expectedMin: 1,
			expectedMax: 6,
		},
		{
			name: "negative values",
			series: []Series{
				{Data: []float64{-5, 0, 5}},
			},
			expectedMin: -5,
			expectedMax: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := line.findGlobalMinMax(tt.series)
			if min != tt.expectedMin {
				t.Errorf("Expected min %v, got %v", tt.expectedMin, min)
			}
			if max != tt.expectedMax {
				t.Errorf("Expected max %v, got %v", tt.expectedMax, max)
			}
		})
	}
}

func TestLineChart_getAllSeries(t *testing.T) {
	t.Run("from Data", func(t *testing.T) {
		line := NewLineChart(WithData([]float64{1, 2, 3}))
		series := line.getAllSeries()
		if len(series) != 1 {
			t.Errorf("Expected 1 series, got %d", len(series))
		}
		if len(series[0].Data) != 3 {
			t.Errorf("Expected 3 data points, got %d", len(series[0].Data))
		}
	})

	t.Run("from Series", func(t *testing.T) {
		s := []Series{
			{Label: "A", Data: []float64{1, 2}},
			{Label: "B", Data: []float64{3, 4}},
		}
		line := NewLineChart(WithSeries(s))
		series := line.getAllSeries()
		if len(series) != 2 {
			t.Errorf("Expected 2 series, got %d", len(series))
		}
	})

	t.Run("empty", func(t *testing.T) {
		line := NewLineChart()
		series := line.getAllSeries()
		if series != nil {
			t.Errorf("Expected nil, got %v", series)
		}
	})
}
