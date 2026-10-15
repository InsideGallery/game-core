package text

import (
	"bytes"
	"fmt"
	"image/color"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// lineSpacingFactor is the line spacing used by Measure, times the face size.
const lineSpacingFactor = 1.4

type faceKey struct {
	size float64
	bold bool
}

// Fonts holds the regular and bold font sources and a per-size face cache.
type Fonts struct {
	regular *text.GoTextFaceSource
	bold    *text.GoTextFaceSource
	mu      sync.Mutex
	faces   map[faceKey]text.Face
}

// Load reads the regular and bold fonts from fsys.
func Load(fsys fs.FS, regularPath, boldPath string) (*Fonts, error) {
	regular, err := loadSource(fsys, regularPath)
	if err != nil {
		return nil, err
	}

	bold, err := loadSource(fsys, boldPath)
	if err != nil {
		return nil, err
	}

	return &Fonts{regular: regular, bold: bold, faces: map[faceKey]text.Face{}}, nil
}

func loadSource(fsys fs.FS, path string) (*text.GoTextFaceSource, error) {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read font %s: %w", path, err)
	}

	source, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse font %s: %w", path, err)
	}

	return source, nil
}

// Face returns the cached face for size and weight.
func (f *Fonts) Face(size float64, bold bool) text.Face { //nolint:ireturn // ebiten face interface
	f.mu.Lock()
	defer f.mu.Unlock()

	key := faceKey{size: size, bold: bold}
	if face, ok := f.faces[key]; ok {
		return face
	}

	source := f.regular
	if bold {
		source = f.bold
	}

	face := &text.GoTextFace{Source: source, Size: size}
	f.faces[key] = face

	return face
}

// Wrap word-wraps s to maxWidth using face's advance.
func Wrap(face text.Face, s string, maxWidth float64) []string {
	return wrapWith(func(t string) float64 {
		return text.Advance(t, face) //nolint:staticcheck // whole candidate lines, not prefixes
	}, s, maxWidth)
}

// Measure returns the width and height of s drawn with face. Line spacing is
// the face size times 1.4.
func Measure(face text.Face, s string) (w, h float64) {
	return text.Measure(s, face, faceSize(face)*lineSpacingFactor)
}

func faceSize(face text.Face) float64 {
	if g, ok := face.(*text.GoTextFace); ok {
		return g.Size
	}

	m := face.Metrics()

	return m.HAscent + m.HDescent
}

// Draw draws lines at x, y, advancing lineHeight per line.
func Draw(dst *ebiten.Image, face text.Face, lines []string, x, y, lineHeight float64, clr color.Color) {
	for i, line := range lines {
		op := &text.DrawOptions{}
		op.GeoM.Translate(x, y+float64(i)*lineHeight)
		op.ColorScale.ScaleWithColor(clr)
		text.Draw(dst, line, face, op)
	}
}
