package termcharts

import (
	"math"
	"strings"
	"testing"
)

func TestNewPieChart(t *testing.T) {
	data := []float64{30, 25, 20, 15, 10}
	pie := NewPieChart(WithData(data))

	if pie == nil {
		t.Fatal("NewPieChart returned nil")
	}

	if pie.opts == nil {
		t.Fatal("PieChart has nil options")
	}

	if len(pie.opts.Data) != len(data) {
		t.Errorf("expected %d data points, got %d", len(data), len(pie.opts.Data))
	}
}

func TestPieChart_Render_BasicData(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		wantLen  bool
		contains []string
	}{
		{
			name:    "simple data",
			data:    []float64{50, 30, 20},
			wantLen: true,
		},
		{
			name:    "single value",
			data:    []float64{100},
			wantLen: true,
		},
		{
			name:    "equal values",
			data:    []float64{25, 25, 25, 25},
			wantLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pie := NewPieChart(
				WithData(tt.data),
				WithStyle(StyleUnicode),
				WithColor(false), // Explicit no-color to use distinct characters
			)
			result := pie.Render()

			if tt.wantLen && len(result) == 0 {
				t.Error("expected non-empty result")
			}

			// Should contain legend characters when colors disabled (●, ○, ◆, etc.)
			if !strings.Contains(result, "●") && !strings.Contains(result, "○") && !strings.Contains(result, "◆") {
				t.Error("expected Unicode slice characters in output")
			}
		})
	}
}

func TestPieChart_Render_EmptyData(t *testing.T) {
	pie := NewPieChart(WithData([]float64{}))
	result := pie.Render()

	if result != "" {
		t.Errorf("expected empty string for empty data, got: %q", result)
	}
}

func TestPieChart_Render_AllZeros(t *testing.T) {
	pie := NewPieChart(WithData([]float64{0, 0, 0}))
	result := pie.Render()

	if result != "" {
		t.Errorf("expected empty string for all zeros, got: %q", result)
	}
}

func TestPieChart_Render_WithLabels(t *testing.T) {
	data := []float64{30, 25, 20, 15, 10}
	labels := []string{"Chrome", "Firefox", "Safari", "Edge", "Other"}

	pie := NewPieChart(
		WithData(data),
		WithLabels(labels),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	// Should contain all labels
	for _, label := range labels {
		if !strings.Contains(result, label) {
			t.Errorf("expected output to contain label %q", label)
		}
	}

	// Should contain percentages
	if !strings.Contains(result, "%") {
		t.Error("expected output to contain percentages")
	}
}

func TestPieChart_Render_WithTitle(t *testing.T) {
	title := "Browser Market Share"
	pie := NewPieChart(
		WithData([]float64{30, 25, 20, 15, 10}),
		WithTitle(title),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	if !strings.Contains(result, title) {
		t.Errorf("expected output to contain title %q", title)
	}
}

func TestPieChart_Render_WithValues(t *testing.T) {
	data := []float64{30.5, 25.3, 20.2}
	pie := NewPieChart(
		WithData(data),
		WithShowValues(true),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	// Should contain value representations
	if !strings.Contains(result, "30.5") {
		t.Error("expected output to contain value 30.5")
	}
}

func TestPieChart_Render_ASCIIMode(t *testing.T) {
	pie := NewPieChart(
		WithData([]float64{50, 30, 20}),
		WithStyle(StyleASCII),
	)
	result := pie.Render()

	// Should contain ASCII slice characters (*, o, #, x, +, etc.)
	if !strings.Contains(result, "*") && !strings.Contains(result, "o") && !strings.Contains(result, "#") {
		t.Error("expected ASCII slice characters in output")
	}

	// Should NOT contain Unicode slice characters
	if strings.Contains(result, "●") || strings.Contains(result, "○") || strings.Contains(result, "◆") {
		t.Error("expected no Unicode characters in ASCII mode")
	}
}

func TestPieChart_Render_UnicodeMode(t *testing.T) {
	pie := NewPieChart(
		WithData([]float64{50, 30, 20}),
		WithStyle(StyleUnicode),
		WithColor(false), // Explicit no-color to use distinct characters
	)
	result := pie.Render()

	// Should contain Unicode legend characters when colors disabled (●, ○, ◆, etc.)
	if !strings.Contains(result, "●") && !strings.Contains(result, "○") && !strings.Contains(result, "◆") {
		t.Error("expected Unicode slice characters in output")
	}
}

func TestPieChart_Render_InvalidData(t *testing.T) {
	tests := []struct {
		name string
		data []float64
	}{
		{
			name: "contains NaN",
			data: []float64{10, math.NaN(), 20},
		},
		{
			name: "contains positive infinity",
			data: []float64{10, math.Inf(1), 20},
		},
		{
			name: "contains negative infinity",
			data: []float64{10, math.Inf(-1), 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pie := NewPieChart(WithData(tt.data))
			result := pie.Render()

			if result != "" {
				t.Errorf("expected empty string for invalid data, got: %q", result)
			}
		})
	}
}

func TestPieChart_Render_NegativeValues(t *testing.T) {
	// Negative values should be treated as 0
	pie := NewPieChart(
		WithData([]float64{30, -10, 20}),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	// Should still render (negative treated as 0)
	if len(result) == 0 {
		t.Error("expected non-empty result even with negative values")
	}
}

func TestPieChart_Render_WithColor(t *testing.T) {
	pie := NewPieChart(
		WithData([]float64{50, 30, 20}),
		WithColor(true),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	// Should contain ANSI escape codes
	if !strings.Contains(result, "\033[") {
		t.Error("expected ANSI color codes in output")
	}
}

func TestPieChart_Render_WithTheme(t *testing.T) {
	pie := NewPieChart(
		WithData([]float64{50, 30, 20}),
		WithTheme(DarkTheme),
		WithColor(true),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	// Should render without error
	if len(result) == 0 {
		t.Error("expected non-empty result with theme")
	}
}

func TestPieChart_Render_WithWidth(t *testing.T) {
	pie := NewPieChart(
		WithData([]float64{50, 30, 20}),
		WithWidth(40),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	// Should render without error
	if len(result) == 0 {
		t.Error("expected non-empty result with custom width")
	}
}

func TestPie_ConvenienceFunction(t *testing.T) {
	result := Pie([]float64{50, 30, 20})

	if len(result) == 0 {
		t.Error("Pie convenience function returned empty string")
	}
}

func TestPieWithLabels_ConvenienceFunction(t *testing.T) {
	result := PieWithLabels(
		[]float64{50, 30, 20},
		[]string{"A", "B", "C"},
	)

	if len(result) == 0 {
		t.Error("PieWithLabels convenience function returned empty string")
	}

	if !strings.Contains(result, "A") {
		t.Error("expected output to contain label 'A'")
	}
}

func TestPieWithValues_ConvenienceFunction(t *testing.T) {
	result := PieWithValues(
		[]float64{50.5, 30.3, 20.2},
		[]string{"A", "B", "C"},
	)

	if len(result) == 0 {
		t.Error("PieWithValues convenience function returned empty string")
	}

	// Should contain value
	if !strings.Contains(result, "50.5") {
		t.Error("expected output to contain value 50.5")
	}
}

func TestPieChart_calculateSlices(t *testing.T) {
	pie := NewPieChart(
		WithData([]float64{50, 25, 25}),
		WithLabels([]string{"Half", "Quarter1", "Quarter2"}),
	)

	slices := pie.calculateSlices()

	if len(slices) != 3 {
		t.Fatalf("expected 3 slices, got %d", len(slices))
	}

	// Check percentages
	if math.Abs(slices[0].Percentage-50.0) > 0.01 {
		t.Errorf("expected first slice to be 50%%, got %.2f%%", slices[0].Percentage)
	}

	if math.Abs(slices[1].Percentage-25.0) > 0.01 {
		t.Errorf("expected second slice to be 25%%, got %.2f%%", slices[1].Percentage)
	}

	// Check labels
	if slices[0].Label != "Half" {
		t.Errorf("expected first slice label 'Half', got %q", slices[0].Label)
	}
}

func TestPieChart_calculateSlices_DefaultLabels(t *testing.T) {
	pie := NewPieChart(WithData([]float64{50, 30, 20}))

	slices := pie.calculateSlices()

	// Should have default labels
	if slices[0].Label != "Item 1" {
		t.Errorf("expected default label 'Item 1', got %q", slices[0].Label)
	}
}

func TestPieChart_Percentages(t *testing.T) {
	// Test that percentages add up to approximately 100%
	pie := NewPieChart(
		WithData([]float64{33.33, 33.33, 33.34}),
		WithLabels([]string{"A", "B", "C"}),
	)

	slices := pie.calculateSlices()

	totalPercent := 0.0
	for _, s := range slices {
		totalPercent += s.Percentage
	}

	if math.Abs(totalPercent-100.0) > 0.1 {
		t.Errorf("expected total percentage ~100%%, got %.2f%%", totalPercent)
	}
}

func TestPieChart_ManySlices(t *testing.T) {
	// Test with many slices
	data := make([]float64, 10)
	labels := make([]string, 10)
	for i := range data {
		data[i] = float64(10)
		labels[i] = string(rune('A' + i))
	}

	pie := NewPieChart(
		WithData(data),
		WithLabels(labels),
		WithStyle(StyleUnicode),
	)
	result := pie.Render()

	if len(result) == 0 {
		t.Error("expected non-empty result with many slices")
	}

	// Each slice should be 10%
	if !strings.Contains(result, "10.0%") {
		t.Error("expected slices to show 10.0%")
	}
}
