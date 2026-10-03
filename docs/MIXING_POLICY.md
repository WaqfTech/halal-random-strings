# Reviewed mixing policy

Each identifier belongs to exactly one permitted generation group. There is no
neutral category or fallback that bypasses group checks. With no category filter,
the engine chooses one of the five groups for each identifier. Selecting categories
requires all of them to be in the same group; every distinct selected category
contributes at least one complete dictionary expression.

| Group | Categories |
| --- | --- |
| Asma Allah only | `asma_allah` |
| Prophets only | `prophets` |
| Personal names only | `muslim_names_male`, `muslim_names_female` |
| Protected Islamic context | `sahaba`, `islamic_art_forms`, `islamic_events`, `islamic_golden_age_scholars`, `islamic_inventions`, `islamic_months`, `islamic_virtues`, `scholarly_terms`, `adab_terms`, `holy_sanctuaries`, `muslim_empires`, `nouns_concepts`, `adjectives` |
| General | `animals`, `arabic_food`, `architectural_elements`, `birds`, `colors_arabic`, `days_of_week_arabic`, `flowers`, `fruits`, `gems_minerals`, `geographic_features`, `jordanian_food`, `nouns_objects`, `nouns_places`, `saudi_food`, `spices`, `suffixes`, `trees`, `vegetables`, `yemeni_food` |

The existing `adjectives` contain Arabic names and religious descriptors, and
`nouns_concepts` contain religious concepts. They stay protected. Caliphates in
`muslim_empires` also stay protected. The historical scholars list includes
contributors of different faiths; its group describes the protected historical
context, not an assertion about each person's faith.

Known prophet name spellings were moved out of personal names. Familiar associated
titles/spellings (`Mustafa`, `Taha`, `Yasin`, `muzammil`, `mudathir`) are also isolated
conservatively. This is not a ruling on disputed names or titles. Sources for the
added name spellings include [Qur'an 6:83–86](https://quran.com/en/al-anam/83-86),
[19:56](https://quran.com/maryam/56), and [21:85](https://quran.com/al-anbya/85).
The list is not an exhaustive list of every transliteration.

Standalone `servant_prefixes` are disabled. Asma Allah cannot mix with prefixes or
other categories. Complete personal names such as `Abdul-Hadi` remain atomic
entries in the personal names group. Exact expressions shared by multiple groups
are retained only in the more protected group.

Mosque elements, religious schools, shrines and other religious expressions were
moved out of general architecture, places and suffixes. Religious architecture
from other traditions was excluded from general generation pending an appropriate
reviewed context. The complete changes are in
[curation.json](reviews/2026-10-03-mixing-policy/curation.json).

The spice entry `Kaffir Lime Leaves` was renamed `Makrut Lime Leaves`; offensive
spellings are blocked. The [American Heritage Dictionary](https://ahdictionary.com/word/search.html?q=kaffir+lime)
labels the former name offensive, and [UC Riverside's citrus collection](https://citrusvariety.ucr.edu/crc2454)
records makrut as an alternative name. Ordinary compound botanical names such as
`bird-of-paradise` remain complete general expressions.

## API changes

Custom dictionaries are disabled. `NewDefaultEngine()` is the supported constructor.
`NewEngine(Words)` retains its signature for source compatibility but returns a
disabled engine: `GenerateWithOptionsE` reports `ErrCustomDictionariesDisabled`,
compatibility generation methods return no output, and no caller data is retained.
Both `LoadWords` methods reject overrides without reading files or modifying state.

Use `GenerateWithOptionsE` to handle errors explicitly. The older wrappers return
`nil` or an empty string on errors. `IsSafe` checks **blocked vocabulary only**;
use `ValidateString(identifier, separator)` for format and mixing-policy checks.

Supported separators are `-`, `_`, `.`, and `/`. Letters, whitespace, newlines,
multiple separators and arbitrary text are rejected. Word bounds count separated
tokens, including all tokens in a compound name, and exclude the optional four
digit suffix. Zero values default to one result and 5–8 words. Other bounds must
satisfy `1 <= min <= max <= 64`; repeat must be 1–1,000,000. Impossible category and
length requests fail before sampling. Nonzero seeds reproduce results within this
version; this change intentionally changes prior seeded outputs.

The embedded `policy` metadata is shared by Go and Python. The embedded `rules`
are retained as validated composition examples for compatibility, not runtime
customization. The engine constructs feasible combinations directly from the
approved category pools. Additions require dictionary review, policy classification,
the explicit category matrix test, rebuilding the binary, and a corpus audit.

## Validation tools

```sh
go test -race ./...
go vet ./...
go build ./...
python3 -m unittest discover -s scripts -p 'test_*.py'
make verify-corpus
python3 scripts/analyze.py output.txt --sep -
python3 scripts/populate_d1.py output.txt --dump-sql /tmp/validated.sql --sep -
```

Python 3.9 or later is required. The analyzer checks every physical line, preserves
invalid whitespace, reports total violations, and exits 1 for unsafe or empty
corpora or missing/invalid policy data. It accepts an identifier only when its
entire text can be segmented into complete expressions from one group. Long names
and apostrophes use the same normalization as generation; blocked phrases have no
four-token cutoff. This avoids interpreting `hadi` inside `abdul-hadi` as a separate
entry while still rejecting `ahmad-cow-olive` and `sad-ibn-abi-waqqas-cow`.

The SQL exporter validates before any Wrangler setup and validates again while
writing. It publishes a dump atomically only after the entire corpus passes, so
failed input leaves an existing dump intact. Its SQL omits transaction statements,
as required by [Cloudflare's D1 import documentation](https://developers.cloudflare.com/d1/best-practices/import-export-data/).
Offline SQL execution does not prove a remote deployment or reservation service.

These rules enforce the project's reviewed list policy. They do not certify every
possible cultural interpretation of a random phrase. Outputs are labels, can
repeat, and are not authentication secrets. Applications must enforce uniqueness
when using them as identifiers.
