# Goal: Eliminate Redundant Rule Filtering and Hot-Path Heap Allocations

## Goal Description
Move rule filtering outside generation loops to eliminate 1M redundant slice allocations. Replace strings.NewReplacer heap churn with strings.Builder and use O(1) blocked word lookup maps.

## Dependencies & Execution Order
- **Mode**: Dependent (requires prerequisite goals)
- **Depends On**: fix-words-embedding-and-init, fix-separator-and-word-count-logic
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/optimize-generator-performance-and-allocations/facts.md`](facts.md)
- **Execution Plan**: [`goals/optimize-generator-performance-and-allocations/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
