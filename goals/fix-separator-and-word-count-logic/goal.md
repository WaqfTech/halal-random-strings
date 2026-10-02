# Goal: Fix Custom Separator Formatting and Word Count Boundaries

## Goal Description
Ensure custom --sep applies to template joints and within compound words without producing hybrid separators. Count words accurately after separator normalization so multi-word entries do not exceed MaxWords.

## Dependencies & Execution Order
- **Mode**: Dependent (requires prerequisite goals)
- **Depends On**: fix-words-embedding-and-init
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-separator-and-word-count-logic/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-separator-and-word-count-logic/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
