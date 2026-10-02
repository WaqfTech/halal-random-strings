# Goal: Harden D1 Population Script and Eliminate CI Blocking Prompts

## Goal Description
Remove interactive input() prompts from populate_d1.py to fix CI/CD execution, add INSERT OR IGNORE, reduce batch sizes to safe CLI limits, and document on-the-fly Worker generation pattern.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/harden-d1-scripts-and-architecture/facts.md`](facts.md)
- **Execution Plan**: [`goals/harden-d1-scripts-and-architecture/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
