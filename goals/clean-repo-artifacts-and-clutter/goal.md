# Goal: Purge Committed Binary, 2bedeleted Folder, and Orphaned Files

## Goal Description
Remove the 3MB compiled binary 'halal' from git tracking, delete the leftover 2bedeleted/ folder, remove or wire cities.txt, and update .gitignore to prevent binary check-ins.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/clean-repo-artifacts-and-clutter/facts.md`](facts.md)
- **Execution Plan**: [`goals/clean-repo-artifacts-and-clutter/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
