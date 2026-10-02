# Facts: Enhance CLI UX, Color Accessibility, and Align Documentation

## Architectural Invariants & Constraints
- UNIX CLI ergonomics: `--help` is informational and must go to `stdout` with exit code 0; usage errors must go to `stderr` with non-zero exit codes.
- Accessibility standards: Terminal applications must support `NO_COLOR` standard (https://no-color.org) and avoid unreadable yellow-on-white text.
- Truth in documentation: Examples shown in `--help` and `README.md` must accurately reflect actual runtime behavior.

## File & Interface Contracts
- `cmd/halal-random-strings/main.go`: `usage()`, flag parsing, color styling.
- `README.md`: Usage section and example outputs.
- `Makefile`: `generate` target.
