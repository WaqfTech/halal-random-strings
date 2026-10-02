# Plan: Prune False-Positive Blocked Words and Scraped Dictionary Noise

## Execution Steps
- [x] **Phase 1 — Audit Self-Blocked Category Words**: Run automated cross-check script to find every word in `categories` that triggers `isSafe() == false`.
- [x] **Phase 2 — Prune Benign Terms from `blocked`**:
  - Remove benign everyday nouns from `blocked`: `ocean`, `mountain`, `world`, `bed`, `love`, `angel`, `friendship`, `health`, `plant`, `song`, `perfume`, `dragon`, `ghost`, `stick`, `journey`, `wedding`, `soul`, `spirit`, `paradise`.
  - Ensure genuine vulgarities, drugs, alcohol, sexual terms, and insults remain strictly blocked.
- [x] **Phase 3 — Clean Up Scraped Dictionary Artifacts**:
  - In `suffixes`, trim hyphenated glosses (`abode-residence` -> `abode`, `refuge-shelter` -> `refuge`, etc.).
  - In `fruits`, `vegetables`, and `trees`, trim biological/taxonomic suffixes (`-Bush`, `-Leaf`, `-Root`, `-Citrus`, `-Stalk`).
- [x] **Phase 4 — Deduplication**:
  - Remove duplicate entries within `arabic_food`, `animals`, `gems_minerals`, and `architectural_elements`.
  - Remove duplicate entries in `blocked` list.
- [x] **Phase 5 — Automated Unit Tests**: Add a test ensuring that for every category $C$ and every word $w \in C$, `isSafe(w) == true`.
- [x] **Phase 6 — Attribution Commit**: Commit with `(goals/prune-blocked-list-and-resolve-self-blocking/goal.md)` using a `fix:` subject, then run `sila goals`.
