package util

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name         string
		data         []float64
		expectedNorm []float64
		expectedMin  float64
		expectedMax  float64
	}{
		{
			name:         "simple range",
			data:         []float64{0, 5, 10},
			expectedNorm: []float64{0, 0.5, 1},
			expectedMin:  0,
			expectedMax:  10,
		},
		{
			name:         "negative values",
			data:         []float64{-10, 0, 10},
			expectedNorm: []float64{0, 0.5, 1},
			expectedMin:  -10,
			expectedMax:  10,
		},
		{
			name:         "all same values",
			data:         []float64{5, 5, 5, 5},
			expectedNorm: []float64{0.5, 0.5, 0.5, 0.5},
			expectedMin:  5,
			expectedMax:  5,
		},
		{
			name:         "empty data",
			data:         []float64{},
			expectedNorm: []float64{},
			expectedMin:  0,
			expectedMax:  0,
		},
		{
			name:         "single value",
			data:         []float64{42},
			expectedNorm: []float64{0.5},
			expectedMin:  42,
			expectedMax:  42,
		},
		{
			name:         "fractional values",
			data:         []float64{1.5, 2.5, 3.5},
			expectedNorm: []float64{0, 0.5, 1},
			expectedMin:  1.5,
			expectedMax:  3.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized, min, max := Normalize(tt.data)

			if min != tt.expectedMin {
				t.Errorf("Normalize() min = %v, want %v", min, tt.expectedMin)
			}

			if max != tt.expectedMax {
				t.Errorf("Normalize() max = %v, want %v", max, tt.expectedMax)
			}

			if len(normalized) != len(tt.expectedNorm) {
				t.Fatalf("Normalize() length = %v, want %v", len(normalized), len(tt.expectedNorm))
			}

			for i, v := range normalized {
				if math.Abs(v-tt.expectedNorm[i]) > 1e-10 {
					t.Errorf("Normalize()[%d] = %v, want %v", i, v, tt.expectedNorm[i])
				}
			}
		})
	}
}

func TestScale(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		dataMin   float64
		dataMax   float64
		targetMin float64
		targetMax float64
		expected  float64
	}{
		{
			name:      "scale 0-10 to 0-100",
			value:     5,
			dataMin:   0,
			dataMax:   10,
			targetMin: 0,
			targetMax: 100,
			expected:  50,
		},
		{
			name:      "scale negative range",
			value:     0,
			dataMin:   -10,
			dataMax:   10,
			targetMin: 0,
			targetMax: 100,
			expected:  50,
		},
		{
			name:      "scale to smaller range",
			value:     50,
			dataMin:   0,
			dataMax:   100,
			targetMin: 0,
			targetMax: 10,
			expected:  5,
		},
		{
			name:      "same min and max",
			value:     5,
			dataMin:   5,
			dataMax:   5,
			targetMin: 0,
			targetMax: 100,
			expected:  0,
		},
		{
			name:      "value at minimum",
			value:     0,
			dataMin:   0,
			dataMax:   10,
			targetMin: 0,
			targetMax: 100,
			expected:  0,
		},
		{
			name:      "value at maximum",
			value:     10,
			dataMin:   0,
			dataMax:   10,
			targetMin: 0,
			targetMax: 100,
			expected:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Scale(tt.value, tt.dataMin, tt.dataMax, tt.targetMin, tt.targetMax)

			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("Scale() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		name        string
		data        []float64
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "positive values",
			data:        []float64{1, 5, 3, 9, 2},
			expectedMin: 1,
			expectedMax: 9,
		},
		{
			name:        "negative values",
			data:        []float64{-5, -2, -10, -1},
			expectedMin: -10,
			expectedMax: -1,
		},
		{
			name:        "mixed values",
			data:        []float64{-5, 0, 5, 10},
			expectedMin: -5,
			expectedMax: 10,
		},
		{
			name:        "single value",
			data:        []float64{42},
			expectedMin: 42,
			expectedMax: 42,
		},
		{
			name:        "empty data",
			data:        []float64{},
			expectedMin: 0,
			expectedMax: 0,
		},
		{
			name:        "all same values",
			data:        []float64{7, 7, 7},
			expectedMin: 7,
			expectedMax: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := MinMax(tt.data)

			if min != tt.expectedMin {
				t.Errorf("MinMax() min = %v, want %v", min, tt.expectedMin)
			}

			if max != tt.expectedMax {
				t.Errorf("MinMax() max = %v, want %v", max, tt.expectedMax)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{
			name:     "value within range",
			value:    5,
			min:      0,
			max:      10,
			expected: 5,
		},
		{
			name:     "value below minimum",
			value:    -5,
			min:      0,
			max:      10,
			expected: 0,
		},
		{
			name:     "value above maximum",
			value:    15,
			min:      0,
			max:      10,
			expected: 10,
		},
		{
			name:     "value at minimum",
			value:    0,
			min:      0,
			max:      10,
			expected: 0,
		},
		{
			name:     "value at maximum",
			value:    10,
			min:      0,
			max:      10,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamp(tt.value, tt.min, tt.max)

			if result != tt.expected {
				t.Errorf("Clamp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClampInt(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		min      int
		max      int
		expected int
	}{
		{
			name:     "value within range",
			value:    5,
			min:      0,
			max:      10,
			expected: 5,
		},
		{
			name:     "value below minimum",
			value:    -5,
			min:      0,
			max:      10,
			expected: 0,
		},
		{
			name:     "value above maximum",
			value:    15,
			min:      0,
			max:      10,
			expected: 10,
		},
		{
			name:     "value at minimum",
			value:    0,
			min:      0,
			max:      10,
			expected: 0,
		},
		{
			name:     "value at maximum",
			value:    10,
			min:      0,
			max:      10,
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClampInt(tt.value, tt.min, tt.max)

			if result != tt.expected {
				t.Errorf("ClampInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected bool
	}{
		{
			name:     "valid positive number",
			value:    42.5,
			expected: true,
		},
		{
			name:     "valid negative number",
			value:    -42.5,
			expected: true,
		},
		{
			name:     "zero",
			value:    0,
			expected: true,
		},
		{
			name:     "NaN",
			value:    math.NaN(),
			expected: false,
		},
		{
			name:     "positive infinity",
			value:    math.Inf(1),
			expected: false,
		},
		{
			name:     "negative infinity",
			value:    math.Inf(-1),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValid(tt.value)

			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAllValid(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		expected bool
	}{
		{
			name:     "all valid",
			data:     []float64{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "contains NaN",
			data:     []float64{1, 2, math.NaN(), 4, 5},
			expected: false,
		},
		{
			name:     "contains infinity",
			data:     []float64{1, 2, math.Inf(1), 4, 5},
			expected: false,
		},
		{
			name:     "empty data",
			data:     []float64{},
			expected: true,
		},
		{
			name:     "all NaN",
			data:     []float64{math.NaN(), math.NaN()},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AllValid(tt.data)

			if result != tt.expected {
				t.Errorf("AllValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected int
	}{
		{
			name:     "round up",
			value:    5.6,
			expected: 6,
		},
		{
			name:     "round down",
			value:    5.4,
			expected: 5,
		},
		{
			name:     "round half",
			value:    5.5,
			expected: 6,
		},
		{
			name:     "negative round up",
			value:    -5.4,
			expected: -5,
		},
		{
			name:     "negative round down",
			value:    -5.6,
			expected: -6,
		},
		{
			name:     "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.value)

			if result != tt.expected {
				t.Errorf("Round() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{
			name:     "positive number",
			value:    5,
			expected: 5,
		},
		{
			name:     "negative number",
			value:    -5,
			expected: 5,
		},
		{
			name:     "zero",
			value:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.value)

			if result != tt.expected {
				t.Errorf("Abs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "a greater",
			a:        10,
			b:        5,
			expected: 10,
		},
		{
			name:     "b greater",
			a:        5,
			b:        10,
			expected: 10,
		},
		{
			name:     "equal",
			a:        5,
			b:        5,
			expected: 5,
		},
		{
			name:     "negative values",
			a:        -5,
			b:        -10,
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)

			if result != tt.expected {
				t.Errorf("Max() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "a smaller",
			a:        5,
			b:        10,
			expected: 5,
		},
		{
			name:     "b smaller",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "equal",
			a:        5,
			b:        5,
			expected: 5,
		},
		{
			name:     "negative values",
			a:        -5,
			b:        -10,
			expected: -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)

			if result != tt.expected {
				t.Errorf("Min() = %v, want %v", result, tt.expected)
			}
		})
	}
}
