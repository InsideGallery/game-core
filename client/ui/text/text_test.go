package text

import (
	"strings"
	"testing"
	"testing/fstest"

	"golang.org/x/image/font/gofont/goregular"
)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"r.ttf": {Data: goregular.TTF},
		"b.ttf": {Data: goregular.TTF},
	}
}

func TestFaceCache(t *testing.T) {
	f, err := Load(testFS(), "r.ttf", "b.ttf")
	if err != nil {
		t.Fatal(err)
	}

	if f.Face(16, false) != f.Face(16, false) {
		t.Fatal("face not cached")
	}

	if f.Face(16, true) == f.Face(16, false) {
		t.Fatal("bold and regular share a face")
	}
}

func TestLoadMissingPath(t *testing.T) {
	_, err := Load(testFS(), "r.ttf", "missing/bold.ttf")
	if err == nil || !strings.Contains(err.Error(), "missing/bold.ttf") {
		t.Fatalf("want error with path, got %v", err)
	}
}

func TestMeasure(t *testing.T) {
	f, err := Load(testFS(), "r.ttf", "b.ttf")
	if err != nil {
		t.Fatal(err)
	}

	w1, h1 := Measure(f.Face(16, false), "hello")
	_, h2 := Measure(f.Face(16, false), "hello\nworld")

	if w1 <= 0 || h1 <= 0 || h2 <= h1 {
		t.Fatalf("unexpected measure %v %v %v", w1, h1, h2)
	}
}
