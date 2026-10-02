# Plan: Fix Delimiter Safety Bypass and Suffix Seed Determinism

## Execution Steps
- [ ] **Phase 1 — Universal Separator Normalization in `isSafeLocked`**:
  - Replace any non-alphanumeric punctuation rune or delimiter (or split on `opts.Sep` and non-alphanumerics) so that `tokens` correctly isolates words regardless of whether `sep` is `-`, `_`, `.`, `/`, or custom.
  - Verify that `isSafe("safe.wine.world")` returns `false`.
- [ ] **Phase 2 — Seed-Aware Numeric Suffix Generation**:
  - In `GenerateWithOptionsE`, check if `opts.Seed != 0`.
  - If seeded, generate the 4-digit number `[1000, 9999]` using `src.Intn(9000) + 1000`.
  - If unseeded (`opts.Seed == 0`), use `crypto/rand` for cryptographic entropy.
- [ ] **Phase 3 — Nil-Safe Indexing in `Engine.Generate()`**:
  - Update `Engine.Generate()` to check `len(r) > 0` before indexing `[0]`.
  - Return an empty string or bubble up errors gracefully through an error-returning variant.
- [ ] **Phase 4 — Fix Multi-Category Rule Filtering**:
  - Ensure that when multiple categories are provided, compound rules containing those categories are prioritized and single-category rules are only used if explicitly requested.
- [ ] **Phase 5 — Automated Unit Tests**: Add tests verifying custom separator safety and exact seed reproducibility across runs.
- [ ] **Phase 6 — Attribution Commit**: Commit with `(goals/fix-delimiter-safety-and-seed-determinism/goal.md)` using a `fix:` subject, then run `sila goals`.
