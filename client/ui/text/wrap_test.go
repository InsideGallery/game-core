package text

import (
	"reflect"
	"strings"
	"testing"
)

func perRune(s string) float64 { return float64(len([]rune(s))) }

func TestWrapWith(t *testing.T) {
	got := wrapWith(perRune, "one two three four", 10)
	if want := []string{"one two", "three four"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}

	got = wrapWith(perRune, strings.Repeat("x", 25), 10)
	if len(got) != 3 || len([]rune(got[0])) != 10 || len([]rune(got[1])) != 10 || len([]rune(got[2])) != 5 {
		t.Fatalf("long word: got %q", got)
	}
}

func TestWrapWithNewlinesAndEmpty(t *testing.T) {
	if got := wrapWith(perRune, "a\nb", 10); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Fatalf("newlines: got %q", got)
	}

	if got := wrapWith(perRune, "", 10); !reflect.DeepEqual(got, []string{""}) {
		t.Fatalf("empty: got %q", got)
	}
}
