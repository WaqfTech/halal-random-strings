# Goal: Sanitize Theology And Corpus Audit

## Goal Description
Audit corpus generation at scale (130,000 strings), isolate sacred entities from mundane objects/animals/food, sanitize dictionary scraping artifacts, fix apostrophe transliteration, and optimize engine throughput.

## Implementation Details
1. **Theological Quarantine & Sanctity Protection**:
   - Purged Holy Qur'an and sacred items (`mushaf`, `quran-stand`, `sajjadah`, `masbaha`, `prayer-beads`, `adhan-clock`, `sajada`) and modesty attire (`hijab`, `niqab`, `thobe`, `miswak`) from `nouns_objects` (a mundane category).
   - Moved sacred religious objects into `holy_sanctuaries`.
   - Classified `muslim_empires` into `sacredCategories` to strictly isolate Caliphates from animals, poultry, and food.
   - Added `rahim` (`al-rahim`) to `asma_allah`.
   - Purged divisive sectarian entries (`Yazid`, `Ghadir Khumm`).
   - Purged Western/Christian transliterations (`Christ's Thorn Jujube`) and inappropriate inventions (`Shampoo`).

2. **Scraped Dictionary Artifact Purge**:
   - Replaced dual English-Arabic joined glossary entries (`peace-salam`, `justice-adl`, `courage-shaja-a`, `favor-ni-ma`) with pure concepts.
   - Stripped glossary category metadata headers (`judgment-day-concept`, `badr-event`, `caliphate-system`, `rahman-attribute`, `fajr-dawn`, `arafat-day`, `ruh-cleanliness`).
   - Cleaned definition glosses in `holy_sanctuaries` (`haram-area-sanctuary` -> `al-haram`, `jameh-mosque-great` -> `jami-masjid`, `mussalla-prayer` -> `musalla`).
   - Cleaned compound descriptors in `suffixes` and `nouns_objects`.

3. **Core Engine Normalization & Optimization**:
   - Updated `normalizeWord` to strip apostrophes instead of turning them into token delimiters (`sa'd` -> `sad`, `mu'min` -> `mumin`, `sa'eed` -> `saeed`).
   - Fixed word count tokenization so single short transliterated phrases do not falsely satisfy `MinWords`.
   - Fixed generator loop discard logic to retry alternative rules instead of destructively wiping accumulated string builders.
   - Pre-computed `defaultSacredCandidates` and `defaultMundaneCandidates` on `Engine` during index rebuilding, cutting generation latency by 63% (from 5.4µs to 2.0µs) and heap allocations by 89% (from 5,432 B/op to 605 B/op).
   - Raised CLI clamp limit from 100,000 to 1,000,000 in `cmd/halal-random-strings/main.go`.

4. **Continuous Quality Gate & Semantic Audit**:
   - Upgraded `scripts/analyze.py` with full semantic audit reporting: checks Holy entity sanctity, cross-domain isolation, duplicate frequencies, length/word distributions, and dictionary scraping artifacts.
   - Added `verify-corpus` target to `Makefile`.
   - Added `TestGoal_21_SanitizeTheologyAndCorpusAudit` to `goal_delivery_test.go`.

## Done Condition
1. All requirements described in the goal and execution plan are implemented.
2. All automated unit and integration tests pass cleanly with `-race` (`go test -race ./...`).
3. Zero architectural regressions; definitions of done satisfied.
