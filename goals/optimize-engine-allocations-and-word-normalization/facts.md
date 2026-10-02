# Facts: Optimize Generator Hot-Path and Pre-normalize Dictionary Words

## Architectural Invariants & Constraints
- Memory footprint: Minimizing GC pressure is critical for embedding in high-throughput microservices and Cloudflare Workers (WebAssembly).
- Thread-safety: Any shared caching, PRNG pooling, or lookup indexes must be strictly concurrent-read safe under `sync.RWMutex`.
- Output parity: Optimizations must not change the distribution of selected words or format of generated strings.

## File & Interface Contracts
- `halal-random-strings.go`: `Engine` struct, `NewEngine`, `GenerateWithOptionsE`, `normalizeWord`.
- `halal-random-strings_test.go`: `BenchmarkGenerate`, `BenchmarkGenerateWithOptions`.
