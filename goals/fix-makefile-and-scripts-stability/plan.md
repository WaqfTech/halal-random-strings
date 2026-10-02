# Plan: Fix Makefile Typo, Process Fork Bomb, and Script Zero-Division

## Execution Steps
- [x] **Phase 1 — Fix Makefile Typo**: Fix `%(GO_APP_NAME)` typo to `$(GO_APP_NAME)` on line 36 of `Makefile`.
- [x] **Phase 2 — Replace Process Fork Bomb**: In `Makefile` line 49, replace the 1,000-process bash loop with a single invocation of `./$(GO_APP_NAME) -r 1000 > $(OUTPUT_FILE)`.
- [x] **Phase 3 — Fix Default Target**: Change `all` in `Makefile` so it only builds, tests, and analyzes locally without attempting cloud deployment (`populate-d1`).
- [x] **Phase 4 — Fix ZeroDivisionError in analyze.py**: Add a check in `scripts/analyze.py` so that empty input files output a clear message rather than raising `ZeroDivisionError`.
- [x] **Phase 5 — Attribution Commit**: Commit with `(goals/fix-makefile-and-scripts-stability/goal.md)` using a `fix:` subject, then run `sila goals`.
