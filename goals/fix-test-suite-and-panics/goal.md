# Goal: Fix Test Suite Panics and Assertions

## Goal Description
Fix the index out of range panic in TestIncludeRandomNumber, remove pass-by-value bugs in test setup, and update all tests to verify embedded vocabulary, word count constraints, and category filters.

## Dependencies & Execution Order
- **Mode**: Dependent (requires prerequisite goals)
- **Depends On**: fix-words-embedding-and-init, fix-scunthorpe-and-blocked-list
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-test-suite-and-panics/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-test-suite-and-panics/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
