// Last modified: 2026-01-02

package internal

import "math"

// Normalize scales data to the range [0, 1].
// Returns normalized data and the original min/max values.
func Normalize(data []float64) ([]float64, float64, float64) {
	if len(data) == 0 {
		return data, 0, 0
	}

	min, max := MinMax(data)
	if min == max {
		// All values are the same
		normalized := make([]float64, len(data))
		for i := range normalized {
			normalized[i] = 0.5
		}
		return normalized, min, max
	}

	normalized := make([]float64, len(data))
	scale := max - min
	for i, v := range data {
		normalized[i] = (v - min) / scale
	}

	return normalized, min, max
}

// Scale maps data from one range to another.
// Maps [dataMin, dataMax] to [targetMin, targetMax].
func Scale(value, dataMin, dataMax, targetMin, targetMax float64) float64 {
	if dataMax == dataMin {
		return targetMin
	}
	normalized := (value - dataMin) / (dataMax - dataMin)
	return targetMin + normalized*(targetMax-targetMin)
}

// MinMax returns the minimum and maximum values in the data.
// Returns (0, 0) for empty data.
func MinMax(data []float64) (min, max float64) {
	if len(data) == 0 {
		return 0, 0
	}

	min = data[0]
	max = data[0]

	for _, v := range data[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	return min, max
}

// Clamp constrains a value to the range [min, max].
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt constrains an integer value to the range [min, max].
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// IsValid returns true if the value is a valid number (not NaN or Inf).
func IsValid(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

// AllValid returns true if all values in the data are valid (not NaN or Inf).
func AllValid(data []float64) bool {
	for _, v := range data {
		if !IsValid(v) {
			return false
		}
	}
	return true
}

// Round rounds a float64 to the nearest integer.
func Round(value float64) int {
	return int(math.Round(value))
}

// Abs returns the absolute value of an integer.
func Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// Max returns the larger of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Min returns the smaller of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
