// Package text loads fonts from an fs.FS and draws word-wrapped text. It is
// taken from the tutmos client text package. Wrapping is ebiten-free so it is
// unit-testable.
package text

import "strings"

// wrapWith breaks s into lines no wider than maxWidth according to measure.
// It breaks on spaces, splits a single over-long word by runes, and keeps
// explicit newlines.
func wrapWith(measure func(string) float64, s string, maxWidth float64) []string {
	var lines []string

	for _, paragraph := range strings.Split(s, "\n") {
		lines = append(lines, wrapParagraph(measure, paragraph, maxWidth)...)
	}

	return lines
}

func wrapParagraph(measure func(string) float64, paragraph string, maxWidth float64) []string {
	var (
		lines   []string
		current string
	)

	for _, word := range strings.Split(paragraph, " ") {
		candidate := word
		if current != "" {
			candidate = current + " " + word
		}

		if measure(candidate) <= maxWidth {
			current = candidate

			continue
		}

		if current != "" {
			lines = append(lines, current)
		}

		if measure(word) <= maxWidth {
			current = word

			continue
		}

		// Word longer than the width: split by runes.
		var piece []rune

		for _, r := range word {
			if len(piece) > 0 && measure(string(piece)+string(r)) > maxWidth {
				lines = append(lines, string(piece))
				piece = piece[:0]
			}

			piece = append(piece, r)
		}

		current = string(piece)
	}

	return append(lines, current)
}
