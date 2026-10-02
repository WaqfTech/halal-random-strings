# Plan: Enhance CLI UX, Color Accessibility, and Align Documentation

## Execution Steps
- [x] **Phase 1 — Update Sample Outputs & Formatting**:
  - Update all `printExample` calls in `cmd/halal-random-strings/main.go` to use realistic 4-digit numbers (e.g. `1234`, `9876`).
  - Update README sample strings to match.
- [x] **Phase 2 — Output Stream Fix**:
  - Modify `usage()` to accept an `io.Writer` target so `--help` directs to `os.Stdout` and invalid flag errors direct to `os.Stderr`.
- [x] **Phase 3 — Contrast & Color Configuration**:
  - Check for `NO_COLOR` environment variable or non-TTY outputs to disable lipgloss styling.
  - Choose colors with high contrast on both light and dark terminal palettes.
- [x] **Phase 4 — Resolve Documentation & Flag Limit Drift**:
  - Align `Makefile` and `README.md` batch count expectations.
  - Review `PrefixThreshold` and `SuffixThreshold` API surface.
- [x] **Phase 5 — Verification & Tests**: Test CLI invocation with `--help`, `--no-color`, and piped stdout.
- [x] **Phase 6 — Attribution Commit**: Commit with `(goals/enhance-cli-ux-and-resolve-documentation-drift/goal.md)` using a `fix:` subject, then run `sila goals`.
