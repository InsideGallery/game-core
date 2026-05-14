# Engine

Import path: `github.com/InsideGallery/game-core/engine`

Package `engine` defines small shared interfaces and components used by game projects. It includes the top-level
`Game` and `Player` contracts plus a thread-safe attribute map for mutable entity state.

Key exports:

- `Game` describes a game object that can be initialized and ticked with a context.
- `Player` combines an ECS entity with the communications contract from `engine/communications`.
- `Attributes` stores arbitrary key/value attributes behind a mutex.
- `NewAttributes` creates an initialized attribute store.
- `SetAttribute`, `UpdateAttribute`, `GetAttribute`, and `RemoveAttribute` manage raw attribute values.
- Typed getters such as `GetInt`, `GetString`, `GetBool`, `GetUint32`, `GetUint64`, `GetUint8`, and `GetFloat64`
  return zero values for missing attributes and type-assert existing values.
