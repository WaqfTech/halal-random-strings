# Plan: Sanitize Theological Terms and Prevent Disrespectful Word Combinations

## Execution Steps
- [x] **Phase 1 — Audit & Purge Malevolent Entities**: Remove `dajjal`, `jahannam`, `fitna`, and `yajuj-majuj` from `nouns_concepts` in `words.json`.
- [x] **Phase 2 — Purge Sacrilegious Rules**:
  - Remove Rule 11 (`['islamic_golden_age_scholars', 'vegetables']`).
  - Remove Rule 43 (`['animals', 'adjectives']`) or segregate adjectives into physical descriptors vs moral virtues/prophetic titles.
  - Remove Rules 41 & 42 (`['colors_arabic', 'muslim_names_male']` and `['colors_arabic', 'muslim_names_female']`).
- [x] **Phase 3 — Complete Sacred Entity Segregation**:
  - Remove sacred sites (`kaba-structure`, `hajar-al-aswad`, `maqam-ibrahim`, `hijr-isma'il`, `masjid-al-haram`, `masjid-al-nabawi`, `zamzam-well`, `safa-marwah-mounts`) from `nouns_concepts` (ensure they are only in `holy_sanctuaries`).
  - Segregate Quranic terms (`quran`, `ayats`, `surahs`) into dedicated respectful rules or category.
- [x] **Phase 4 — Sanitize Divine Names & Culinary Transliterations**:
  - Ensure male names with exclusive Divine Attributes are prefixed with `abd-` (e.g. `abd-al-qadir`, `abd-al-ghani`, `abd-al-wahid`).
  - Correct `fattah` in `saudi_food` and `arabic_food` to `fatteh` or `fatta` to avoid collision with *Al-Fattah*.
  - Remove `jannah-end` from `suffixes`.
- [x] **Phase 5 — Automated Test Verification**: Add regression tests verifying no scholar pairs with vegetables, no animal pairs with prophetic titles, and no malevolent concepts exist.
- [x] **Phase 6 — Attribution Commit**: Commit with `(goals/fix-theological-sanitation-and-combinatorial-rules/goal.md)` using a `fix:` subject, then run `sila goals`.
