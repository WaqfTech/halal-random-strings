# Plan: Eliminate Redundant Rule Filtering and Hot-Path Heap Allocations

## Execution Steps
- [x] **Phase 1 — Hoist Invariant Computations**: Move category-to-rule filtering outside the `Repeat` and retry loops so it executes once per generation request instead of up to 1,000,000 times.
- [x] **Phase 2 — Replace Heap-Heavy Replacer**: Replace per-word `fmt.Sprintf` and `strings.NewReplacer` allocations with a reusable `strings.Builder` or pre-tokenized rule representation.
- [x] **Phase 3 — O(1) Blocked Word Set**: Convert `Blocked` slice lookups into a hash map / set lookup, eliminating $O(N)$ linear scans on every generated candidate.
- [x] **Phase 4 — Benchmark Verification**: Add Go benchmark tests (`BenchmarkGenerate`) and verify significant reductions in ns/op and B/op.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/optimize-generator-performance-and-allocations/goal.md)` using a `perf:` subject, then run `sila goals`.
