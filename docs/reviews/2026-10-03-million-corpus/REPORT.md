# Two-million-identifier audit — 2026-10-03

**Verdict: the implemented category policy passes; the sensitivity review requests changes.**

I generated **2,000,000 real CLI identifiers** after building the current source,
then audited every physical line with the production analyzer and a separately
implemented parser. Neither batch contained `ahmad-donkey-food`, blocked vocabulary,
unknown expressions, mixed groups, invalid suffixes, or a violation of the requested
5–8 token bounds. Both validators reported zero mechanical violations.

Reading actual output exposed a remaining gap inside the allowed Islamic group:
`amr-ibn-al-as-jinn`, `abdur-rahman-ibn-abu-bakr-makruh`, and
`prophet-said-ibn-zayd-illumination-qadr` were actually generated. The group contract
allows these associations. It does not adequately protect named people from
unintended descriptions or religious attribution. This run should **not** be used
as a claim that all generated phrases are culturally safe.

This follow-up changes audit evidence and reproducibility tools. The production
generator and dictionary remain at the reviewed source snapshot.

## Actual generation and complete audits

Source commit: `24138744a6281fcd99ef06fc6447b9660a0c83c6`.
Fresh binary: `/tmp/hrs-million-cli`, built with Go 1.27.1.
Generation started at `2026-10-03T09:51:10.263929+00:00`; both batches existed before
the first audit started. The final replay reproduced the initial run's hashes.
Replays of the same two seeded batches are not counted as additional distinct
corpora in this report.

| Metric | Default with suffix | Default without suffix |
| --- | ---: | ---: |
| Seed | 1003202601 | 1003202602 |
| Actual physical lines | 1,000,000 | 1,000,000 |
| Unique complete identifiers | 999,998 | 997,385 |
| Duplicate occurrences beyond the first | 2 | 2,615 |
| Distinct identifiers that repeated | 2 | 1,053 |
| Most occurrences of one identifier | 2 | 55 |
| Exact uniqueness percentage | 99.9998% | 99.7385% |
| Bytes including newlines | 47,996,694 | 42,998,812 |
| Character lengths, excluding newline | 22–101 | 18–99 |
| Generation elapsed seconds | 1.896 | 1.686 |
| Production audit elapsed seconds | 15.250 | 15.105 |
| Independent scan elapsed seconds | 12.416 | 11.416 |
| Mechanical violations, each parser | 0 | 0 |

Every identifier had 5–8 alphabetic tokens, excluding the optional four-digit
suffix. Compound expressions consume multiple tokens. All five groups appeared
in both corpora. The independent parser found no complete identifier with more
than one group reading.

| Group | With suffix | Without suffix |
| --- | ---: | ---: |
| Asma Allah | 199,987 | 200,203 |
| Prophets | 199,833 | 199,822 |
| Personal names | 199,119 | 200,420 |
| Islamic | 200,661 | 199,911 |
| General | 200,400 | 199,644 |

The real corpora are retained locally:

- `/tmp/hrs-million-audit-2026-10-03/default-numbered.txt`
- `/tmp/hrs-million-audit-2026-10-03/default-unnumbered.txt`

Their SHA-256 hashes are respectively:

```text
0b202bd3dd89313cd044fda08de0705bc5b535a4cd3206d0e4954f1f867875c8
4cbd19a0a21bcea2e51e8e9d373421953ba19c47082dfb423f0ac4fa9bfc3465
```

The dictionary SHA-256 is
`5571a76622f3820f8f2d45007cd6872d7e133edd220043dfbfd7de68dfac024f`.
[runtime.json](runtime.json) preserves binary identity, exact arguments, exits,
hashes, full analyzer output, distributions, duplicate examples, and actual samples.
Large corpora are not committed; `/tmp` files are temporary and can be recreated.

## Sensitivity and quality findings

### 1. High: one broad Islamic group still permits sensitive person/concept associations

These are exact lines from the unnumbered million-line corpus:

| Line | Actual generated output | Review concern |
| ---: | --- | --- |
| 21,007 | `amr-ibn-al-as-jinn` | Can read as an inhuman description of a named companion |
| 17,604 | `abdur-rahman-ibn-abu-bakr-makruh` | Can read as a negative judgment attached to a named companion |
| 11,314 | `prophet-said-ibn-zayd-illumination-qadr` | Can read as an unintended religious title attached to a named companion |

A direct `Engine.ValidateString` probe accepted all three; see
[go-validation.json](go-validation.json). The production corpus analyzer and the
independent parser also accept them under the existing group contract. This is
a limitation in the contract's semantic precision.

After inspecting the initial samples, I scanned **both entire million-line files**
for complete dictionary person expressions from `sahaba` or
`islamic_golden_age_scholars` accompanied by selected standalone concept tokens:

| Cue beside a named person | With suffix | Without suffix | Total |
| --- | ---: | ---: | ---: |
| `jinn` | 267 | 275 | 542 |
| `creator` | 275 | 248 | 523 |
| `prophet` | 284 | 282 | 566 |
| `messenger` | 290 | 267 | 557 |
| `prophethood` | 270 | 285 | 555 |
| `makruh` | 277 | 273 | 550 |

There are **3,267 distinct candidate lines**: 1,649 numbered and 1,618 unnumbered.
Cue totals overlap. These are review candidates, not a count of established
offenses. For example, `creator` and `messenger` also have ordinary meanings.
The scanner excludes cue tokens inside atomic dictionary phrases such as
`birth-of-prophet-muhammad`. [semantic-candidates.json](semantic-candidates.json)
records the counts, exclusions, examples, person expressions, and corpus hashes.

**Recommended next change:** prevent named companions and historical scholars
from combining freely with religious concepts or descriptors. Emit reviewed
complete expressions for sensitive categories, or use an explicit list of reviewed
templates. General-only random generation is a smaller default surface while
religious compositions receive contextual review. Expanding a blocklist with
individual observed combinations will leave many equivalent associations possible.

### 2. Medium: possible glossary annotation in the personal-name dictionary

`muslim_names_female` contains both `Yumna` and `Yumna-Success`. The latter appears
to contain an English annotation rather than another reviewed name; this is a
dictionary-quality inference requiring confirmation.

Actual output at unnumbered line 704,891:
`abdul-hafiz-yumna-success-salwa-sumaira`.
The compound appeared in 3,667 numbered and 3,622 unnumbered lines, **7,289 total**.
Go accepts the complete compound as a name expression. A category label cannot
establish that every part of its entries is an appropriate personal name.

**Recommended next change:** review and remove annotation text from entries,
deduplicate the resulting names, and retain entry-level provenance.

### 3. Medium: a historical category contains a demonstrable classification error

`muslim_empires` includes `Nabataeans`; the output samples include it in the Islamic
group. The Met dates the Nabataean kingdom to antiquity and its Roman incorporation
to 106 AD, while its Islamic-history exhibition describes Islam's emergence in the
seventh century. Classifying that kingdom as a Muslim empire is inconsistent with
those dates. See [Nabataean Kingdom and Petra](https://www.metmuseum.org/essays/nabataean-kingdom-and-petra)
and [Byzantium and Islam](https://www.metmuseum.org/exhibitions/listings/2012/byzantium-and-islam).

**Recommended next change:** correct historical classifications with sourced
entry-level review. Do not use these random labels as historical assertions or
inventor attribution.

### 4. Medium: numeric suffixes reduce collisions but do not ensure uniqueness

The numbered corpus repeated `habib-ibn-zayd-al-ansari-1328` and
`umm-hakim-bint-al-harith-8629` twice each. Without suffixes,
`nasir-al-din-al-tusi` appeared 55 times and `abu-kamil-shuja-ibn-aslam` 52 times.
Several long atomic names and events already meet the token target, creating
frequently repeated complete labels. Repetition is allowed by the documented API;
applications needing unique identifiers must enforce it at storage time.

**Recommended next change:** use a unique constraint and bounded collision retries
where uniqueness is required. Measure distribution by complete expression before
changing sampling weights. The analyzer's two-decimal percentage prints `100.00%`
for 999,998/1,000,000; report exact counts or additional precision.

### 5. Low: the bulk path and uniqueness audit retain large in-memory collections

The CLI/library produces the result slice before printing, and the Python analyzer
keeps a `Counter` keyed by complete identifiers. Both retain data proportional to
batch size or unique output count. The actual million-line batches completed
successfully and generation took under two seconds each on this machine; this
follow-up did not measure peak memory.

**Optimization opportunity:** consider streaming/cancellation for service use and
disk-backed or external-sort uniqueness analysis for larger corpora. Validate every
line even if uniqueness accounting is separated. Semantic restrictions take
priority over further hot-path speed work.

## Auditor controls and verification

A separate copy of the numbered corpus retained exactly 1,000,000 physical lines
and replaced only these positions:

| Position | Injected value |
| ---: | --- |
| 1 | `ahmad-donkey-food-1234` |
| 500,000 | `aisha-cow-olive-2345` |
| 1,000,000 | `rahman-ramadan-3456` |

The production analyzer scanned the whole file, reported **exactly three**
violations with all three correct line numbers, and exited **1**. This control file
is deliberately corrupted and is excluded from the two-million generated-corpus
counts. Twelve targeted probes also failed the independent parser. Go directly
rejected the three analogous prohibited identifiers in
[go-validation.json](go-validation.json).

The independent scanner uses an explicit category table, a token trie, and group
hypothesis bitmasks. It does not import `NamingPolicy` or the production analyzer.
It checks every line, suffix, requested token bounds, blocked phrases, complete
group segmentation, probe hits, physical count, bytes, and SHA-256. It shares the
reviewed vocabulary, so both parsers can agree on a semantically weak entry or
composition; the manual findings above demonstrate this limitation.

I inspected 200 reservoir-selected real outputs, 20 group length extrema, and 96
targeted candidate examples. Manual inspection covered samples; automated passes
covered all two million lines. The complete replay generated and audited the same
two seeded corpora successfully, then reproduced the semantic counts.

The fresh build, `go test -race ./...` (cached result), `go vet ./...`, six Python
unit tests, and `git diff --check` passed. Exact receipts are in
[checks.json](checks.json). The broader original review remains available in
[the initial report](../2026-10-03-generator-audit/REPORT.md), with the completed
policy changes in [the implementation report](../2026-10-03-mixing-policy/REPORT.md).

## Reproduce

From the repository root, using the installed toolchain available in this session:

```sh
PATH=/home/jad/.local/share/mise/installs/go/1.27.1/bin:$PATH GOCACHE=/tmp/hrs-review-go-cache go build -o /tmp/hrs-million-cli ./cmd/halal-random-strings
python3 docs/reviews/2026-10-03-million-corpus/replay.py --cli /tmp/hrs-million-cli --corpus-dir /tmp/hrs-million-replay --receipts /tmp/hrs-million-replay/runtime.json --semantic-receipts /tmp/hrs-million-replay/semantic-candidates.json
```

The replay leaves both real corpora, generation and analyzer manifests, logs,
independent receipts, sensitivity candidates, and the corrupted control file in
the requested output directory. Its mechanical success does not certify the
semantic candidates. Reproduction should use the same product source/dictionary
snapshot when comparing these exact seeded hashes.
