# client/app

Import path: `github.com/InsideGallery/game-core/client/app`

Thin Ebitengine game-loop host. It owns the event bus and the scene manager.

- `Config`: window `Width`, `Height`, `Title`, `Setup`, and `Responsive`.
- `Responsive`: when true, `Layout` ignores `Width` and `Height` and returns the outside size times `DeviceScale()`, and `Run` makes the window resizable.
- `DeviceScale()`: the monitor device scale factor, or 1 when unknown.
- `New`, `Run`: build the app and start the loop.
