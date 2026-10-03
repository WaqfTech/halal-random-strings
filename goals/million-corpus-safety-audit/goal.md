# Goal: Generate and audit at least one million real identifiers

## Goal Description
Generate fresh numbered and unnumbered million-line CLI corpora, audit every line, inspect group samples, verify auditor negative controls, and preserve reproducible receipts.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/million-corpus-safety-audit/facts.md`](facts.md)
- **Execution Plan**: [`goals/million-corpus-safety-audit/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
