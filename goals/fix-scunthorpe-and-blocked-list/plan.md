# Plan: Fix Scunthorpe Problem and Sanitize Blocked Word Filter

## Execution Steps
- [x] **Phase 1 — Token-Boundary Safety Check**: Replace naive `strings.Contains` in `isSafe()` with token-boundary matching (split by hyphens/underscores/spaces) and compare against a normalized blocked hash set.
- [x] **Phase 2 — Purge Innocent Words**: Remove everyday words from `blocked` in `words.json` (`book`, `family`, `marriage`, `baby`, `child`, `kid`, `tree`, `flower`, `river`, `sky`, `sun`, `moon`, `star`, `nature`, `science-museum`, `laboratory`, `planetarium`, `observatory`).
- [x] **Phase 3 — Remove Hallucinated `-toy` Suffixes**: Clean up the 230 spurious `-toy` entries in `blocked`.
- [x] **Phase 4 — Verify Zero False Positives**: Write unit tests proving that venerable names (`Muhammad`, `Ahmad`, `Aisha`, `Khadijah`, `Uthman ibn Affan`, `Abu Hurairah`) and virtues (`Compassion`, `Steadfastness`) are no longer blocked.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/fix-scunthorpe-and-blocked-list/goal.md)` using a `fix:` subject, then run `sila goals`.
