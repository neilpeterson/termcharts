// Last modified: 2026-01-02

package termcharts

import (
	"math"
	"strings"
	"testing"
)

func TestNewSparkline(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	spark := NewSparkline(WithData(data))

	if spark == nil {
		t.Fatal("NewSparkline returned nil")
	}
	if spark.opts == nil {
		t.Fatal("Options not initialized")
	}
	if len(spark.opts.Data) != len(data) {
		t.Errorf("Expected data length %d, got %d", len(data), len(spark.opts.Data))
	}
}

func TestSparkline_Render_BasicData(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		expected int // Expected length of output
	}{
		{
			name:     "simple ascending",
			data:     []float64{1, 2, 3, 4, 5},
			expected: 5,
		},
		{
			name:     "simple descending",
			data:     []float64{5, 4, 3, 2, 1},
			expected: 5,
		},
		{
			name:     "mixed values",
			data:     []float64{1, 5, 2, 8, 3, 7},
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spark := NewSparkline(WithData(tt.data))
			result := spark.Render()

			// Count runes (not bytes) for Unicode characters
			if len([]rune(result)) != tt.expected {
				t.Errorf("Expected %d characters, got %d: %s", tt.expected, len([]rune(result)), result)
			}

			// Should not be empty
			if result == "" {
				t.Error("Render returned empty string")
			}
		})
	}
}

func TestSparkline_Render_EmptyData(t *testing.T) {
	spark := NewSparkline(WithData([]float64{}))
	result := spark.Render()

	if result != "" {
		t.Errorf("Expected empty string for empty data, got: %s", result)
	}
}

func TestSparkline_Render_SingleValue(t *testing.T) {
	spark := NewSparkline(WithData([]float64{42}))
	result := spark.Render()

	// Should render a single character
	if len([]rune(result)) != 1 {
		t.Errorf("Expected 1 character, got %d", len([]rune(result)))
	}
}

func TestSparkline_Render_AllSameValues(t *testing.T) {
	spark := NewSparkline(WithData([]float64{5, 5, 5, 5, 5}))
	result := spark.Render()

	// Should render all middle characters (since all normalized to 0.5)
	if len([]rune(result)) != 5 {
		t.Errorf("Expected 5 characters, got %d", len([]rune(result)))
	}

	// All characters should be the same
	runes := []rune(result)
	for i := 1; i < len(runes); i++ {
		if runes[i] != runes[0] {
			t.Error("Expected all characters to be the same for identical values")
			break
		}
	}
}

func TestSparkline_Render_ASCIIMode(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	spark := NewSparkline(
		WithData(data),
		WithStyle(StyleASCII),
	)
	result := spark.Render()

	// Check that only ASCII characters are used
	for _, r := range result {
		if r > 127 {
			t.Errorf("Expected ASCII only, got Unicode character: %c", r)
		}
	}

	// Verify it uses the ASCII character set
	for _, r := range result {
		found := false
		for _, asciiChar := range sparkCharsASCII {
			if r == asciiChar {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Character %c not in ASCII sparkline character set", r)
		}
	}
}

func TestSparkline_Render_UnicodeMode(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	spark := NewSparkline(
		WithData(data),
		WithStyle(StyleUnicode),
	)
	result := spark.Render()

	// Verify it uses the Unicode character set
	runes := []rune(result)
	for _, r := range runes {
		found := false
		for _, unicodeChar := range sparkChars {
			if r == unicodeChar {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Character %c not in Unicode sparkline character set", r)
		}
	}
}

func TestSparkline_Render_InvalidData(t *testing.T) {
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
			spark := NewSparkline(WithData(tt.data))
			result := spark.Render()

			// Should return empty string for invalid data
			if result != "" {
				t.Errorf("Expected empty string for invalid data, got: %s", result)
			}
		})
	}
}

func TestSparkline_Render_WithWidth(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	width := 5

	spark := NewSparkline(
		WithData(data),
		WithWidth(width),
	)
	result := spark.Render()

	// Should be limited to specified width
	if len([]rune(result)) != width {
		t.Errorf("Expected width %d, got %d", width, len([]rune(result)))
	}
}

func TestSparkline_Render_WithColor(t *testing.T) {
	data := []float64{1, 5, 2, 8, 3}
	colorEnabled := true

	spark := NewSparkline(
		WithData(data),
		WithColor(colorEnabled),
	)
	result := spark.Render()

	// Should contain ANSI color codes
	if !strings.Contains(result, "\033[") {
		t.Error("Expected ANSI color codes in output")
	}
}

func TestSparkline_ConvenienceFunction(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	result := Spark(data)

	if result == "" {
		t.Error("Spark() returned empty string")
	}

	if len([]rune(result)) != len(data) {
		t.Errorf("Expected %d characters, got %d", len(data), len([]rune(result)))
	}
}

func TestSparkASCII_ConvenienceFunction(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	result := SparkASCII(data)

	if result == "" {
		t.Error("SparkASCII() returned empty string")
	}

	// Should only contain ASCII characters
	for _, r := range result {
		if r > 127 {
			t.Errorf("Expected ASCII only, got Unicode character: %c", r)
		}
	}
}

func TestSparkColor_ConvenienceFunction(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	result := SparkColor(data)

	if result == "" {
		t.Error("SparkColor() returned empty string")
	}

	// Result may or may not contain colors depending on terminal detection
	// Just verify it returns something
}

func TestSparkline_Render_CharacterMapping(t *testing.T) {
	// Test that min value maps to lowest character and max to highest
	data := []float64{0, 100}

	spark := NewSparkline(
		WithData(data),
		WithStyle(StyleUnicode),
	)
	result := spark.Render()

	runes := []rune(result)
	if len(runes) != 2 {
		t.Fatalf("Expected 2 characters, got %d", len(runes))
	}

	// First should be minimum character
	if runes[0] != sparkChars[0] {
		t.Errorf("Expected min character %c, got %c", sparkChars[0], runes[0])
	}

	// Last should be maximum character
	if runes[1] != sparkChars[len(sparkChars)-1] {
		t.Errorf("Expected max character %c, got %c", sparkChars[len(sparkChars)-1], runes[1])
	}
}
