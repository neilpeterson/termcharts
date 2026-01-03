package termcharts

import (
	"strings"
	"testing"
)

func TestRenderStyle_String(t *testing.T) {
	tests := []struct {
		name     string
		style    RenderStyle
		expected string
	}{
		{
			name:     "StyleAuto",
			style:    StyleAuto,
			expected: "auto",
		},
		{
			name:     "StyleASCII",
			style:    StyleASCII,
			expected: "ascii",
		},
		{
			name:     "StyleUnicode",
			style:    StyleUnicode,
			expected: "unicode",
		},
		{
			name:     "StyleBraille",
			style:    StyleBraille,
			expected: "braille",
		},
		{
			name:     "unknown style",
			style:    RenderStyle(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.style.String()
			if result != tt.expected {
				t.Errorf("RenderStyle.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		color        string
		colorEnabled bool
		expectColor  bool
	}{
		{
			name:         "color enabled with blue",
			text:         "test",
			color:        "blue",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "color disabled",
			text:         "test",
			color:        "blue",
			colorEnabled: false,
			expectColor:  false,
		},
		{
			name:         "empty color",
			text:         "test",
			color:        "",
			colorEnabled: true,
			expectColor:  false,
		},
		{
			name:         "invalid color",
			text:         "test",
			color:        "invalidcolor",
			colorEnabled: true,
			expectColor:  false,
		},
		{
			name:         "red color",
			text:         "test",
			color:        "red",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "green color",
			text:         "test",
			color:        "green",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "yellow color",
			text:         "test",
			color:        "yellow",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "cyan color",
			text:         "test",
			color:        "cyan",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "magenta color",
			text:         "test",
			color:        "magenta",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "white color",
			text:         "test",
			color:        "white",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "gray color",
			text:         "test",
			color:        "gray",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "grey color (alternative spelling)",
			text:         "test",
			color:        "grey",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "orange color alias",
			text:         "test",
			color:        "orange",
			colorEnabled: true,
			expectColor:  true,
		},
		{
			name:         "purple color alias",
			text:         "test",
			color:        "purple",
			colorEnabled: true,
			expectColor:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.text, tt.color, tt.colorEnabled)

			if tt.expectColor {
				// Should contain ANSI escape codes
				if !strings.Contains(result, "\033[") {
					t.Errorf("Colorize() expected to contain ANSI codes but got: %v", result)
				}
				// Should contain the original text
				if !strings.Contains(result, tt.text) {
					t.Errorf("Colorize() should contain original text %v", tt.text)
				}
				// Should contain reset code
				if !strings.Contains(result, colorReset) {
					t.Errorf("Colorize() should contain reset code")
				}
			} else {
				// Should return text unchanged
				if result != tt.text {
					t.Errorf("Colorize() = %v, want %v", result, tt.text)
				}
			}
		})
	}
}

func TestTheme_GetSeriesColor(t *testing.T) {
	tests := []struct {
		name     string
		theme    *Theme
		index    int
		expected string
	}{
		{
			name:     "first color",
			theme:    DefaultTheme,
			index:    0,
			expected: "blue",
		},
		{
			name:     "second color",
			theme:    DefaultTheme,
			index:    1,
			expected: "green",
		},
		{
			name:     "cycle to first color",
			theme:    DefaultTheme,
			index:    6,
			expected: "blue",
		},
		{
			name:     "cycle multiple times",
			theme:    DefaultTheme,
			index:    12,
			expected: "blue",
		},
		{
			name: "empty series colors",
			theme: &Theme{
				Primary: "red",
				Series:  []string{},
			},
			index:    0,
			expected: "red",
		},
		{
			name: "single color series",
			theme: &Theme{
				Primary: "red",
				Series:  []string{"cyan"},
			},
			index:    5,
			expected: "cyan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.theme.GetSeriesColor(tt.index)
			if result != tt.expected {
				t.Errorf("GetSeriesColor(%d) = %v, want %v", tt.index, result, tt.expected)
			}
		})
	}
}

func TestPredefinedThemes(t *testing.T) {
	tests := []struct {
		name  string
		theme *Theme
	}{
		{
			name:  "DefaultTheme",
			theme: DefaultTheme,
		},
		{
			name:  "DarkTheme",
			theme: DarkTheme,
		},
		{
			name:  "LightTheme",
			theme: LightTheme,
		},
		{
			name:  "MonochromeTheme",
			theme: MonochromeTheme,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.theme == nil {
				t.Errorf("%s is nil", tt.name)
			}

			if tt.theme.Primary == "" {
				t.Errorf("%s.Primary is empty", tt.name)
			}

			if len(tt.theme.Series) == 0 {
				t.Errorf("%s.Series is empty", tt.name)
			}
		})
	}
}

func TestColorMap(t *testing.T) {
	// Test that all expected colors are in the map
	expectedColors := []string{
		"black", "red", "green", "yellow", "orange",
		"blue", "magenta", "purple", "cyan", "white",
		"gray", "grey", "brown",
	}

	for _, color := range expectedColors {
		if _, ok := colorMap[color]; !ok {
			t.Errorf("colorMap missing expected color: %s", color)
		}
	}
}

func TestColorConstants(t *testing.T) {
	// Test that color constants are properly defined
	if colorReset != "\033[0m" {
		t.Errorf("colorReset = %q, want %q", colorReset, "\033[0m")
	}

	if colorBlue != "\033[34m" {
		t.Errorf("colorBlue = %q, want %q", colorBlue, "\033[34m")
	}

	if colorRed != "\033[31m" {
		t.Errorf("colorRed = %q, want %q", colorRed, "\033[31m")
	}
}
