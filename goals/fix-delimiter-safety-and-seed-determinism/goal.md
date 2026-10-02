# Goal: Fix Delimiter Safety Bypass and Suffix Seed Determinism

## Goal Description
Resolve critical bugs in safety filtering, PRNG determinism, and error handling:
1. **Safety Filter Delimiter Bypass**: Fix `isSafeLocked` so that non-hyphen custom separators (such as `.`, `/`, `_`, or any arbitrary delimiter) are properly normalized before tokenization, preventing blocked terms (e.g. `wine.beer.pork`) from bypassing the filter.
2. **Deterministic Seed Suffix**: When `opts.Seed != 0`, ensure the numeric suffix is derived deterministically from the seeded PRNG (`src`), fulfilling the CLI promise `(Output will be consistent for seed 12345)`. Only invoke `crypto/rand` when `opts.Seed == 0`.
3. **Panic Protection in `Engine.Generate()`**: Prevent unchecked slice indexing `[0]` on empty or nil results when generation fails or encounters errors. Return a safe error or empty fallback without panicking.
4. **Category Subset Matching**: Ensure rule matching in `GenerateWithOptionsE` does not select single-category rules when multiple categories were explicitly requested.

## Dependencies & Execution Order
- **Mode**: Independent ⚡
- **Depends On**: -
- **Sequence**: 16
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-delimiter-safety-and-seed-determinism/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-delimiter-safety-and-seed-determinism/plan.md`](plan.md)

## Done Condition
1. `isSafe("hello.wine.world")` and any string with custom separators correctly rejects blocked terms.
2. Generating strings with the same non-zero seed produces identical word sequences AND identical numeric suffixes across runs.
3. `Engine.Generate()` never panics even if generation fails or retry limit is exhausted.
4. All tests pass with `go test -race ./...`.
