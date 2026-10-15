# client/input

Import path: `github.com/InsideGallery/game-core/client/input`

One pointer for touch and mouse that reports taps and drags. Source is the first touch when any exists, else the left mouse button.

- `NewPointer(slop)`: slop in logical pixels; callers pass 10 times the device scale.
- `Pointer.Update()`: call once per frame.
- `Pointer.Tap()`: position and true on the frame a press ends within slop of where it started.
- `Pointer.Drag()`: movement since the last frame while a press has moved beyond slop.
- `Pointer.Wheel()`: vertical mouse wheel delta.

The state machine lives in the unexported `tracker`, tested without ebiten input. Out of scope: multi-touch, gestures, keyboard.
