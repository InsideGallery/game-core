# Cards

Import path: `github.com/InsideGallery/game-core/cards`

Package `cards` provides utilities for a standard 52-card deck and binary poker hand evaluation. Cards and
combinations are represented as integer bit masks, with helper functions for converting between IDs and short
string names such as `Ah`, `Tc`, and `straightFlush`.

Key exports:

- `Deck` stores the remaining cards in a deck and is guarded by a mutex.
- `NewDeck` creates a full 52-card deck.
- `Deck.Shuffle`, `Deck.Top`, `Deck.GetTopCards`, and `Deck.GetRandomCard` mutate the deck while drawing cards.
- `ErrNoCardsInDeck` is returned when a draw cannot be satisfied.
- `BinaryEvaluation.Execute` evaluates a card slice and returns a `Combination`.
- `Combination.Score` combines the combination, weight, and kicker weight into a comparable score.
- `GetCardID`, `GetCardName`, `GetCardsIDs`, and `GetCardsNames` convert card IDs and names.
- `GetCombinationID` and `GetCombinationName` convert poker combination IDs and names.
