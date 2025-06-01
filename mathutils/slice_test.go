package mathutils

import (
	"math"
	"reflect"
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestMin(t *testing.T) {
	r := []float64{1, 0.5, 0.0000001, -0.000001, 32.00001}
	v := Min(r...)
	testutils.Equal(t, v, -0.000001)

	r = []float64{1, 0.5, 0.0000001, -0.000001, 32.00001, math.NaN(), math.Inf(-1)}
	v = Min(r...)
	if !reflect.DeepEqual(v, math.Inf(-1)) {
		t.Fatalf("unexpected value %f != inf-", v)
	}

	r = []float64{0, math.NaN()}
	v = Min(r...)
	testutils.Equal(t, v, 0.0)
}

func TestMax(t *testing.T) {
	r := []float64{-1, -0.5, 0.0000001, -0.000001, -32.00001}
	v := Max(r...)
	testutils.Equal(t, v, 0.0000001)

	r = []float64{-1, -0.5, 0.0000001, -0.000001, -32.00001, math.NaN(), math.Inf(-1)}
	v = Max(r...)
	testutils.Equal(t, v, 0.0000001)

	r = []float64{-1, -0.5, 0.0000001, -0.000001, -32.00001, math.NaN(), math.Inf(1)}
	v = Max(r...)
	if !reflect.DeepEqual(v, math.Inf(1)) {
		t.Fatalf("unexpected value %f != inf+", v)
	}

	r = []float64{0, math.NaN()}
	v = Max(r...)
	testutils.Equal(t, v, 0.0)
}

/*
BenchmarkMin-12                                 231729692                4.83 ns/op            0 B/op          0 allocs/op
BenchmarkMax-12                                 253363023                4.66 ns/op            0 B/op          0 allocs/op
*/
var (
	globalValue float64
)

func BenchmarkMin(b *testing.B) {
	r := []float64{1, 0.5, 0.0000001, -0.000001, 32.00001}
	for i := 0; i < b.N; i++ {
		globalValue = Min(r...)
	}
}

func BenchmarkMax(b *testing.B) {
	r := []float64{-1, -0.5, 0.0000001, -0.000001, -32.00001}
	for i := 0; i < b.N; i++ {
		globalValue = Max(r...)
	}
}
