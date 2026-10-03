# Mixing-policy implementation and verification — 2026-10-03

The approved restrictions are implemented. Custom dictionaries are disabled;
Asma Allah and prophets have separate self-only groups; personal names have a
names-only group; protected Islamic and general categories are segregated. There
are no neutral categories or cross-group fallbacks. The default generator chooses
one approved group per identifier.

The engine owns immutable embedded data, validates it at startup, proves requested
lengths feasible before sampling, and requires every selected category to supply
a complete expression. Arbitrary separators and invalid bounds are rejected.
Blocked phrases use a token trie, consistent apostrophe normalization and no
four-token cutoff. The numeric suffix is included in the final blocked check.

See [the full category policy and API migration](../../MIXING_POLICY.md). This
implements the user's approved policy; it is not an independent religious ruling.

## Actual testing

| Check | Observed result |
| --- | --- |
| Full Go race suite | Passed, including concurrent default reads and mutation of rejected caller-owned dictionaries |
| `go vet ./...`, `go build ./...` | Passed |
| Go category matrix | All 1,296 ordered pairs checked: 536 permitted and 760 rejected |
| Python policy/audit/export tests | Six tests passed on Python 3.13 and 3.11 |
| CLI replay | 29 positive/negative/invalid-separator/impossible-request cases passed |
| Large generated corpora | 340,000 identifiers; zero policy violations |
| Analyzer negative fixtures | All 11 fixtures returned exit 1 |
| Unsafe offline SQL export | Exit 1; existing dump preserved |
| `make verify-corpus` | Additional 10,000 generated identifiers passed |

The 340,000 replay identifiers comprise 200,000 default numbered identifiers,
100,000 default unnumbered identifiers, 10,000 animal/architecture identifiers,
and 10,000 default numbered identifiers for each of `_`, `.`, and `/`.
The numbered default batch had 200,000 distinct lines. The unnumbered batch had
99,915 distinct lines: the generator does not guarantee uniqueness.

The exact `ahmad-donkey-food` example fails. Named-person/prophet combinations
with ordinary animals and food also fail without depending on donkey being
blocked. Previously observed leaks (`mihrab-cow`,
`shrine-holy-dwelling-green-peppercorn-falafel`), long Sahaba names plus animals,
divine names plus religious terms, and personal names plus religious terms fail.
`donkey`, `wine`, letters, whitespace, newline and double-hyphen separators fail
before generating any output.

A one-million-result request for one-token Sahaba names failed immediately with
exit 1 and no output in approximately 6.3 ms in the recorded CLI run. Feasibility
is checked before result allocation; the prior nested million-attempt retry path
has been removed.

Actual permitted CLI strings include `iyad-shaista`, `beaver-mashrabiya`,
`fox-alley`, `cricket-frieze`, and `hyrax-muqarnas`. Whole names and expressions stay
intact; `abdul-hadi-maher` and `sad-ibn-abi-waqqas-ramadan` pass. A botanical name
such as `bird-of-paradise` is matched as a complete general expression rather
than interpreted as the religious word `paradise`.

Human inspection of CLI samples exposed `kaffir-lime-leaves`. It was replaced by
`makrut-lime-leaves`, with offensive spellings blocked. [American Heritage](https://ahdictionary.com/word/search.html?q=kaffir+lime)
labels the former name offensive; [UC Riverside](https://citrusvariety.ucr.edu/crc2454)
records makrut as an alternative. The final corpora were generated after this fix.

## Data and compatibility

The curated dataset contains 36 categories and 3,526 entries. Standalone prefixes
were removed, prophet spellings and associated titles were isolated, religious
architecture/places/suffixes were moved out of general lists, and ambiguous whole
expressions were assigned a single group. The curation log records each change.
Historical scholars of different faiths remain in the protected historical group;
this does not label their personal faith.

The custom-constructor signature remains available, but the returned engine is
disabled and reports `ErrCustomDictionariesDisabled`. `LoadWords` fails without
reading or changing anything. Deprecated generation wrappers still return empty
results on errors; the error-returning API is documented and tested.

The prior tests requiring servant-prefix/divine combinations, arbitrary runtime
dictionaries, silent inverted-bound rewriting, successful empty audits, and D1
transaction statements were superseded by the current approved policy and error
contract. They were changed to meaningful rejection or compatibility checks.
The original audit and its evidence are unchanged.

The Python analyzer no longer silently skips missing policy data, omits the
blocklist, hides all but ten violations in the total, or exits successfully after
finding violations. The SQL exporter validates before Wrangler setup and during
writing, publishes only a complete validated dump, supports Python 3.11, checks
batch bounds, and omits D1-incompatible transaction statements. The latter follows
[Cloudflare's import guidance](https://developers.cloudflare.com/d1/best-practices/import-export-data/).
SQLite execution and export rejection were tested locally; no remote D1 writes
or deployment were performed.

## Performance and remaining scope

On the same Go 1.27.1 toolchain and i7-9700 host, three benchmark repetitions
measured median `Generate()` at 1.505 μs, 608 B and 13 allocations per identifier.
The earlier audit measured 1.973 μs, about 605 B and 25 allocations. The new
composition distribution differs, so these are observed benchmark results rather
than a universal throughput guarantee. The trie removes repeated blocked-phrase
string construction; default feasibility plans are precomputed.

The approved implementation addresses the safety bypasses in the
[original audit](../2026-10-03-generator-audit/REPORT.md), with the old divine-prefix
requirement replaced by the user's newer self-only instruction. Documentation now
states the actual category contract, disabled custom dictionaries, error behavior,
probabilistic uniqueness and absent serving API.

Remaining work from the broader audit includes rationale and expert review for
the blocklist and religious dictionary, especially transliteration aliases and
inconsistent object exclusions; CI security/test gates; explicit Wrangler
local/remote targeting and reliable structured setup checks; inserted-versus-
skipped reporting; and optional cancellation/streaming for large service batches.
Random scholar/invention pairs are labels, not verified historical attributions.
The dictionary and grouping do not certify every possible cultural interpretation.
No uniqueness, authorization-service or remote D1 behavior is claimed.

## Receipts and replay

- [runtime.json](runtime.json): CLI arguments, exits, examples, corpus hashes and auditor output.
- [checks.json](checks.json): full race/vet/build and Python check output.
- [benchmark.txt](benchmark.txt): all three benchmark repetitions.
- [make-verify.txt](make-verify.txt): repository corpus target result.
- [curation.json](curation.json): dictionary changes and rationale.

```sh
go build -o /tmp/hrs-policy-cli ./cmd/halal-random-strings
python3 docs/reviews/2026-10-03-mixing-policy/replay.py --cli /tmp/hrs-policy-cli --receipts /tmp/hrs-policy-replay.json
```

The replay creates temporary corpora and retains only hashes, samples and metrics.
It performs no network calls or Cloudflare writes.
