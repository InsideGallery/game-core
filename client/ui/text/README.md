# client/ui/text

Import path: `github.com/InsideGallery/game-core/client/ui/text`

Font loading, word wrap, measure and draw for Ebitengine. Taken from the tutmos client text package; font paths are parameters.

- `Load(fsys, regularPath, boldPath)`: reads two fonts from an `fs.FS`; errors name the path.
- `Fonts.Face(size, bold)`: cached face per size and weight, safe for concurrent use.
- `Wrap(face, s, maxWidth)`: breaks on spaces, splits over-long words by runes, keeps newlines; empty input gives one empty line.
- `Measure(face, s)`: width and height, line spacing is the face size times 1.4.
- `Draw(dst, face, lines, x, y, lineHeight, clr)`: draws the lines.
