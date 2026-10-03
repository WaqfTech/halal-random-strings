# Facts: Enforce reviewed category mixing policy

## User directive

The user approved personal names in a names-only group and disabling custom
dictionaries, and authorized implementation while AFK. Asma Allah and prophets
each mix only within their own category. Religious categories mix only with
religious categories; general categories require removal of protected entries.

## Architectural invariants

- Embedded, validated data only; no caller-owned dictionary references or runtime swaps.
- Every category belongs to an explicit group; no neutral/fallback bypass.
- Every selected category contributes a complete expression; no fragment splicing.
- Common names, titles and religious/general homonyms require conservative curation.
- Final outputs pass blocked vocabulary checks after the numeric suffix is appended.
- Audit and SQL export reject unsafe, empty, unknown or malformed input.
- Local implementation and commit only; no remote push or Cloudflare writes.

## Scope and contracts

Changed paths: `words.json`, `go.mod`, `go.sum`, `halal-random-strings.go`, `policy.go`,
`cmd/halal-random-strings/main.go`, Go tests, `scripts/`, `Makefile`, `.gitignore`,
`README.md`, `docs/`, and this goal package/index.

`NewEngine` and `LoadWords` now fail closed. `GenerateWithOptionsE` is the supported
error-reporting path; legacy wrappers return no output on invalid requests.
`IsSafe` retains blocked-vocabulary-only semantics; `ValidateString` enforces the
complete policy. See [MIXING_POLICY.md](../../docs/MIXING_POLICY.md).

Prior prefix-based divine generation, arbitrary separators, runtime overrides,
silent bound rewriting and the empty-audit success assertions are superseded by
the current user directive and explicit error contract. Their regression tests
were replaced with rejection checks; historical audit evidence was preserved.

## Verification

See [remediation report](../../docs/reviews/2026-10-03-mixing-policy/REPORT.md) and
its runtime JSON/benchmark receipts. Full race suite, vet, build, Python tests,
1,296 category pairs and 340,000 generated CLI identifiers passed. CLI failures
produce exit 1 and no output; unsafe SQL input leaves the existing dump intact.
No remote D1 or authorization service behavior is claimed.
