# Goal: Correct Historical and Religious Attribution of Non-Muslim Scholars

## Goal Description
In words.json, rename muslim_scientists to islamic_golden_age_scholars or categorize polymaths so that Jewish (Maimonides), Christian (Hunayn ibn Ishaq), and Sabian (Thabit ibn Qurra) scholars are respectfully and accurately attributed.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-historical-scholar-attributions/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-historical-scholar-attributions/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
