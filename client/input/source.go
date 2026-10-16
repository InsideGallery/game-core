package input

import "github.com/hajimehoshi/ebiten/v2"

// Sample is one frame of raw input.
type Sample struct {
	X, Y    float64
	Pressed bool
	Touch   bool
	Wheel   float64
}

// Source supplies one Sample per frame. Tests replace it with a fake.
type Source interface {
	Read() Sample
}

type ebitenSource struct {
	ids []ebiten.TouchID
}

// EbitenSource reads the first touch, else the cursor and the left mouse button.
func EbitenSource() Source {
	return &ebitenSource{}
}

// Read returns the current ebiten input.
func (s *ebitenSource) Read() Sample {
	_, wheel := ebiten.Wheel()

	s.ids = ebiten.AppendTouchIDs(s.ids[:0])
	if len(s.ids) > 0 {
		x, y := ebiten.TouchPosition(s.ids[0])

		return Sample{X: float64(x), Y: float64(y), Pressed: true, Touch: true, Wheel: wheel}
	}

	x, y := ebiten.CursorPosition()

	return Sample{
		X: float64(x), Y: float64(y),
		Pressed: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		Wheel:   wheel,
	}
}

// SetHand shows the hand cursor when on, else the default arrow. It changes nothing on a phone.
func SetHand(on bool) {
	if on {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)

		return
	}

	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
}
