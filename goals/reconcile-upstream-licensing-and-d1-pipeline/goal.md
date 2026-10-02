# Goal: Reconcile Upstream Licensing Attribution and Harden D1 Pipeline

## Goal Description
Resolve open-source licensing compliance and modernize the database pipeline:
1. **Restore Charmbracelet Upstream MIT Attribution**: Since `halal-random-strings` is a fork of `charmbracelet/hotdiva2000`, restore the original Charmbracelet MIT copyright and license notice in `LICENSE` or a dedicated `NOTICE` file as required by the MIT license.
2. **Clarify Repository Branding**: Remove the misleading "Part of Charm" banner and Charm logo from `README.md` to avoid falsely implying this is an official Charmbracelet release, while preserving full, respectful credit to the upstream authors.
3. **Batch D1 Inserts into SQL File**: `scripts/populate_d1.py` currently executes `wrangler d1 execute --command` in a loop of 50 items. For 100k items, this spawns 2,000 separate CLI processes. Refactor to generate a consolidated `batch.sql` file and execute via `wrangler d1 execute --file` (or D1 HTTP bulk API), reducing execution time from hours to seconds.
4. **Harden SQL Formatting**: Replace string concatenation SQL construction with robust escaping to prevent query syntax breakages.
5. **Streaming in `scripts/analyze.py`**: Refactor `analyze.py` to stream lines from disk instead of reading the entire file into memory via `read().splitlines()`.

## Dependencies & Execution Order
- **Mode**: Independent ⚡
- **Depends On**: -
- **Sequence**: 20
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/reconcile-upstream-licensing-and-d1-pipeline/facts.md`](facts.md)
- **Execution Plan**: [`goals/reconcile-upstream-licensing-and-d1-pipeline/plan.md`](plan.md)

## Done Condition
1. Upstream Charmbracelet MIT copyright notice is legally compliant and preserved.
2. README branding clearly identifies WaqfTech as the fork maintainer without claiming to be an official Charm project.
3. `populate_d1.py` inserts batches via SQL file execution or transactions without spawning thousands of Wrangler subprocesses.
4. `analyze.py` operates in streaming mode and displays duplicate counts clearly.
5. All tests pass with `go test -race ./...`.
