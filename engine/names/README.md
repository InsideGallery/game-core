# Names

Import path: `github.com/InsideGallery/game-core/engine/names`

Package `names` provides a thread-safe component for storing and updating an entity name.

Key exports:

- `NameComponent` stores one name string behind a mutex.
- `NewNameComponent` creates a component with an initial name.
- `SetName` replaces the stored name.
- `GetName` returns the stored name.
