# Goal: Optimize Generator Hot-Path and Pre-normalize Dictionary Words

## Goal Description
Drastically reduce memory allocations and latency on the generator hot path:
1. **Pre-normalize Words at Load Time**: Currently, `normalizeWord` runs strings replacements, lowercasing, and trim loops dynamically for every candidate word on every generation attempt. Pre-normalize all dictionary words during engine initialization (`NewEngine`, `NewDefaultEngine`, `LoadWords`) and pre-compute word token counts.
2. **Eliminate Quadratic String Building & Repeated Splitting**: Replace repetitive `strings.Split(currentOutput)` and string concatenations (`+`) with single-pass buffer or slice operations.
3. **Eliminate Heap Allocations for 4-Digit Numbers**: Avoid creating `big.NewInt(9000)` on every generated number; read raw bytes from `crypto/rand` or use fast integer bounds.
4. **Pool or Optimize PRNG Instantiation**: Eliminate re-instantiating `rand.New(rand.NewSource(...))` on every call to `GenerateWithOptionsE`.
5. **Target Benchmark**: Reduce allocation from 6,400 B/op and 40 allocs/op to < 500 B/op and < 5 allocs/op, achieving a > 10x throughput improvement.

## Dependencies & Execution Order
- **Mode**: Dependent 🔗
- **Depends On**: fix-delimiter-safety-and-seed-determinism
- **Sequence**: 18
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/optimize-engine-allocations-and-word-normalization/facts.md`](facts.md)
- **Execution Plan**: [`goals/optimize-engine-allocations-and-word-normalization/plan.md`](plan.md)

## Done Condition
1. All dictionary words are pre-normalized and token counts pre-calculated upon engine initialization.
2. `BenchmarkGenerate` demonstrates < 1,000 ns/op and < 10 allocs/op (targeting < 500 B/op).
3. Zero semantic regressions in generated output format and separator support.
4. All tests pass with `go test -race ./...`.
