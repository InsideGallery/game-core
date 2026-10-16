# client/input

Import path: `github.com/InsideGallery/game-core/client/input`

One pointer for touch and mouse that reports taps and drags. Source is the first touch when any exists, else the left mouse button.

- `NewPointer(slop)`: slop in logical pixels; callers pass 10 times the device scale.
- `Pointer.Update()`: call once per frame.
- `Pointer.Tap()`: position and true on the frame a press ends within slop of where it started.
- `Pointer.Drag()`: movement since the last frame while a press has moved beyond slop.
- `Pointer.Wheel()`: vertical mouse wheel delta.
- `Pointer.Hover()`: mouse position and true while no button is down; never true for a touch.
- `Pointer.Press()`: current position and true while a press is down and has not become a drag.
- `NewPointerFrom(slop, src)`: a Pointer on a replaceable `Source`; `NewPointer` uses `EbitenSource()`.
- `Source`, `Sample`: one frame of raw input (`X, Y, Pressed, Touch, Wheel`); tests pass a fake.
- `SetHand(on)`: hand cursor when on, default arrow when off; no effect on a phone. Call once per `Update`.

The state machine lives in the unexported `tracker`, tested without ebiten input. Out of scope: multi-touch, gestures, keyboard.
