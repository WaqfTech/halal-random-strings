# Goal: Prune False-Positive Blocked Words and Scraped Dictionary Noise

## Goal Description
Resolve dictionary self-blocking and clean up scraped dictionary artifacts:
1. **Unblock 110+ Self-Blocked Words**: The `blocked` list currently contains ordinary English words (`bed`, `ocean`, `mountain`, `world`, `love`, `angel`, `friendship`, `health`, `plant`, `song`, `perfume`, `music`, `dragon`, `ghost`, `stick`, `journey`, `wedding`, `soul`, `spirit`, `paradise`, etc.) which causes the engine to reject valid words from its own categories (`geographic_features`, `holy_sanctuaries`, `trees`, `fruits`, `islamic_inventions`).
2. **Handle Contextual `haram`**: Prevent the blocking of sacred terms like `masjid-al-haram` while still preventing illicit terms.
3. **Clean Up Scraped Dictionary Artifacts**:
   - Strip parenthetical glosses from `suffixes`: `dwelling-place` -> `dwelling`, `abode-residence` -> `abode`, `oasis-green` -> `oasis`, `moment-now` -> `moment`, `night-dark` -> `night`, `day-light` -> `day`.
   - Strip taxonomic qualifiers from `vegetables`, `fruits`, and `trees`: `Spinach-Leaf` -> `Spinach`, `Turnip-Root` -> `Turnip`, `Carrot-Root` -> `Carrot`, `Orange-Citrus` -> `Orange`, `Blackberry-Bush` -> `Blackberry`.
4. **Deduplicate Categories and Blocked List**: Remove duplicate entries (`Lahmacun` in `arabic_food`, `Badger` and `Marten` in `animals`, `Cupola` and `Loggia` in `architectural_elements`, and duplicates in `blocked`).

## Dependencies & Execution Order
- **Mode**: Independent ⚡
- **Depends On**: -
- **Sequence**: 17
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/prune-blocked-list-and-resolve-self-blocking/facts.md`](facts.md)
- **Execution Plan**: [`goals/prune-blocked-list-and-resolve-self-blocking/plan.md`](plan.md)

## Done Condition
1. All 36 categories in `words.json` have 0 words blocked by their own `blocked` list.
2. Scraped parenthetical suffixes and taxonomic tags are stripped to single clean words.
3. Category and blocked duplicates are completely removed.
4. All tests pass with `go test -race ./...`.
