# Plan: Harden D1 Population Script and Eliminate CI Blocking Prompts

## Execution Steps
- [x] **Phase 1 — Non-Interactive Execution**: Remove interactive `input()` calls in `scripts/populate_d1.py` or gate behind `--interactive` flag to prevent CI crashes with `EOFError`.
- [x] **Phase 2 — SQL Safety & Deduplication**: Use `INSERT OR IGNORE` in `scripts/populate_d1.py` so single duplicate collisions don't fail entire batches.
- [x] **Phase 3 — Safe Batch Sizing**: Constrain batch sizes to stay well within CLI argument limits (`ARG_MAX`) and Wrangler/D1 statement limits (e.g. 50-100 rows per batch instead of 1,000).
- [x] **Phase 4 — Architectural Guidance**: Document in README that Cloudflare Workers should generate strings on-the-fly to avoid D1 read-after-write concurrency races and billing overhead.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/harden-d1-scripts-and-architecture/goal.md)` using a `chore:` or `fix:` subject, then run `sila goals`.
