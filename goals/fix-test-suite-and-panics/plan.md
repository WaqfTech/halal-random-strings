# Plan: Fix Test Suite Panics and Assertions

## Execution Steps
- [x] **Phase 1 — Eliminate Panics**: Fix the out-of-bounds slice index `-1` in `TestIncludeRandomNumber` by properly handling empty separators and empty result strings.
- [x] **Phase 2 — Pass-by-Reference / Option Normalization**: Ensure tests assert against the actual default separator (`-`) rather than reading uninitialized local copies.
- [x] **Phase 3 — Comprehensive Test Coverage**: Update tests for `TestWordsLoad`, `TestIsSafe`, `TestWordCount`, and `TestCategories` to test against the embedded engine and sanitized vocabulary.
- [x] **Phase 4 — Race Detection**: Run `go test -v -race ./...` and ensure all tests pass with zero failures and zero race conditions.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/fix-test-suite-and-panics/goal.md)` using a `fix:` subject, then run `sila goals`.
