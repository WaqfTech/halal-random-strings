# Plan: Fix Category Filtering and Add Single-Category Fallback Logic

## Execution Steps
- [x] **Phase 1 — Diagnose Filter Logic**: Trace rule matching where single categories lack compound rules, resulting in `len(filteredRules) == 0`.
- [x] **Phase 2 — Single-Category Dynamic Fallback**: Implement dynamic rule creation (e.g. `"{category}"`) when a user specifies categories that do not have multi-part rules, enabling `--categories fruits` or `--categories adjectives` to work as expected.
- [x] **Phase 3 — Flexible Multi-Category Matching**: Allow rules where *any* pattern category matches the selected categories when exact matches do not exist, or validate that at least one rule matches before entering retry loops.
- [x] **Phase 4 — Error Reporting**: Return explicit errors if unknown categories are provided instead of silently returning empty strings.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/fix-category-filtering-and-fallback/goal.md)` using a `fix:` subject, then run `sila goals`.
