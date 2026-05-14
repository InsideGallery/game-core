# Relations

Import path: `github.com/InsideGallery/game-core/engine/relations`

Package `relations` provides parent/child relationship contracts and a thread-safe relation component. Parents and
children are ECS entities; the component tracks parent IDs by entity type and child lists by entity type.

Key exports:

- `Child` describes an entity that can reference parents and construct or destroy parent links.
- `Parent` describes an entity that can attach, detach, and list child entities.
- `RelationComponent` stores parent IDs and attached child entities behind a mutex.
- `NewRelationComponent` creates a relation component from an initial parent map and an `ignoreZeroIDs` setting.
- `SetParent`, `RemParent`, `GetParentID`, and `GetParents` manage parent references.
- `Parent` resolves a parent entity by type and returns `ErrNotFoundParent` when no parent is available.
- `ConstructChild` and `DestroyChild` attach or detach a child from all configured parents.
- `Attach`, `Detach`, `GetChildren`, and `GetAllChildren` manage child lists.
