# Facts: Fix Delimiter Safety Bypass and Suffix Seed Determinism

## Architectural Invariants & Constraints
- The safety filter `isSafe` must remain zero-tolerance: any blocked term must be rejected regardless of formatting, casing, or punctuation delimiters.
- Deterministic reproducibility: Given a specific seed `S`, two executions of `GenerateWithOptionsE` with identical options must yield byte-for-byte identical output strings.
- Panic-freedom: Public library methods (`Generate`, `GenerateN`, `GenerateWithOptions`) must never panic on invalid inputs or runtime errors.

## File & Interface Contracts
- `halal-random-strings.go`: `isSafeLocked`, `GenerateWithOptionsE`, `Generate`.
- `halal-random-strings_test.go`: Add `TestDelimiterBypassSafety` and `TestDeterministicSeedReproducibility`.
