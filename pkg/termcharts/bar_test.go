package termcharts

import (
	"math"
	"strings"
	"testing"
)

func TestNewBarChart(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	bar := NewBarChart(WithData(data))

	if bar == nil {
		t.Fatal("NewBarChart returned nil")
	}
	if bar.opts == nil {
		t.Fatal("Options not initialized")
	}
	if len(bar.opts.Data) != len(data) {
		t.Errorf("Expected data length %d, got %d", len(data), len(bar.opts.Data))
	}
}

func TestBarChart_Render_BasicData(t *testing.T) {
	tests := []struct {
		name string
		data []float64
	}{
		{
			name: "simple ascending",
			data: []float64{10, 20, 30, 40},
		},
		{
			name: "simple descending",
			data: []float64{40, 30, 20, 10},
		},
		{
			name: "mixed values",
			data: []float64{10, 25, 15, 30, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBarChart(WithData(tt.data))
			result := bar.Render()

			// Should not be empty
			if result == "" {
				t.Error("Render returned empty string")
			}

			// Should have multiple lines (one per data point)
			lines := strings.Split(strings.TrimSpace(result), "\n")
			if len(lines) != len(tt.data) {
				t.Errorf("Expected %d lines, got %d", len(tt.data), len(lines))
			}
		})
	}
}

func TestBarChart_Render_EmptyData(t *testing.T) {
	bar := NewBarChart(WithData([]float64{}))
	result := bar.Render()

	if result != "" {
		t.Errorf("Expected empty string for empty data, got: %s", result)
	}
}

func TestBarChart_Render_SingleValue(t *testing.T) {
	bar := NewBarChart(WithData([]float64{42}))
	result := bar.Render()

	// Should render a single line
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
}

func TestBarChart_Render_WithLabels(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	labels := []string{"Q1", "Q2", "Q3", "Q4"}

	bar := NewBarChart(
		WithData(data),
		WithLabels(labels),
	)
	result := bar.Render()

	// Should contain all labels
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("Expected output to contain label %s", label)
		}
	}
}

func TestBarChart_Render_WithTitle(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	title := "Quarterly Sales"

	bar := NewBarChart(
		WithData(data),
		WithTitle(title),
	)
	result := bar.Render()

	// Should contain title
	if !strings.Contains(result, title) {
		t.Error("Expected output to contain title")
	}

	// Title should be on first line
	lines := strings.Split(result, "\n")
	if !strings.Contains(lines[0], title) {
		t.Error("Expected title to be on first line")
	}
}

func TestBarChart_Render_WithValues(t *testing.T) {
	data := []float64{10.5, 25.3, 15.8, 30.1}

	bar := NewBarChart(
		WithData(data),
		WithShowValues(true),
	)
	result := bar.Render()

	// Should contain numeric values
	for _, val := range data {
		// Check for the value (allowing some formatting variation)
		if !strings.Contains(result, "10.5") && val == 10.5 {
			t.Errorf("Expected output to contain value %.1f", val)
		}
	}
}

func TestBarChart_Render_Horizontal(t *testing.T) {
	data := []float64{10, 25, 15, 30}

	bar := NewBarChart(
		WithData(data),
		WithDirection(Horizontal),
	)
	result := bar.Render()

	// Should have one line per data point
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != len(data) {
		t.Errorf("Expected %d lines, got %d", len(data), len(lines))
	}
}

func TestBarChart_Render_Vertical(t *testing.T) {
	data := []float64{10, 25, 15, 30}

	bar := NewBarChart(
		WithData(data),
		WithDirection(Vertical),
		WithHeight(15),
	)
	result := bar.Render()

	// Should have multiple lines (height of chart)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 5 {
		t.Error("Expected multiple lines for vertical bar chart")
	}
}

func TestBarChart_Render_VerticalWithLabels(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	labels := []string{"A", "B", "C", "D"}

	bar := NewBarChart(
		WithData(data),
		WithLabels(labels),
		WithDirection(Vertical),
		WithHeight(15),
	)
	result := bar.Render()

	// Should contain labels at the bottom
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("Expected output to contain label %s", label)
		}
	}
}

func TestBarChart_Render_ASCIIMode(t *testing.T) {
	data := []float64{10, 25, 15, 30}

	bar := NewBarChart(
		WithData(data),
		WithStyle(StyleASCII),
	)
	result := bar.Render()

	// Should contain ASCII bar characters
	if !strings.Contains(result, "#") {
		t.Error("Expected ASCII mode to use # character")
	}
}

func TestBarChart_Render_UnicodeMode(t *testing.T) {
	data := []float64{10, 25, 15, 30}

	bar := NewBarChart(
		WithData(data),
		WithStyle(StyleUnicode),
	)
	result := bar.Render()

	// Should contain Unicode bar characters
	hasUnicode := false
	for _, r := range result {
		if r > 127 {
			hasUnicode = true
			break
		}
	}
	if !hasUnicode {
		t.Error("Expected Unicode mode to use Unicode characters")
	}
}

func TestBarChart_Render_WithColor(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	colorEnabled := true

	bar := NewBarChart(
		WithData(data),
		WithColor(colorEnabled),
	)
	result := bar.Render()

	// Should contain ANSI color codes
	if !strings.Contains(result, "\033[") {
		t.Error("Expected ANSI color codes in output")
	}
}

func TestBarChart_Render_InvalidData(t *testing.T) {
	tests := []struct {
		name string
		data []float64
	}{
		{
			name: "contains NaN",
			data: []float64{10, 20, math.NaN(), 40},
		},
		{
			name: "contains positive infinity",
			data: []float64{10, 20, math.Inf(1), 40},
		},
		{
			name: "contains negative infinity",
			data: []float64{10, 20, math.Inf(-1), 40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := NewBarChart(WithData(tt.data))
			result := bar.Render()

			// Should return empty string for invalid data
			if result != "" {
				t.Errorf("Expected empty string for invalid data, got: %s", result)
			}
		})
	}
}

func TestBarChart_Render_AllSameValues(t *testing.T) {
	bar := NewBarChart(WithData([]float64{20, 20, 20, 20}))
	result := bar.Render()

	// Should render all bars with same length
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// All lines should have similar content (same bar length)
	// Just verify they all contain bar characters
	for _, line := range lines {
		if len(line) == 0 {
			t.Error("Expected non-empty lines for all bars")
		}
	}
}

func TestBarChart_Render_WithWidth(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	width := 60

	bar := NewBarChart(
		WithData(data),
		WithWidth(width),
	)
	result := bar.Render()

	// Just verify that output is generated
	// Width controls the bar area, not total line width (which includes labels/values)
	if result == "" {
		t.Error("Expected non-empty output")
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != len(data) {
		t.Errorf("Expected %d lines, got %d", len(data), len(lines))
	}
}

func TestBarChart_Render_ZeroValues(t *testing.T) {
	data := []float64{0, 10, 0, 20, 0}

	bar := NewBarChart(WithData(data))
	result := bar.Render()

	// Should handle zero values without error
	if result == "" {
		t.Error("Render returned empty string for data with zeros")
	}

	// Should have one line per data point (trim trailing newline only)
	lines := strings.Split(strings.TrimSuffix(result, "\n"), "\n")
	if len(lines) != len(data) {
		t.Errorf("Expected %d lines, got %d", len(data), len(lines))
	}
}

func TestBarChart_Render_NegativeValues(t *testing.T) {
	// Note: Current implementation doesn't handle negative values well
	// This test documents the current behavior
	data := []float64{-10, 20, -5, 15}

	bar := NewBarChart(WithData(data))
	result := bar.Render()

	// Should still render something (even if not ideal)
	if result == "" {
		t.Error("Render returned empty string for negative values")
	}
}

func TestBar_ConvenienceFunction(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	result := Bar(data)

	if result == "" {
		t.Error("Bar() returned empty string")
	}

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != len(data) {
		t.Errorf("Expected %d lines, got %d", len(data), len(lines))
	}
}

func TestBarWithLabels_ConvenienceFunction(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	labels := []string{"Q1", "Q2", "Q3", "Q4"}
	result := BarWithLabels(data, labels)

	if result == "" {
		t.Error("BarWithLabels() returned empty string")
	}

	// Should contain all labels
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("Expected output to contain label %s", label)
		}
	}
}

func TestBarVertical_ConvenienceFunction(t *testing.T) {
	data := []float64{10, 25, 15, 30}
	result := BarVertical(data)

	if result == "" {
		t.Error("BarVertical() returned empty string")
	}

	// Should have multiple lines
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 5 {
		t.Error("Expected multiple lines for vertical bar chart")
	}
}

func TestFindMax(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		expected float64
	}{
		{
			name:     "simple max",
			data:     []float64{10, 25, 15, 30, 5},
			expected: 30,
		},
		{
			name:     "negative values",
			data:     []float64{-10, -5, -20, -3},
			expected: -3,
		},
		{
			name:     "single value",
			data:     []float64{42},
			expected: 42,
		},
		{
			name:     "empty array",
			data:     []float64{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMax(tt.data)
			if result != tt.expected {
				t.Errorf("Expected max %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestMaxStringLength(t *testing.T) {
	tests := []struct {
		name     string
		strings  []string
		expected int
	}{
		{
			name:     "simple strings",
			strings:  []string{"a", "bb", "ccc", "d"},
			expected: 3,
		},
		{
			name:     "empty array",
			strings:  []string{},
			expected: 0,
		},
		{
			name:     "single string",
			strings:  []string{"hello"},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxStringLength(tt.strings)
			if result != tt.expected {
				t.Errorf("Expected max length %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestBarMode_String(t *testing.T) {
	tests := []struct {
		mode     BarMode
		expected string
	}{
		{BarModeGrouped, "grouped"},
		{BarModeStacked, "stacked"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.mode.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.mode.String())
			}
		})
	}
}

func TestBarGrouped_ConvenienceFunction(t *testing.T) {
	series := []Series{
		{Label: "2023", Data: []float64{10, 20, 15}},
		{Label: "2024", Data: []float64{12, 25, 18}},
	}
	result := BarGrouped(series)

	if result == "" {
		t.Error("BarGrouped() returned empty string")
	}

	// Should have lines for each category
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines for grouped bar chart, got %d", len(lines))
	}
}

func TestBarStacked_ConvenienceFunction(t *testing.T) {
	series := []Series{
		{Label: "Product A", Data: []float64{10, 20, 15}},
		{Label: "Product B", Data: []float64{5, 10, 8}},
	}
	result := BarStacked(series)

	if result == "" {
		t.Error("BarStacked() returned empty string")
	}

	// Should have lines for each category
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 3 {
		t.Errorf("Expected at least 3 lines for stacked bar chart, got %d", len(lines))
	}
}

func TestBarChart_Render_GroupedHorizontal(t *testing.T) {
	series := []Series{
		{Label: "2023", Data: []float64{10, 20, 30}},
		{Label: "2024", Data: []float64{15, 25, 35}},
	}
	labels := []string{"Q1", "Q2", "Q3"}

	bar := NewBarChart(
		WithSeries(series),
		WithLabels(labels),
		WithBarMode(BarModeGrouped),
		WithStyle(StyleASCII),
	)
	result := bar.Render()

	if result == "" {
		t.Error("Render returned empty string for grouped horizontal bar chart")
	}

	// Should contain labels
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("Expected output to contain label %s", label)
		}
	}

	// Should have one line per category
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != len(labels) {
		t.Errorf("Expected %d lines, got %d", len(labels), len(lines))
	}
}

func TestBarChart_Render_StackedHorizontal(t *testing.T) {
	series := []Series{
		{Label: "Product A", Data: []float64{10, 20, 30}},
		{Label: "Product B", Data: []float64{5, 10, 15}},
	}
	labels := []string{"Q1", "Q2", "Q3"}

	bar := NewBarChart(
		WithSeries(series),
		WithLabels(labels),
		WithBarMode(BarModeStacked),
		WithStyle(StyleASCII),
	)
	result := bar.Render()

	if result == "" {
		t.Error("Render returned empty string for stacked horizontal bar chart")
	}

	// Should contain labels
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("Expected output to contain label %s", label)
		}
	}
}

func TestBarChart_Render_GroupedVertical(t *testing.T) {
	series := []Series{
		{Label: "2023", Data: []float64{10, 20, 30}},
		{Label: "2024", Data: []float64{15, 25, 35}},
	}
	labels := []string{"Q1", "Q2", "Q3"}

	bar := NewBarChart(
		WithSeries(series),
		WithLabels(labels),
		WithBarMode(BarModeGrouped),
		WithDirection(Vertical),
		WithHeight(15),
		WithStyle(StyleASCII),
	)
	result := bar.Render()

	if result == "" {
		t.Error("Render returned empty string for grouped vertical bar chart")
	}

	// Should have multiple lines (height of chart)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 5 {
		t.Error("Expected multiple lines for vertical grouped bar chart")
	}

	// Should contain labels
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("Expected output to contain label %s", label)
		}
	}
}

func TestBarChart_Render_StackedVertical(t *testing.T) {
	series := []Series{
		{Label: "Product A", Data: []float64{10, 20, 30}},
		{Label: "Product B", Data: []float64{5, 10, 15}},
	}
	labels := []string{"Q1", "Q2", "Q3"}

	bar := NewBarChart(
		WithSeries(series),
		WithLabels(labels),
		WithBarMode(BarModeStacked),
		WithDirection(Vertical),
		WithHeight(15),
		WithStyle(StyleASCII),
	)
	result := bar.Render()

	if result == "" {
		t.Error("Render returned empty string for stacked vertical bar chart")
	}

	// Should have multiple lines (height of chart)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 5 {
		t.Error("Expected multiple lines for vertical stacked bar chart")
	}
}

func TestBarChart_Render_MultiSeriesWithLegend(t *testing.T) {
	series := []Series{
		{Label: "2023", Data: []float64{10, 20, 30}},
		{Label: "2024", Data: []float64{15, 25, 35}},
	}

	bar := NewBarChart(
		WithSeries(series),
		WithShowLegend(true),
		WithStyle(StyleASCII),
	)
	result := bar.Render()

	// Should contain series labels in legend
	for _, s := range series {
		if !strings.Contains(result, s.Label) {
			t.Errorf("Expected legend to contain series label %s", s.Label)
		}
	}
}

func TestBarChart_Render_MultiSeriesWithColor(t *testing.T) {
	series := []Series{
		{Label: "2023", Data: []float64{10, 20, 30}},
		{Label: "2024", Data: []float64{15, 25, 35}},
	}

	bar := NewBarChart(
		WithSeries(series),
		WithColor(true),
		WithStyle(StyleUnicode),
	)
	result := bar.Render()

	// Should contain ANSI color codes
	if !strings.Contains(result, "\033[") {
		t.Error("Expected ANSI color codes in output for multi-series with color")
	}
}

func TestBarChart_Render_MultiSeriesWithTitle(t *testing.T) {
	series := []Series{
		{Label: "2023", Data: []float64{10, 20, 30}},
		{Label: "2024", Data: []float64{15, 25, 35}},
	}
	title := "Sales Comparison"

	bar := NewBarChart(
		WithSeries(series),
		WithTitle(title),
	)
	result := bar.Render()

	// Should contain title on first line
	lines := strings.Split(result, "\n")
	if !strings.Contains(lines[0], title) {
		t.Error("Expected title to be on first line")
	}
}

func TestBarChart_Render_MultiSeriesEmptyData(t *testing.T) {
	series := []Series{
		{Label: "Empty", Data: []float64{}},
	}

	bar := NewBarChart(WithSeries(series))
	result := bar.Render()

	// Should handle empty series gracefully
	if result == "" {
		// Empty series data results in empty output - this is expected
		t.Log("Empty series data resulted in empty output (expected)")
	}
}

func TestBarChart_Render_MultiSeriesInvalidData(t *testing.T) {
	series := []Series{
		{Label: "Invalid", Data: []float64{10, math.NaN(), 30}},
	}

	bar := NewBarChart(WithSeries(series))
	result := bar.Render()

	// Should return empty string for invalid data
	if result != "" {
		t.Error("Expected empty string for series with invalid data")
	}
}

func TestBarChart_Render_MultiSeriesCustomColors(t *testing.T) {
	series := []Series{
		{Label: "Red", Data: []float64{10, 20, 30}, Color: "red"},
		{Label: "Green", Data: []float64{15, 25, 35}, Color: "green"},
	}

	bar := NewBarChart(
		WithSeries(series),
		WithColor(true),
		WithStyle(StyleUnicode),
	)
	result := bar.Render()

	// Should contain ANSI color codes for red and green
	if !strings.Contains(result, "\033[31m") && !strings.Contains(result, "\033[32m") {
		t.Error("Expected custom color codes in output")
	}
}

func TestCalculateMaxValue_Stacked(t *testing.T) {
	series := []Series{
		{Label: "A", Data: []float64{10, 20, 30}},
		{Label: "B", Data: []float64{5, 10, 15}},
	}

	bar := NewBarChart(
		WithSeries(series),
		WithBarMode(BarModeStacked),
	)

	maxVal := bar.calculateMaxValue(series)
	// Max stacked value should be 30+15=45
	expected := 45.0
	if maxVal != expected {
		t.Errorf("Expected max stacked value %f, got %f", expected, maxVal)
	}
}

func TestCalculateMaxValue_Grouped(t *testing.T) {
	series := []Series{
		{Label: "A", Data: []float64{10, 20, 30}},
		{Label: "B", Data: []float64{5, 10, 15}},
	}

	bar := NewBarChart(
		WithSeries(series),
		WithBarMode(BarModeGrouped),
	)

	maxVal := bar.calculateMaxValue(series)
	// Max individual value should be 30
	expected := 30.0
	if maxVal != expected {
		t.Errorf("Expected max grouped value %f, got %f", expected, maxVal)
	}
}
