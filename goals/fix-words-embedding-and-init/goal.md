# Goal: Embed words.json and Provide Thread-Safe Generator Engine

## Goal Description
Embed words.json using //go:embed and replace the uninitialized global var words Words with a thread-safe Engine. Fixes zero-output generation and binary portability.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-words-embedding-and-init/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-words-embedding-and-init/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
