# Plan: Reconcile Upstream Licensing Attribution and Harden D1 Pipeline

## Execution Steps
- [ ] **Phase 1 — Upstream Notice Restoration**:
  - Add Charmbracelet's original MIT copyright notice to a `NOTICE` file or append to `LICENSE` under an Upstream Credits section.
  - Revise `README.md` to remove "Part of Charm" and replace with clear "Forked from charmbracelet/hotdiva2000" attribution.
- [ ] **Phase 2 — D1 SQL File Batching in `populate_d1.py`**:
  - Write inserts in blocks of 500-1000 statements wrapped in transactions (`BEGIN TRANSACTION ... COMMIT;`) to a temporary `.sql` file.
  - Execute via `wrangler d1 execute <db> --file=<temp.sql>` to minimize process spawning overhead.
  - Implement robust escaping for special characters.
- [ ] **Phase 3 — Stream Processing in `analyze.py`**:
  - Read `output.txt` line by line with a generator/iterator.
  - Output duplicate summary with frequency counts (`item (x3)`).
- [ ] **Phase 4 — Integration Verification**: Run `make generate` followed by `python3 scripts/analyze.py output.txt`.
- [ ] **Phase 5 — Attribution Commit**: Commit with `(goals/reconcile-upstream-licensing-and-d1-pipeline/goal.md)` using a `chore:` subject, then run `sila goals`.
