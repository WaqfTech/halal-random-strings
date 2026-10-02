# Goal: Harmonize License Inconsistency, Update Security Policy and Docs

## Goal Description
Clarify MIT vs Waqf GPL copyleft conflict, attribute WaqfTech in the license, replace placeholder email in SECURITY.md, and update README/TODOS documentation drift.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/harmonize-licensing-and-security-policy/facts.md`](facts.md)
- **Execution Plan**: [`goals/harmonize-licensing-and-security-policy/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
