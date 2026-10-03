# Halal Random Strings

A random identifier generator using reviewed general, personal-name, and Islamic word lists. Category groups prevent protected names and religious expressions from being combined with animals, food, or other general words.

This project aims to provide strings that are:

*   **Halal-friendly:** Utilizing Islamic terms and concepts, transliterated into Latin script.
*   **Kid-friendly:** Avoiding inappropriate or offensive language.
*   **Varied:** Produces many combinations; applications must enforce uniqueness when required.
*   **Reviewed mixing:** Each identifier uses one approved category group. Random combinations still require ongoing cultural review.

## Features

*   Generates random strings from curated lists of adjectives, nouns (objects, places, concepts), Islamic virtues, Sahaba names, and food items.
*   Checks complete outputs against blocked words and phrases.
*   Supports configurable string length (minimum 5, maximum 8 words by default).
*   Appends an optional four digit random number (deterministic when seeded). Identifiers are not authentication secrets.
*   Provides a validated SQL export script for Cloudflare D1. A serving API is not implemented.

## Project Structure

```
halal-random-strings/
├── cmd/
│   └── halal-random-strings/
│       └── main.go             # CLI entry point for string generation
├── db/
│   └── schema.sql              # Database schema for Cloudflare D1
├── docs/
│   ├── CONTRIBUTING.md         # Guidelines for contributing to the project
│   └── TODOS.md                # Future enhancements and known limitations
├── scripts/
│   ├── analyze.py              # Python script to analyze string uniqueness
│   └── populate_d1.py          # Python script to populate D1 database
├── halal-random-strings.go     # Core logic for string generation
├── halal-random-strings_test.go # Unit tests for the core logic
├── words.json                  # Centralized JSON file containing all word categories and generation rules
├── LICENSE                     # GNU Affero General Public License v3.0 (AGPL-3.0)
├── NOTICE                      # Upstream copyright and MIT license attribution
├── WaqfDPL-1.0.md              # Waqf Digital Public License (Waqf-DPL 1.0) draft
├── SECURITY.md                 # Guidelines for reporting security vulnerabilities
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

*   **Go (1.21 or later):** [Install Go](https://golang.org/doc/install)
*   **Cloudflare Wrangler CLI:** Used for interacting with Cloudflare D1 and Workers. [Install Wrangler](https://developers.cloudflare.com/workers/wrangler/install-update/)
*   **Node.js (LTS recommended):** Required for `npm` to install Wrangler. [Download Node.js](https://nodejs.org/)

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/WaqfTech/halal-random-strings.git
    cd halal-random-strings
    ```

2.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Build the Go executable:**
    ```bash
    go build -o halal-random-strings ./cmd/halal-random-strings
    ```

### Usage

#### Generating Strings (CLI)

To generate a single random string:

```bash
./halal-random-strings
```

To generate multiple strings (e.g., 10 strings):

```bash
./halal-random-strings -r 10
```

Example output:

```
olive-apple-pear-plum-peach-1234
aisha-maher-amira-fatima-layla-9876
# ... and so on
```

#### Mixing policy and errors

Asma Allah mixes only with Asma Allah; prophets only with prophets; personal names
only with personal names; protected Islamic categories only within their group;
and general categories only within theirs. Standalone servant prefixes and custom
dictionaries are disabled. Some formerly general entries were moved or removed.
See [the complete mixing policy and API migration](docs/MIXING_POLICY.md).

```bash
# Allowed general mix
./halal-random-strings --categories colors_arabic,nouns_places -r 10
# Allowed names-only mix
./halal-random-strings --categories muslim_names_male,muslim_names_female -r 10
# Rejected: protected names cannot mix with animals or food
./halal-random-strings --categories prophets,animals,arabic_food
```

Use `GenerateWithOptionsE` to receive validation errors. `NewDefaultEngine` loads
only the embedded reviewed dictionary. Separators are limited to `-`, `_`, `.`, `/`.
`ValidateString` checks the full policy; `IsSafe` checks blocked vocabulary only.

#### `words.json` Structure

The `words.json` file is the heart of the generator. It contains:

*   `categories`: A map where keys are category names (e.g., `adjectives`, `fruits`, `islamic_golden_age_scholars`) and values are arrays of strings (the words themselves).
*   `policy`: The required group classification, separator allowlist and limits, shared by generation and auditing.
*   `rules`: Validated composition examples, retained for compatibility. Runtime generation uses the approved category pools.
*   `blocked`: An array of strings that should never appear in generated names.

#### Cloudflare D1 Integration

This project is designed to pre-generate a large pool of unique strings and store them in a Cloudflare D1 database, which can then be accessed via a Cloudflare Worker API.

1.  **Authenticate Wrangler:**
    ```bash
    wrangler login
    ```

2.  **Create a D1 database:**
    ```bash
    wrangler d1 create my-halal-strings-db # Choose your desired database name
    ```

3.  **Apply the schema:**
    ```bash
    wrangler d1 execute my-halal-strings-db --file db/schema.sql
    ```

4.  **Generate a large batch of strings:**
    ```bash
    # Generates 1,000 strings into output.txt and runs uniqueness analysis
    make generate
    ```

5.  **Populate D1 database:**
    ```bash
    make populate-d1
    ```
    This script rejects unsafe input before Wrangler setup and batches validated strings. Use `--dump-sql /tmp/validated.sql` for an offline export. Python 3.9 or later is required.

A Worker serving or reservation API is future work. This repository contains the
generator, schema, analyzer and import tooling; it does not deliver that service.

> **Architectural Recommendation (Edge Workers):** On-the-fly generation is an
> option for readable labels after implementing an appropriate runtime generator.
> A preallocated D1 pool can support single-use invitations, but requires atomic
> allocation, collision handling and separate authorization checks. Neither a
> serving Worker nor those reservation guarantees are implemented here.

## Contributing

We welcome contributions! Please see `docs/CONTRIBUTING.md` for guidelines on how to set up your development environment, add new words, and propose changes.

## Security

For information on how to report security vulnerabilities, please see `SECURITY.md`.

## Future Enhancements

Check `docs/TODOS.md` for a list of planned features, potential improvements, and known limitations.

## Credit & Upstream Attribution

This project is a fork of and derivative work based on [charmbracelet/hotdiva2000](https://github.com/charmbracelet/hotdiva2000), originally authored by Charmbracelet, Inc. under the MIT License. We express our gratitude to the original authors for their open-source contributions.

For full upstream MIT copyright and permission notices, see [NOTICE](NOTICE).

## License

This project is licensed under the [GNU Affero General Public License v3.0 (AGPL-3.0)](LICENSE) with copyright (c) 2024-2026 WaqfTech.

The moral foundation, spiritual principles, and ethical terms of endowment guiding this work are governed by the [Waqf Digital Public License (Waqf-DPL 1.0)](WaqfDPL-1.0.md), drafted by WaqfTech at [github.com/WaqfTech/waqf-license-draft](https://github.com/WaqfTech/waqf-license-draft).
