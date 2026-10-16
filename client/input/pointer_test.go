package input

import "testing"

type fakeSource struct {
	samples []Sample
	i       int
}

func (f *fakeSource) Read() Sample {
	s := f.samples[f.i]
	f.i++

	return s
}

func TestPointerHover(t *testing.T) {
	cases := []struct {
		name   string
		sample Sample
		ok     bool
	}{
		{"mouse", Sample{X: 5, Y: 5}, true},
		{"touch", Sample{X: 5, Y: 5, Touch: true}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewPointerFrom(10, &fakeSource{samples: []Sample{c.sample}})
			p.Update()

			x, y, ok := p.Hover()
			if ok != c.ok || (ok && (x != 5 || y != 5)) {
				t.Fatalf("hover = %v,%v,%v want ok %v", x, y, ok, c.ok)
			}

			if _, _, ok := p.Press(); ok {
				t.Fatal("press must not be ok")
			}
		})
	}
}

func TestPointerPress(t *testing.T) {
	p := NewPointerFrom(10, &fakeSource{samples: []Sample{
		{X: 0, Y: 0, Pressed: true},
		{X: 4, Y: 4, Pressed: true},
		{X: 0, Y: 30, Pressed: true},
	}})
	cases := []struct {
		x, y float64
		ok   bool
	}{{0, 0, true}, {4, 4, true}, {0, 0, false}}

	for i, c := range cases {
		p.Update()

		x, y, ok := p.Press()
		if ok != c.ok || (ok && (x != c.x || y != c.y)) {
			t.Fatalf("frame %d: press = %v,%v,%v want %v,%v,%v", i, x, y, ok, c.x, c.y, c.ok)
		}
	}
}

func TestPointerTapFromSource(t *testing.T) {
	p := NewPointerFrom(10, &fakeSource{samples: []Sample{
		{X: 5, Y: 5, Pressed: true},
		{X: 5, Y: 5},
	}})

	p.Update()

	if _, _, ok := p.Tap(); ok {
		t.Fatal("tap on press frame")
	}

	p.Update()

	if x, y, ok := p.Tap(); !ok || x != 5 || y != 5 {
		t.Fatalf("tap = %v,%v,%v", x, y, ok)
	}
}
