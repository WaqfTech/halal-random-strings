# Generator code and sensitivity review — 2026-10-03

**Verdict: request changes. The current code does not enforce the requested “never generate insulting or sensitive combinations” requirement across its supported paths.**

Reviewed baseline: `124f1c3` in `/mnt/Jad/github/projects/waqftech/halal-random-strings`. This delivery adds review evidence and reproduction tools; the generator and dictionary have not been fixed.

The exact `ahmad-donkey-food` string was produced by an actual `NewEngine` call with a custom dictionary. It was **not** observed from the shipped default dictionary. Default generation and explicit name/animal category requests have useful protections, but those protections are bypassable through separators, overrides, and semantic misclassification.

## Actual testing and scope

I ran the CLI in a terminal, read real samples, generated two deterministic default corpora of 100,000 strings each, generated two targeted corpora of 10,000 strings each, exercised public library APIs, and ran an independent adversarial race probe. Automated corpus scans helped locate candidates; reading samples and examining the word lists supplied the semantic review. This is a product-policy review, not a claim that every flagged combination is universally offensive or a religious ruling.

| Check | Actual result |
| --- | --- |
| `go test -race ./...` | Passed; root package 1.475 s; CLI has no package tests |
| `go build ./...` and CLI build | Passed |
| `go vet ./...` | Passed |
| `govulncheck ./...` | Passed after network escalation: 0 reachable known vulnerabilities; 1 finding in a required module was not reachable |
| Default corpus, seed 42, no number | 100,000 lines; 98,283 distinct strings; 1,717 duplicate occurrences beyond the first |
| Default corpus, seed 2026, no number | 100,000 lines; 98,277 distinct strings; 1,723 duplicate occurrences beyond the first |
| Expanded name/sanctuary versus animal/food matching | 0 matches in the two default samples; 102 matches in the 10,000 `animals,architectural_elements` outputs |
| Expanded religious-place/caliphate marker review | 43 + 37 candidates in the default samples; 92 in `animals,nouns_places`; 200 in `animals,architectural_elements` |
| Existing corpus analyzer, seed 42 | Printed `[FAIL]` for a glossary artifact but exited 0 |
| Adversarial race probe | Reported a race at `halal-random-strings.go:501`; program exited 66 (`go run` exited 1) |
| Python 3.11 importer invocation | Failed to parse `scripts/populate_d1.py` at line 56 |
| `golangci-lint run ./...` | Incomplete: v1.64.8 reports unresolved `halalrandomstrings` references despite successful compiler/vet runs; not treated as a confirmed source compilation defect |
| `gofmt -l` | Reports the generator, its main test file, and CLI source as unformatted |

Corpus candidates use full phrase boundaries and a separate marker list; ordinary substrings inside longer phrases are not automatically insults. For example, “clock” inside `adhan-clock` is not evidence of disrespect. The expanded-marker counts are review candidates, not counts of adjudicated religious offenses. Sampling cannot prove all possible outputs safe.

Environment: Go 1.27.1, Python 3.13.5, Linux amd64, Intel i7-9700. The successful dependency scan needed network access; no production database was modified. Go 1.21 compatibility and remote D1 execution were not established by these checks.

## Findings, in priority order

### 1. P1 — Neutral categories leak religious terms into animal and food strings

Locations: `halal-random-strings.go:67`, `:123`, `:175`, `:431`; `words.json:730`, `:735`, `:778`; `architectural_elements` around `words.json:2700`.

The safety partition uses category names, rather than the semantic identity of the selected words. `nouns_places`, `architectural_elements`, and `suffixes` can be treated as neutral and appended to either domain. `nouns_places` still contains `shrine-holy`, `zawiya-shrine`, `khanqah-monastery`, `mausoleum-tomb`, and `caliphate`. Architecture repeats `Mihrab` and `Minbar`, which also exist in the protected sanctuary list. The current tests check sacred/mundane category overlap, but miss protected words inside neutral categories.

Actual default outputs:

```text
civet-jungle-shrine-holy-silver
shrine-holy-dwelling-green-peppercorn-falafel
caliphate-clean-butterfly-wadi-valley
ant-lake-zawiya-shrine-source
```

Actual category override outputs:

```text
theater-squirrel-shrine-holy-ox
sheep-canal-dome-koala-mihrab
minbar-bison-buffalo-hare-nave
```

These break the intended separation or create sensitive associations. Their words are individually allowed, so the blocked-word filter cannot repair them. Audit every word's semantic tags, including neutral categories and compound terms; preserve token provenance through assembly and validate the completed combination.

### 2. P1 — Category fallback bypasses the divine-name construction policy

Locations: `halal-random-strings.go:420`; `words.json` rules for `servant_prefixes,asma_allah`; `goal_delivery_test.go:459`.

The approved rule uses a servant prefix, but `--categories asma_allah` dynamically creates an unrestricted single-category rule. Actual output:

```text
ghaffar-fattah
sabur-razzaq
ghaffar-mumit
adl-mutakabbir
```

Command:

```sh
/tmp/hrs-review-cli --categories asma_allah --min-words 2 --max-words 2 --seed 42 -r 20 --no-random-number
```

The protection needs to be an invariant, rather than one preferred rule that users can bypass. Reject incomplete or unreviewed sacred-name category combinations. The single `servant_prefixes` category also allows meaningless repeated fragments instead of reviewed complete names.

### 3. P1 — Arbitrary separators evade blocking and corrupt identifier output

Locations: `halal-random-strings.go:136`, `:261`, `:336`, `:511`, `:550`; `cmd/halal-random-strings/main.go:171`.

Blocked matching recognizes punctuation as boundaries, but alphabetic separators concatenate tokens. The actual shipped CLI accepts:

```sh
/tmp/hrs-review-cli --categories muslim_names_male --sep donkey --min-words 2 --max-words 2 --seed 42 -r 10 --no-random-number
```

Actual outputs include `isadonkeymazen`, `nasserdonkeyamin`, and `huddonkeymaher`. A custom engine with `wine` explicitly blocked emits `winexdate` using separator `x`. A newline separator makes one requested identifier print as two lines.

Validate tokens before formatting, restrict separators to a documented harmless set, reject control characters and embedded offensive terms, and validate the final serialized identifier. The number suffix is appended after the current safety check, so validation should cover the complete result.

### 4. P1 — Custom dictionaries have no structural or mandatory safety validation

Locations: `halal-random-strings.go:123`, `:198`, `:230`, `:355`, `:503`.

`NewEngine` accepts a mixed-domain rule without checking it; `getRuleDomain` returns the first recognized domain, so a rule containing a name, animal, and food is classified as only one domain. Category-option validation is skipped when the caller uses the supplied rules without `Options.Categories`.

Actual API results:

```text
custom-engine-exact-prohibited-example              -> ahmad-donkey-food
custom-engine-mixed-domain-rule-with-blocklist      -> ahmad-cow-food
```

The first fixture has no custom blocked list; no mandatory core policy is inherited. The second fixture blocks `donkey` but still violates name/animal separation. These demonstrate public API behavior with custom data, not default-dictionary generation. Even trusted custom datasets can make accidental rule mistakes.

Validate every rule and normalized word at construction/load time. Mixed domains, unknown references, incomplete sacred forms, and unknown semantic classifications should fail closed. Keep a mandatory core policy that dictionary overrides cannot silently erase. `LoadWords("{}")` currently succeeds and replaces a healthy engine with unusable data; validate before swapping snapshots.

### 5. P1 — The sensitivity audit gives false passes and does not gate publication/import

Locations: `scripts/analyze.py:20`, `:53`, `:115`, `:132`, `:139`, `:178`, `:196`; `Makefile:60`; `scripts/populate_d1.py:68`.

Actual analyzer fixtures:

| Input | Existing analyzer result |
| --- | --- |
| `ahmad-donkey-food` | `[PASS]`, exit 0: removed/blocked dictionary words such as `donkey` are outside its semantic vocabulary |
| `ahmad-cow-rice` | `[FAIL]`, exit 0 |
| `ahmad_cow_rice` | `[PASS]`, exit 0: only hyphen-separated output is parsed |
| `abu-ubaidah-ibn-al-jarrah-cow` | `[PASS]`, exit 0: the entity has five tokens, while matching stops at four |
| `sad-ibn-abi-waqqas-cow` | `[PASS]`, exit 0: apostrophe handling differs from generator normalization |

The default seed-42 corpus also triggered a glossary-artifact failure for `umm-kulthum-bint-muhammad-fajr-dawn`, while returning exit 0. This occurrence can arise across word boundaries; whether the phrase is actually undesirable needs review. It still proves that a printed failure cannot fail `make verify-corpus`.

Other weaknesses: missing/unreadable dictionaries silently disable semantic checking; two broad domain sets are loaded but unused; only the first ten violations are retained and then presented as the total; blocked terms are not checked; the importer has no sensitivity gate. An offline importer probe successfully produced an INSERT containing the exact `ahmad-donkey-food` fixture. No database execution occurred.

Use one canonical normalization/policy definition, independent prohibited-combination fixtures, variable phrase lengths, explicit separator parsing, complete counts with capped examples, nonzero failure status, and a mandatory check before exporting/importing identifiers. Retain an independent audit oracle so the generator and checker cannot share the same blind spot.

### 6. P2 — Engine construction retains caller-owned mutable data; thread safety is incomplete

Locations: `halal-random-strings.go:175`, `:198`, `:501`; `halal-random-strings_test.go:165`.

The constructor keeps maps/slices and nested rule patterns from the caller. Mutating a pattern after construction changed a previously name-only rule into one producing `ahmad-cow-food`. A concurrent mutation probe reported a real data race between caller writes and generation at line 501. The internal `RWMutex` cannot protect writes made through those external aliases.

The existing concurrency test uses only concurrent readers of an unmodified engine, which explains its green result. Deep-copy inputs or compile immutable rule/word snapshots. Test constructor ownership and generation concurrent with supported loading operations.

### 7. P2 — Impossible requests perform up to a million rule attempts per output and fail late

Locations: `halal-random-strings.go:324`, `:474`, `:485`, `:496`, `:581`; CLI limits at `cmd/halal-random-strings/main.go:158`.

The nested 1,000-by-1,000 retry loops repeat impossible work. Failure is checked only after processing the entire batch, even when index 0 already failed. Requesting exactly one word from the two-token servant-prefix/divine-name rule took approximately 0.138 s for one result, 1.33 s for ten, and 25.8 s for 100 under this session's load. These timings are observations, not a controlled throughput comparison.

The CLI allows one million results; the public API has no repeat/output-size cap, and the read lock is held for the whole batch. If exposed through a service, unbounded options can become a resource-exhaustion path. Inverted bounds are silently rewritten: `--min-words 9 --max-words 2` succeeds with nine words.

Check feasible token lengths before generating, reject invalid bounds explicitly, return immediately when an output cannot be made, impose allocation/work budgets, and add cancellation for service/batch use. Avoid retaining the engine lock throughout an entire batch by using an immutable snapshot.

### 8. P2 — The documented rule contract does not match execution

Locations: `halal-random-strings.go:35`, `:389`, `:416`, `:501`, `:503`; `README.md` dictionary description.

`Rule.Template` is ignored. A template `constant-{custom}-suffix` yields only `safe`. Missing categories inside a rule are silently skipped: a rule referencing `nonexistent,custom` also emits `safe`.

Multi-category requests can silently omit a requested category. `--categories colors_arabic,muslim_names_female` succeeds with female names only because a single-category name rule matches. Another fallback, `colors_arabic,holy_sanctuaries`, generates arbitrary combinations such as `turquoize-kaaba` and `mecca-mushaf`. The options are neither a dependable required-category contract nor a reviewed semantic composition contract.

Define whether category selection means allowed categories or required categories; document and test that contract. Implement templates with validation or remove the unused field and its advertised behavior. Reject malformed rules instead of degrading them into different patterns. The compatibility `GenerateWithOptions` method also discards errors (`halal-random-strings.go:592`); several tests iterate its results without requiring the requested result count, so an empty slice can make those checks pass vacuously. Prefer the error-returning API and explicit result-count assertions.

### 9. P2 — Blocking has normalization and phrase-length holes

Locations: `halal-random-strings.go:136`, `:155`, `:261`, `:280`, `:303`.

Dictionary normalization strips apostrophes; blocked-list normalization turns them into boundaries. A custom dictionary with `s'mores` blocked generated `smores`. This term is already present in the shipped blocked list, although absent from its generation categories. Five-token blocked phrases are never matched; a fixture blocking `alpha-beta-gamma-delta-epsilon` still generates that exact phrase.

Unicode letters and non-ASCII apostrophes have inconsistent treatment between generator, filter, and Python auditor. Define the supported alphabet and transliteration rules, normalize all policy inputs consistently, and use the maximum configured blocked-phrase length rather than a hard-coded four-token window. Preserve word-boundary matching to avoid blocking innocent substrings such as those inside personal names.

### 10. P2 — “Unique” and “cryptographically secure” claims exceed the implementation

Locations: `halal-random-strings.go:56`, `:476`, `:553`, `:558`; `README.md` features; CLI help at `cmd/halal-random-strings/main.go:68`.

Generation has no uniqueness check. The two no-number default batches each had about 1.7% duplicate excess, and low-cardinality single categories can collide much more often. A 4-digit suffix has only 9,000 possibilities, about 13.14 bits, regardless of whether its source is cryptographically secure.

Words use time-seeded `math/rand`; setting a seed makes the suffix deterministic too. On supported older toolchains, the crypto-error branch silently falls back to that PRNG. These properties are useful for reproducible labels but unsuitable as a promise of unguessable access credentials. The Go documentation explicitly excludes `math/rand` from security-sensitive use: [official package documentation](https://pkg.go.dev/math/rand).

For labels, describe probabilistic collision resistance honestly and retry collisions at an authoritative unique store. For invitation authorization, use a separate high-entropy secret and rate limits; do not treat the readable label as the sole bearer credential. No live authorization service exists here to test.

### 11. P2 — D1 pipeline tests bless SQL that conflicts with D1 import guidance

Locations: `scripts/populate_d1.py:55`, `:65`, `:79`, `:84`, `:193`; `goal_delivery_test.go:800`.

An actual offline dump emits `BEGIN TRANSACTION;` and `COMMIT;`; the tests require those strings. Current Cloudflare import guidance says to remove those statements before `wrangler d1 execute --file`, and identifies nested-transaction failures: [official D1 import documentation](https://developers.cloudflare.com/d1/best-practices/import-export-data/). This is a verified generated-SQL/documentation conflict; a production D1 failure was not reproduced.

The Python importer also fails to parse on Python 3.11 because a backslash appears inside an f-string expression. It works in this session's Python 3.13, so the Go tests do not reveal the portability defect. Precompute the joined SQL text or explicitly require Python 3.12+. Test D1-compatible imports through the intended Wrangler path, not just SQL substrings or ordinary SQLite.

Additional weak points: database/auth/table checks parse human-readable text and use substring matches; the schema path depends on the current directory; batch size is unchecked; duplicate inserts are ignored without reporting inserted versus skipped totals; local versus remote execution is not explicit in the script.

### 12. P3 — Corpus policy, documentation, and verification need clearer ownership

Locations: `words.json` blocked list and `nouns_objects`; `README.md`; `docs/TODOS.md`; `cmd/halal-random-strings/main.go:91`; `.github/workflows/dependabot-sync.yml`.

The blocked list mixes profanity and prohibited products with ordinary words such as `teenager`, `depression`, `exam`, `computer`, and `cute`, without reasons or context. Meanwhile allowed objects include `dagger`, `spear`, and `hoe`, although `knife`, `sword`, and `weapon` are blocked. These are policy inconsistencies and potential sensitivity issues, not proof that every listed word is offensive. Opaque categories such as “impure animals” should not substitute for a documented culturally reviewed naming policy.

Sacred phrases and religious events remain arbitrary identifier material. Ordinary suffixes include `creator`, `master`, and `tool`; inventing name/epithet or religious-term/product associations can cause avoidable confusion. Random scholar/invention pairings can look like historical attributions even though there is no reviewed attribution mapping. Preserve complete reviewed names and phrases rather than splicing fragments.

README promises D1/API integration but points to an absent `api_worker.js`/`.ts`; it should distinguish the working generator/import tooling from a proposed API. Several CLI samples violate the advertised word counts or show old category combinations. The only tracked GitHub workflow handles dependency syncing; there is no test/vet/sensitivity CI gate. Its reusable upstream workflow uses mutable `@main` with a PAT secret; pinning and least-privilege access deserve a separate review before reliance on it.

## Measured optimization opportunities

Six benchmark repetitions measured:

| Operation | Median | Allocation |
| --- | --- | --- |
| `Generate()` | 1.973 μs/string | approximately 605 B, 25 allocations |
| `GenerateWithOptions()` | 1.929 μs/string | approximately 605 B, 25 allocations |

The normal path is already fast. A short allocation profile attributed about 67% of sampled allocated bytes cumulatively to `isSafeLocked`, including repeated phrase builders. Improve that path only after correcting its semantics: compile blocked phrases into a trie or match token spans, reuse canonical prevalidated token metadata, and size builders for the actual output. Keep a final semantic combination check; prevalidating individual words alone does not solve these findings.

Other useful opportunities are feasibility-aware rule selection, immediate failure instead of full-batch retries, immutable snapshots instead of long-held read locks, and an optional streaming batch API. The Python analyzer reads incrementally but still retains every distinct line, every length, and sorted duplicate collections; calculate length metrics online and bound examples. Exact deduplication still needs storage proportional to unique inputs unless deliberately replaced with approximate statistics.

No implementation was changed, so no speedup is claimed. Profile text and benchmark readings are included in the evidence folder; large generated corpora and profile binaries remain temporary.

## Ordered remediation and acceptance checks

1. **Define and enforce one mandatory naming policy.** Assign semantic tags to each normalized word/phrase, including formerly neutral categories. Preserve full personal/sacred phrases and allow only reviewed compositions. Consider an everyday identifier mode with a small neutral word list and a separate mode that emits a complete reviewed respectful phrase. Acceptance: every factory/load/category/separator path rejects the user's exact example and all named-animal-food variants; all protected phrases stay in their approved form.
2. **Validate and own dataset snapshots.** Deep-copy inputs, validate rules/templates/references, reject mixed and unknown semantic domains, and swap only validated snapshots. Acceptance: malformed load leaves the previous engine intact; caller mutations cannot affect results; independent race probes are clean.
3. **Make audit/export/import failures enforceable.** Share canonical normalization but maintain independent prohibited fixtures and exhaustive legal-rule checks. Acceptance: the five analyzer fixtures above fail with nonzero status; alternative separators, long names, blocked-only words, missing dictionary data, and more than ten violations are handled correctly; unsafe rows never reach SQL export.
4. **Bound work and honor the option contract.** Reject impossible/inverted lengths before allocation, cap output/work, and fail at the first unsuccessful item. Acceptance: impossible large batches fail promptly, cancellation releases resources, and category/template behavior is explicitly tested.
5. **Correct identity/security and D1 documentation.** Separate readable labels from authorization secrets and uniqueness enforcement; produce D1-compatible SQL and verify supported Python versions. Acceptance: collision retry behavior and actual intended Wrangler import path have receipts; docs do not advertise absent APIs.
6. **Add CI and measure only remaining bottlenecks.** Gate compiler, vet, race/ownership tests, dataset validation, and independent sensitivity checks. Acceptance: deliberately broken safety fixtures fail the gate; benchmark/profiling changes preserve policy behavior and demonstrate measured benefit on the same toolchain/hardware.

## Evidence and replay

See [evidence](evidence/) for CLI/API/analyzer receipts, corpus hashes and counts, race output, SQL dump, and benchmark/profile text. [replay.py](replay.py) rebuilds the CLI and replays bounded local probes; it does not access D1 or modify production data. The Go probe sources use a `.go.txt` suffix to keep adversarial fixtures out of the normal package build.

```sh
python3 docs/reviews/2026-10-03-generator-audit/replay.py --go /home/jad/.local/share/mise/installs/go/1.27.1/bin/go
```

Sila reported all 21 existing goals implemented. That tracker state does not close the defects reproduced here. `sila heal` initially failed on a read-only shared directory; healing, the review decision, and the friction report subsequently succeeded through approved escalation. Healing added managed `GEMINI.md` wiring, committed separately from this audit. Required tests and evidence are recorded above independently of tracker state.
