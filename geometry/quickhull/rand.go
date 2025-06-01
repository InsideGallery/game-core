package quickhull

import "math/rand/v2"

func FastRandFloat64(precision, min, max float64) float64 {
	if precision <= 0 {
		return 0
	}

	return float64(FastRand(int(min*precision), int(max*precision))) / precision
}

func FastRand(min, max int) int {
	if min > max {
		return 0
	}

	return rand.IntN(max-min) + min // nolint:gosec
}
