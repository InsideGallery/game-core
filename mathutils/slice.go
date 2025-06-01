package mathutils

// SumValues sum values
func SumValues(value []float64) (result float64) {
	for _, v := range value {
		result += v
	}

	return result
}

// Max return maximum value in slice
func Max(vs ...float64) float64 {
	if len(vs) == 0 {
		return 0
	}

	if len(vs) == 1 {
		return vs[0]
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// Min return minimum value in slice
func Min(vs ...float64) float64 {
	if len(vs) == 0 {
		return 0
	}

	if len(vs) == 1 {
		return vs[0]
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}
