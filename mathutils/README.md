# Math Utils

Import path: `github.com/InsideGallery/game-core/mathutils`

Package `mathutils` provides small numeric helpers used by game-core packages.

Key exports:

- `SumValues` returns the sum of a float64 slice.
- `Max` and `Min` return the maximum or minimum value from variadic float64
  inputs. Empty input returns zero.
- `RoundWithPrecision` rounds a value to the requested precision.
- `ApproximatelyEqual` compares two float64 values using a very small epsilon.
