# Plan: Harmonize License Inconsistency, Update Security Policy and Docs

## Execution Steps
- [x] **Phase 1 — License Reconciliation**: Resolve the contradiction between MIT and Waqf GPL 2.0 copyleft terms in README and docs, clarifying that code is under Waqf / MIT as intended. Add WaqfTech to the copyright notice.
- [x] **Phase 2 — Finalize SECURITY.md**: Replace the unedited boilerplate `(Replace with actual email address)` with the designated contact channel.
- [x] **Phase 3 — Reconcile Documentation Drift**: Update `README.md` to remove broken references to non-existent scripts (`scripts/test_uniqueness.sh`). Update `docs/TODOS.md` to remove already-implemented items (`min_words`, `max_words`, `sep`).
- [x] **Phase 4 — Attribution Commit**: Commit with `(goals/harmonize-licensing-and-security-policy/goal.md)` using a `docs:` or `chore:` subject, then run `sila goals`.
