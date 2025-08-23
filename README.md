# Halal Random Strings

A wholesome, family-friendly random name generator (inspired by hotdiva2000).

Example output:

```
$ halal-random-strings -r 3
nur-sabr-masjid
sadiq-hikma-kitab
karim-ukhuwa-safina
```

## Usage

Use it as a library:

```go
import "github.com/WaqfTech/halal-random-strings"

fmt.Println(halalrandomstrings.Generate()) // nur-sabr-masjid
```

Use it on the CLI:

```bash
# Generate a random string
halal-random-strings

# Generate 5 random strings
halal-random-strings -r 5

# Use a different separator
halal-random-strings --sep _

# Use a seed for reproducibility
halal-random-strings --seed 42

# See also
halal-random-strings -h
```

## Credit

This project is a fork of [charmbracelet/hotdiva2000](https://github.com/charmbracelet/hotdiva2000), and we are grateful to the original authors for their work.

## License

[MIT](https://github.com/charmbracelet/hotdiva2000/raw/main/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" width="400" src="https://stuff.charm.sh/charm-badge.jpg" /></a>

Charm 热爱开源 • Charm loves open source
