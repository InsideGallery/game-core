# client/ui/scroll

Import path: `github.com/InsideGallery/game-core/client/ui/scroll`

Clamped vertical scroll offset driven by drag and wheel. No Ebitengine import: callers pass input values in.

- `View{Offset, ViewportH, ContentH}`: scroll state.
- `Apply(dragDY, wheel, wheelStep)`: dragging up scrolls down; wheel moves by `wheel*wheelStep`; clamps.
- `Clamp()`: keeps `Offset` within `[0, max(0, ContentH-ViewportH)]`.
- `AtEnd(margin)`: true near the content end, used to load the next page.
- `Visible(top, height)`: true when the item (content coordinates) crosses the viewport.
- `Bar(minLen)`: `(top, length, ok)` of the scroll bar in viewport coordinates; `ok` is false when nothing scrolls. Callers pass 24 times the device scale as `minLen`; each scene draws the bar.
