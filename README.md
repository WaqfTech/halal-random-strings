# Halal Random Strings

A wholesome, family-friendly random string generator designed for creating unique and meaningful identifiers, particularly suitable for applications like game lobby invitations.

This project aims to provide strings that are:

*   **Halal-friendly:** Utilizing Islamic terms and concepts, transliterated into Latin script.
*   **Kid-friendly:** Avoiding inappropriate or offensive language.
*   **Unique:** Designed to generate a high volume of distinct strings for reliable identification.
*   **Meaningful:** Combining words in a way that can evoke positive and inspiring associations.

## Features

*   Generates random strings from curated lists of adjectives, nouns (objects, places, concepts), Islamic virtues, Sahaba names, and food items.
*   Ensures strings are free from a comprehensive list of blocked words.
*   Supports configurable string length (minimum 5, maximum 8 words by default).
*   Appends a cryptographically secure random number for enhanced uniqueness.
*   Integrates with Cloudflare D1 for persistent storage and API access.

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
│   ├── populate_d1.py          # Python script to populate D1 database
│   └── test_uniqueness.sh      # Bash script to test uniqueness at scale
├── halal-random-strings.go     # Core logic for string generation
├── halal-random-strings_test.go # Unit tests for the core logic
├── words.json                  # Centralized JSON file containing all word categories and generation rules
├── LICENSE                     # Project license (MIT)
├── Waqf.en.md                  # Full text of the "Waqf" General Public License
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
al-khansa-salih-siraj-123456789
khalid-ibn-al-walid-haleem-wasif-987654321
# ... and so on
```

#### `words.json` Structure

The `words.json` file is the heart of the generator. It contains:

*   `categories`: A map where keys are category names (e.g., `adjectives`, `fruits`, `islamic_golden_age_scholars`) and values are arrays of strings (the words themselves).
*   `rules`: An array of objects, each defining a `pattern` (which categories to combine) and a `template` for how to combine them (e.g., `"{adjectives}-{nouns_concepts}"`).
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
    # Generates 100,000 strings into output.txt and runs uniqueness analysis
    make generate
    ```

5.  **Populate D1 database:**
    ```bash
    make populate-d1
    ```
    This script reads `output.txt` and inserts the strings into your D1 database. It includes batching for efficiency and basic error handling.

6.  **Deploy Cloudflare Worker API:**
    Refer to the `api_worker.js` (or `api_worker.ts`) outline in the project for setting up your Cloudflare Worker to serve strings from D1 via API endpoints like `/invite/new` and `/invite/status`.

> [!TIP]
> **Architectural Recommendation (Edge Workers):**
> For high-throughput production services, consider generating halal strings **on-the-fly** directly inside the Cloudflare Worker (by embedding the JSON dictionary or compiling the Go generator to WebAssembly) rather than reading from a pre-populated D1 table. On-the-fly generation eliminates D1 read-after-write concurrency races, removes database query latency, and avoids D1 row read/write billing overhead. Pre-populating D1 is best reserved for pre-allocated claim or single-use invitation tokens that require durable tracking.

## Contributing

We welcome contributions! Please see `docs/CONTRIBUTING.md` for guidelines on how to set up your development environment, add new words, and propose changes.

## Security

For information on how to report security vulnerabilities, please see `SECURITY.md`.

## Future Enhancements

Check `docs/TODOS.md` for a list of planned features, potential improvements, and known limitations.

## Credit

This project is a fork of [charmbracelet/hotdiva2000](https://github.com/charmbracelet/hotdiva2000), and we are grateful to the original authors for their work.

## License

This project is licensed under the [MIT License](LICENSE). The core principles and moral intentions behind this project are further elaborated in the ["Waqf" General Public License](Waqf.en.md).

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" width="400" src="https://stuff.charm.sh/charm-badge.jpg" /></a>

Charm 热爱开源 • Charm loves open source
