package termcharts

import (
	"errors"
	"testing"
)

func TestDirection_String(t *testing.T) {
	tests := []struct {
		name     string
		dir      Direction
		expected string
	}{
		{
			name:     "horizontal direction",
			dir:      Horizontal,
			expected: "horizontal",
		},
		{
			name:     "vertical direction",
			dir:      Vertical,
			expected: "vertical",
		},
		{
			name:     "unknown direction",
			dir:      Direction(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dir.String()
			if result != tt.expected {
				t.Errorf("Direction.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrEmptyData",
			err:      ErrEmptyData,
			expected: "data cannot be empty",
		},
		{
			name:     "ErrInvalidData",
			err:      ErrInvalidData,
			expected: "data contains invalid values",
		},
		{
			name:     "ErrInvalidDimensions",
			err:      ErrInvalidDimensions,
			expected: "chart dimensions too small",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("error message = %v, want %v", tt.err.Error(), tt.expected)
			}
		})
	}
}

func TestErrors_Is(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "ErrEmptyData is ErrEmptyData",
			err:      ErrEmptyData,
			target:   ErrEmptyData,
			expected: true,
		},
		{
			name:     "ErrInvalidData is ErrInvalidData",
			err:      ErrInvalidData,
			target:   ErrInvalidData,
			expected: true,
		},
		{
			name:     "ErrEmptyData is not ErrInvalidData",
			err:      ErrEmptyData,
			target:   ErrInvalidData,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("errors.Is() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSeries(t *testing.T) {
	s := Series{
		Label: "Test Series",
		Data:  []float64{1, 2, 3, 4, 5},
		Color: "blue",
	}

	if s.Label != "Test Series" {
		t.Errorf("Series.Label = %v, want %v", s.Label, "Test Series")
	}

	if len(s.Data) != 5 {
		t.Errorf("len(Series.Data) = %v, want %v", len(s.Data), 5)
	}

	if s.Color != "blue" {
		t.Errorf("Series.Color = %v, want %v", s.Color, "blue")
	}
}
