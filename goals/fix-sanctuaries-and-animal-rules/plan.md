# Plan: Segregate Holy Sanctuaries and Fix Disrespectful Animal Pairings

## Execution Steps
- [x] **Phase 1 — Separate Sanctuaries**: Extract sacred sites (`kaaba`, `masjid`, `madina`, `al-aqsa`, `minbar`, `mihrab`) from `nouns_places` in `words.json` into a separate `holy_sanctuaries` category.
- [x] **Phase 2 — Guard Rule 16**: Modify Rule 16 (`{animals}-{nouns_places}`) so that animals only combine with geographic features (`desert`, `oasis`, `mountain`, `river`, `valley`).
- [x] **Phase 3 — Audit Animal Vocabulary**: Review `animals` category to prevent ritual impurities or derogatory associations in family-friendly strings.
- [x] **Phase 4 — Attribution Commit**: Commit with `(goals/fix-sanctuaries-and-animal-rules/goal.md)` using a `fix:` subject, then run `sila goals`.
