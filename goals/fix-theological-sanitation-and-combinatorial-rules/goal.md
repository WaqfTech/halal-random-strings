# Goal: Sanitize Theological Terms and Prevent Disrespectful Word Combinations

## Goal Description
Purge theological hazards, sacrilegious rules, and disrespectful word combinations from `words.json`:
1. Remove Rule 11 pairing Islamic Golden Age scholars with vegetables (`ibn-sina-cucumber`).
2. Remove Rule 43 pairing animals with adjectives containing Prophetic titles/moral honors (`ape-muzammil`, `donkey-amin`, `hyena-faruq`).
3. Remove Rules 41 & 42 pairing Arabic colors with personal human names to avoid racially charged descriptors (`aswad-fatima`, `asfar-omar`).
4. Purge malevolent concepts (`dajjal`, `jahannam`, `fitna`, `yajuj-majuj`) from `nouns_concepts` to prevent generation of blessed evil concepts (`mubarak-dajjal`, `sadiq-jahannam`, `dajjal-hub`).
5. Complete sanctuary and sacred text segregation by removing remaining sanctuaries and Quranic entities (`kaba-structure`, `hajar-al-aswad`, `maqam-ibrahim`, `hijr-isma'il`, `masjid-al-haram`, `masjid-al-nabawi`, `zamzam-well`, `safa-marwah-mounts`, `quran`, `ayats`, `surahs`) from `nouns_concepts`.
6. Fix unprefixed Divine Names in `muslim_names_male` (`qadir`, `ghani`, `wahid`, `mu'izz`, etc.) and correct culinary transliteration `fattah` to `fatteh`/`fatta`.
7. Eliminate eschatological error `jannah-end` from `suffixes`.

## Dependencies & Execution Order
- **Mode**: Independent ⚡
- **Depends On**: -
- **Sequence**: 15
- **Shape**: ship
- **Tier**: roadmap

## References
- **Shared Understanding & Fact Sheet**: [`goals/fix-theological-sanitation-and-combinatorial-rules/facts.md`](facts.md)
- **Execution Plan**: [`goals/fix-theological-sanitation-and-combinatorial-rules/plan.md`](plan.md)

## Done Condition
1. Rules 11, 41, 42, and 43 are either safely removed or replaced with culturally respectful patterns.
2. `nouns_concepts` is fully purged of malevolent concepts (`dajjal`, `jahannam`, `fitna`, `yajuj-majuj`) and sacred sanctuaries/texts (`kaba-structure`, `quran`, etc.).
3. Unprefixed Divine Names are removed or properly prefixed in `muslim_names_male`.
4. Automated tests verify zero occurrences of disrespectful pairings and complete category sanitation.
5. All tests pass with `go test -race ./...`.
