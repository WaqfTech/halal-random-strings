# Goal: Fix Category Filtering and Add Single-Category Fallback Logic

## Goal Description
Fix the rule matching logic in halal-random-strings.go so specifying any of the 29 categories without compound rules does not fail silently. Provide single-word category selection fallback.

## Dependencies & Execution Order
- **Mode**: Dependent (requires prerequisite goals)
- **Depends On**: fix-words-embedding-and-init
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-category-filtering-and-fallback/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-category-filtering-and-fallback/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
