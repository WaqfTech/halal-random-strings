# Goal: Enhance CLI UX, Color Accessibility, and Align Documentation

## Goal Description
Fix usability bugs and eliminate documentation inconsistencies:
1. **Accurate CLI Examples**: Update `usage()` sample outputs in `cmd/halal-random-strings/main.go` and `README.md` to show actual 4-digit numbers (`1000-9999`) rather than outdated 9-digit numbers (`123456789`).
2. **Standard Help Stream**: Ensure `--help` outputs to `os.Stdout` instead of `os.Stderr` so users can pipe help text (`| less`).
3. **Color Accessibility & `NO_COLOR`**: Implement adaptive terminal styling or respect `NO_COLOR`/`TERM=dumb` so that commands and examples are readable on both dark and light terminal backgrounds.
4. **Align Generation Limits**: Reconcile `README.md` claim of generating 100,000 strings with the CLI `maxResults` cap (1,000) and `Makefile` target.
5. **Handle Dead Fields**: Address unused struct fields `PrefixThreshold` and `SuffixThreshold`, and unparsed `Rule.Template` (either implement them or cleanly document/deprecate them).

## Dependencies & Execution Order
- **Mode**: Independent ⚡
- **Depends On**: -
- **Sequence**: 19
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/enhance-cli-ux-and-resolve-documentation-drift/facts.md`](facts.md)
- **Execution Plan**: [`goals/enhance-cli-ux-and-resolve-documentation-drift/plan.md`](plan.md)

## Done Condition
1. CLI help examples display realistic 4-digit numeric suffixes matching actual output.
2. `halal-random-strings --help` prints to stdout and exits with code 0.
3. Terminal styling complies with `NO_COLOR` and provides readable contrast.
4. CLI documentation in `README.md`, `Makefile`, and `main.go` are completely aligned.
5. All tests pass with `go test -race ./...`.
