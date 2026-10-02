# Plan: Embed words.json and Provide Thread-Safe Generator Engine

## Execution Steps
- [x] **Phase 1 — Embedding & Model**: Embed `words.json` using `//go:embed` in `halal-random-strings.go`. Define a thread-safe `Engine` struct with `sync.RWMutex`.
- [x] **Phase 2 — Default Initialization**: Create `NewDefaultEngine()` and initialize a package-level default engine so that calls to `Generate()`, `GenerateN()`, and `GenerateWithOptions()` work immediately without requiring manual `LoadWords()` calls.
- [x] **Phase 3 — CLI & Package Export**: Update `cmd/halal-random-strings/main.go` to use the embedded engine, ensuring `halal-random-strings` runs portably without external file dependencies.
- [x] **Phase 4 — Attribution Commit**: Commit with `(goals/fix-words-embedding-and-init/goal.md)` using a `fix:` subject, then run `sila goals`.
