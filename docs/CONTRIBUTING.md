# Contributing to Halal Random Strings

We welcome contributions to the Halal Random Strings project! By contributing, you help us create a more comprehensive, diverse, and useful tool for generating kid-friendly and halal-friendly random strings.

## How to Contribute

### 1. Set Up Your Development Environment

To get started, you'll need:

*   **Go (1.21 or later):** Follow the official Go installation guide: [https://golang.org/doc/install](https://golang.org/doc/install)
*   **Cloudflare Wrangler CLI:** Used for interacting with Cloudflare D1. Install it globally:
    ```bash
    npm install -g wrangler
    ```
*   **Node.js (LTS recommended):** Required for `npm` to install Wrangler. Download from [https://nodejs.org/](https://nodejs.org/)

### 2. Clone the Repository

```bash
git clone https://github.com/WaqfTech/halal-random-strings.git
cd halal-random-strings
```

### 3. Install Go Dependencies

```bash
go mod tidy
```

### 4. Run Tests

Before making any changes, ensure all existing tests pass:

```bash
go test ./...
```

### 5. Add New Words and Categories

The core of this project is the `words.json` file. To contribute new words or categories:

*   **Locate `words.json`:** It's in the project root directory.
*   **Understand the Structure:**
    *   `categories`: A dictionary where keys are category names (e.g., `adjectives`, `fruits`, `islamic_golden_age_scholars`) and values are arrays of strings (the words themselves).
    *   `rules`: An array of objects, each defining a pattern (which categories to combine) and a template for how to combine them (e.g., `"{adjectives}-{nouns_concepts}"`).
    *   `blocked`: An array of strings that should never appear in generated names.
*   **Adding New Words:** Simply add your new words to the appropriate category array in `words.json`. Maintain alphabetical order within categories for readability.
*   **Adding New Categories:** If your words don't fit existing categories, create a new key-value pair under `categories`. Remember to also add new rules under the `rules` array if you want these new categories to be used in string generation.
*   **Maintain Halal and Kid-Friendly Principles:** Ensure all new words and combinations adhere to the project's core values: positive, encouraging, inviting, inspiring, and suitable for children and a faithful audience.

### 6. Propose Changes

1.  **Create a new branch:**
    ```bash
    git checkout -b feature/your-feature-name
    ```
2.  **Make your changes.**
3.  **Test your changes:** Run `go test ./...` and ensure everything passes.
4.  **Format your code:**
    ```bash
    go fmt ./...
    ```
5.  **Commit your changes:** Use clear and concise commit messages. Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification (e.g., `feat: add new fruit words`, `fix: resolve build issue`).
6.  **Push your branch:**
    ```bash
    git push origin feature/your-feature-name
    ```
7.  **Open a Pull Request:** Go to the GitHub repository and open a pull request from your branch to the `main` branch. Describe your changes clearly.

## Code Style

*   Follow standard Go formatting (`go fmt`).
*   Keep code clear, concise, and well-commented where necessary.

## Reporting Issues

If you find a bug or have a feature request, please open an issue on the GitHub repository.

Thank you for your contributions!