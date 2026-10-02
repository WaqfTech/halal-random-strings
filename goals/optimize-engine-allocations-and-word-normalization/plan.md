# Plan: Optimize Generator Hot-Path and Pre-normalize Dictionary Words

## Execution Steps
- [ ] **Phase 1 — Benchmark Baseline Profiling**: Capture CPU and memory profiles using `go test -bench=. -benchmem -cpuprofile=cpu.pprof -memprofile=mem.pprof`.
- [ ] **Phase 2 — In-Memory Dictionary Pre-normalization**:
  - Add an internal pre-normalized category index to `Engine` where words are stored trimmed, lowercased, and pre-tokenized.
  - Store token length alongside words so word count calculation requires no `strings.Split`.
- [ ] **Phase 3 — Generator Loop Allocation Reductions**:
  - Replace quadratic string concatenation with a pre-sized `strings.Builder`.
  - Replace `big.NewInt` with lightweight uniform integer generator.
- [ ] **Phase 4 — PRNG Optimization**:
  - Use a sync.Pool or thread-safe source when `opts.Seed == 0` rather than allocating a new `rand.Rand` struct on every invocation.
- [ ] **Phase 5 — Benchmark Verification**: Run `go test -bench=. -benchmem` to verify allocation reduction and latency gains.
- [ ] **Phase 6 — Attribution Commit**: Commit with `(goals/optimize-engine-allocations-and-word-normalization/goal.md)` using a `perf:` subject, then run `sila goals`.
