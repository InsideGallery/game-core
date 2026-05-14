package mathutils

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

// ---------------------------------------------------------------------------
// SumValues
// ---------------------------------------------------------------------------

func TestSumValuesMultiple(t *testing.T) {
	result := SumValues([]float64{1, 2, 3, 4, 5})
	testutils.Equal(t, result, 15.0)
}

func TestSumValuesEmpty(t *testing.T) {
	result := SumValues([]float64{})
	testutils.Equal(t, result, 0.0)
}

func TestSumValuesSingle(t *testing.T) {
	result := SumValues([]float64{42.5})
	testutils.Equal(t, result, 42.5)
}

func TestSumValuesNegative(t *testing.T) {
	result := SumValues([]float64{-1, -2, -3})
	testutils.Equal(t, result, -6.0)
}

func TestSumValuesMixed(t *testing.T) {
	result := SumValues([]float64{10, -5, 3, -8})
	testutils.Equal(t, result, 0.0)
}

func TestSumValuesNil(t *testing.T) {
	result := SumValues(nil)
	testutils.Equal(t, result, 0.0)
}

// ---------------------------------------------------------------------------
// Max – edge cases not in existing tests
// ---------------------------------------------------------------------------

func TestMaxSingleValue(t *testing.T) {
	testutils.Equal(t, Max(7.0), 7.0)
}

func TestMaxEmpty(t *testing.T) {
	testutils.Equal(t, Max(), 0.0)
}

func TestMaxAllEqual(t *testing.T) {
	testutils.Equal(t, Max(3.0, 3.0, 3.0), 3.0)
}

func TestMaxNegativeValues(t *testing.T) {
	testutils.Equal(t, Max(-10.0, -5.0, -1.0), -1.0)
}

func TestMaxTwoValues(t *testing.T) {
	testutils.Equal(t, Max(1.0, 2.0), 2.0)
}

func TestMaxLargeSpread(t *testing.T) {
	testutils.Equal(t, Max(-1000.0, 0.0, 1000.0), 1000.0)
}

// ---------------------------------------------------------------------------
// Min – edge cases not in existing tests
// ---------------------------------------------------------------------------

func TestMinSingleValue(t *testing.T) {
	testutils.Equal(t, Min(7.0), 7.0)
}

func TestMinEmpty(t *testing.T) {
	testutils.Equal(t, Min(), 0.0)
}

func TestMinAllEqual(t *testing.T) {
	testutils.Equal(t, Min(3.0, 3.0, 3.0), 3.0)
}

func TestMinPositiveValues(t *testing.T) {
	testutils.Equal(t, Min(10.0, 5.0, 1.0), 1.0)
}

func TestMinTwoValues(t *testing.T) {
	testutils.Equal(t, Min(1.0, 2.0), 1.0)
}

func TestMinLargeSpread(t *testing.T) {
	testutils.Equal(t, Min(-1000.0, 0.0, 1000.0), -1000.0)
}
