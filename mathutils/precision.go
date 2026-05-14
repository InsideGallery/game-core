package mathutils

import "math"

const defaultPrecision = 0.0001

// RoundWithPrecision rounds value to the given precision.
func RoundWithPrecision(value float64, precision float64) float64 {
	if ApproximatelyEqual(precision, 0) {
		precision = defaultPrecision
	}

	precision = 1 / precision

	return math.Round(value*precision) / precision
}

// ApproximatelyEqual reports whether two floating point values are almost equal.
func ApproximatelyEqual(a, b float64) bool {
	difference := a - b

	return difference < math.SmallestNonzeroFloat64 && difference > -math.SmallestNonzeroFloat64
}
