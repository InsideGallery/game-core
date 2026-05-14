package mathutils

import (
	"testing"

	"github.com/FrogoAI/testutils"
)

func TestRoundWithPrecision(t *testing.T) {
	cases := []struct {
		name      string
		value     float64
		precision float64
		want      float64
	}{
		{
			name:      "zero precision uses default",
			value:     1.55,
			precision: 0,
			want:      1.55,
		},
		{
			name:      "zero value",
			value:     0,
			precision: 0.1,
			want:      0,
		},
		{
			name:      "round down",
			value:     1.2345,
			precision: 0.1,
			want:      1.2,
		},
		{
			name:      "round up",
			value:     1.7654,
			precision: 0.1,
			want:      1.8,
		},
		{
			name:      "small precision",
			value:     1.8555,
			precision: 0.0001,
			want:      1.8555,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, RoundWithPrecision(tc.value, tc.precision), tc.want)
		})
	}
}

func TestApproximatelyEqual(t *testing.T) {
	cases := []struct {
		name string
		a    float64
		b    float64
		want bool
	}{
		{
			name: "equal",
			a:    1,
			b:    1,
			want: true,
		},
		{
			name: "different",
			a:    1,
			b:    2,
			want: false,
		},
		{
			name: "near zero",
			a:    0,
			b:    0,
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, ApproximatelyEqual(tc.a, tc.b), tc.want)
		})
	}
}
