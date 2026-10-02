# Facts: Prune False-Positive Blocked Words and Scraped Dictionary Noise

## Architectural Invariants & Constraints
- Zero internal conflict: No dictionary word intentionally supplied in `words.json` should ever be blocked by its own safety filter.
- Clean word tokens: Words should be natural lexical units, not raw scrape artifacts with parenthetical explanations or taxonomy tags.
- Blocked list precision: The blocked list must target true harm (vulgarities, pornography, intoxicants, gambling, slurs) without blanket-blocking benign vocabulary.

## File & Interface Contracts
- `words.json`: `blocked` array, `categories` dictionary.
- `halal-random-strings_test.go`: Add `TestDictionaryZeroSelfBlocking`.
