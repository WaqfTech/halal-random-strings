# Goal: Segregate Holy Sanctuaries and Fix Disrespectful Animal Pairings

## Goal Description
Segregate holy Islamic sanctuaries (Kaaba, Masjid, Madina) from general places in words.json. Restrict Rule 16 ({animals}-{nouns_places}) to prevent pairing impure/derogatory animals with sacred sites.

## Dependencies & Execution Order
- **Mode**: Independent (disjoint, parallelizable)
- **Depends On**: none
- **Sequence**: -
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-sanctuaries-and-animal-rules/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-sanctuaries-and-animal-rules/plan.md`](plan.md)

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...` or stack equivalent).
3. Zero architectural regressions; definitions of done satisfied.
