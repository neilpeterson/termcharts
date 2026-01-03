package internal

import (
	"os"
	"testing"
)

func TestGetTerminalSize(t *testing.T) {
	// Basic test - should not panic and return valid dimensions
	size := GetTerminalSize()

	if size.Width <= 0 {
		t.Errorf("GetTerminalSize().Width = %d, want > 0", size.Width)
	}

	if size.Height <= 0 {
		t.Errorf("GetTerminalSize().Height = %d, want > 0", size.Height)
	}
}

func TestGetTerminalSize_DefaultSize(t *testing.T) {
	// Test default size values
	if DefaultSize.Width != 80 {
		t.Errorf("DefaultSize.Width = %d, want 80", DefaultSize.Width)
	}

	if DefaultSize.Height != 24 {
		t.Errorf("DefaultSize.Height = %d, want 24", DefaultSize.Height)
	}
}

func TestGetSizeFromEnv(t *testing.T) {
	tests := []struct {
		name           string
		columns        string
		lines          string
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "valid environment variables",
			columns:        "100",
			lines:          "30",
			expectedWidth:  100,
			expectedHeight: 30,
		},
		{
			name:           "only columns set",
			columns:        "120",
			lines:          "",
			expectedWidth:  120,
			expectedHeight: 0,
		},
		{
			name:           "only lines set",
			columns:        "",
			lines:          "40",
			expectedWidth:  0,
			expectedHeight: 40,
		},
		{
			name:           "invalid values",
			columns:        "invalid",
			lines:          "invalid",
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name:           "empty values",
			columns:        "",
			lines:          "",
			expectedWidth:  0,
			expectedHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origColumns := os.Getenv("COLUMNS")
			origLines := os.Getenv("LINES")
			defer func() {
				os.Setenv("COLUMNS", origColumns)
				os.Setenv("LINES", origLines)
			}()

			// Set test env vars
			if tt.columns != "" {
				os.Setenv("COLUMNS", tt.columns)
			} else {
				os.Unsetenv("COLUMNS")
			}

			if tt.lines != "" {
				os.Setenv("LINES", tt.lines)
			} else {
				os.Unsetenv("LINES")
			}

			width, height := getSizeFromEnv()

			if width != tt.expectedWidth {
				t.Errorf("getSizeFromEnv() width = %d, want %d", width, tt.expectedWidth)
			}

			if height != tt.expectedHeight {
				t.Errorf("getSizeFromEnv() height = %d, want %d", height, tt.expectedHeight)
			}
		})
	}
}

func TestSupportsColor(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name: "NO_COLOR set",
			envVars: map[string]string{
				"NO_COLOR": "1",
				"TERM":     "xterm-256color",
			},
			expected: false,
		},
		{
			name: "FORCE_COLOR set",
			envVars: map[string]string{
				"FORCE_COLOR": "1",
			},
			expected: true,
		},
		{
			name: "TERM is dumb",
			envVars: map[string]string{
				"TERM": "dumb",
			},
			expected: false,
		},
		{
			name: "TERM is empty",
			envVars: map[string]string{
				"TERM": "",
			},
			expected: false,
		},
		{
			name: "TERM contains xterm",
			envVars: map[string]string{
				"TERM": "xterm-256color",
			},
			expected: true,
		},
		{
			name: "TERM contains screen",
			envVars: map[string]string{
				"TERM": "screen-256color",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origVars := make(map[string]string)
			envKeys := []string{"NO_COLOR", "FORCE_COLOR", "TERM", "WT_SESSION", "ConEmuANSI"}
			for _, key := range envKeys {
				origVars[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				for key, val := range origVars {
					if val != "" {
						os.Setenv(key, val)
					} else {
						os.Unsetenv(key)
					}
				}
			}()

			// Set test env vars
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			result := SupportsColor()
			if result != tt.expected {
				t.Errorf("SupportsColor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSupportsUnicode(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name: "LANG is C",
			envVars: map[string]string{
				"LANG": "C",
			},
			expected: false,
		},
		{
			name: "LC_ALL is C",
			envVars: map[string]string{
				"LC_ALL": "C",
			},
			expected: false,
		},
		{
			name: "UTF-8 locale",
			envVars: map[string]string{
				"LANG": "en_US.UTF-8",
			},
			expected: true,
		},
		{
			name: "UTF8 locale (no dash)",
			envVars: map[string]string{
				"LANG": "en_US.UTF8",
			},
			expected: true,
		},
		{
			name: "LC_ALL UTF-8",
			envVars: map[string]string{
				"LC_ALL": "en_US.UTF-8",
			},
			expected: true,
		},
		{
			name: "LC_CTYPE UTF-8",
			envVars: map[string]string{
				"LC_CTYPE": "en_US.UTF-8",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env vars
			origVars := make(map[string]string)
			envKeys := []string{"LANG", "LC_ALL", "LC_CTYPE"}
			for _, key := range envKeys {
				origVars[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				for key, val := range origVars {
					if val != "" {
						os.Setenv(key, val)
					} else {
						os.Unsetenv(key)
					}
				}
			}()

			// Set test env vars
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			result := SupportsUnicode()
			if result != tt.expected {
				t.Errorf("SupportsUnicode() = %v, want %v (env: %v)", result, tt.expected, tt.envVars)
			}
		})
	}
}

func TestIsTTY(t *testing.T) {
	// Basic test - should not panic
	result := IsTTY()

	// Result depends on test environment, just ensure it returns a bool
	_ = result
}

func TestTerminalSize_Struct(t *testing.T) {
	size := TerminalSize{
		Width:  100,
		Height: 50,
	}

	if size.Width != 100 {
		t.Errorf("TerminalSize.Width = %d, want 100", size.Width)
	}

	if size.Height != 50 {
		t.Errorf("TerminalSize.Height = %d, want 50", size.Height)
	}
}
