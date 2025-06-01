package shapes

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestPointOperation(t *testing.T) {
	p := NewPoint(0, 0, 0)
	p2 := NewPoint(1, 0, 0)
	p3 := NewPoint(0, 1, 0)

	mp := p.Copy().Add(p2).Add(p3)

	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	testutils.Equal(t, p.Add(p2).Coordinates(), [3]float64{1, 0, 0})
	testutils.Equal(t, p.Add(p3).Coordinates(), [3]float64{0, 1, 0})

	mp = mp.Scale(2)
	testutils.Equal(t, mp.Coordinates(), [3]float64{2, 2, 0})
	mp = mp.Divide(NewPoint(2, 2, 0))
	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	mp = mp.Divide(NewPoint(0, 0, 0))
	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	mp = mp.Increase(2)
	testutils.Equal(t, mp.Coordinates(), [3]float64{3, 3, 2})
	mp = mp.Invert()
	testutils.Equal(t, mp.Coordinates(), [3]float64{-3, -3, -2})
	mp = mp.Invert()
	mp = mp.Decrease(2)
	testutils.Equal(t, mp.Coordinates(), [3]float64{1, 1, 0})
	nmp := mp.Copy()
	nmp = nmp.Scale(20)
	testutils.Equal(t, nmp.Coordinates(), [3]float64{20, 20, 0})
	testutils.Equal(t, nmp.Cross(p2).Coordinates(), [3]float64{0, 0, -20})
	mp = mp.Scale(2)
	mp = mp.Add(NewPoint(0, 0, 1))
	testutils.Equal(t, mp.Coordinates(), [3]float64{2, 2, 1})
	p4 := mp.Copy()
	testutils.Equal(t, p4.Decrease(2).Coordinates(), [3]float64{0, 0, -1})
	testutils.Equal(t, p4.Subtract(NewPoint(1, 0, 0)).Coordinates(), [3]float64{1, 2, 1})
	testutils.Equal(t, p4.Multiply(NewPoint(1, 1, 2)).Coordinates(), [3]float64{2, 2, 2})
	testutils.Equal(t, p4.Divide(NewPoint(2, 1, 0)).Coordinates(), [3]float64{1, 2, 1})
	testutils.Equal(t, p4.Increase(2).Coordinates(), [3]float64{4, 4, 3})
	testutils.Equal(t, p4.Decrease(1).Coordinates(), [3]float64{1, 1, 0})
	testutils.Equal(t, p4.Invert().Coordinates(), [3]float64{-2, -2, -1})
	testutils.Equal(t, p4.Cross(p2).Coordinates(), [3]float64{0, 1, -2})
}

/*
BenchmarkNonModifiablePoint-4                           20000000                60.0 ns/op            32 B/op          1 allocs/op
BenchmarkNonModifiablePointManyUpdates-4                10000000               183 ns/op              32 B/op          1 allocs/op
*/
var (
	globalVector Point
)

func BenchmarkNonModifiablePoint(b *testing.B) {
	p := NewPoint(1, 1, 1)
	p1 := NewPoint(0, 1, 0)
	for i := 0; i < b.N; i++ {
		p = p.Add(p1).Scale(1.1)
	}
	globalVector = p
}

func BenchmarkNonModifiablePointManyUpdates(b *testing.B) {
	p := NewPoint(1, 1, 1)
	p1 := NewPoint(0, 1, 0)
	for i := 0; i < b.N; i++ {
		p = p.Add(p1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1).Scale(1.1)
	}
	globalVector = p
}
