# Plan: Fix Custom Separator Formatting and Word Count Boundaries

## Execution Steps
- [x] **Phase 1 — Dynamic Separator Injection**: Eliminate hardcoded hyphens in `words.json` templates or ensure separator replacement transforms all compound delimiters into `opts.Sep`.
- [x] **Phase 2 — Unified Word Counting**: Count words *after* whitespace and internal punctuation normalization so that multi-word phrases (e.g. Sahaba names) are counted accurately against `MinWords` and `MaxWords`.
- [x] **Phase 3 — Fixed-Width Numeric Suffix**: Replace unpredictable 1-to-18-digit suffixes with a consistent, configurable fixed-width integer (e.g. 4 to 6 digits) that does not overflow 32-bit systems.
- [x] **Phase 4 — Verification**: Verify that `--sep _` produces pure underscore delimiters and outputs strictly obey `--min-words` and `--max-words`.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/fix-separator-and-word-count-logic/goal.md)` using a `fix:` subject, then run `sila goals`.
