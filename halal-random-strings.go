// Package halalrandomstrings provides a human-readable random string generator.
package halalrandomstrings

import (
	_ "embed"
	"math/rand"
	"strings"

	"github.com/charmbracelet/x/exp/ordered"
)

//go:generate sort -u modifiers.txt -o modifiers.txt
//go:generate sort -u nouns.txt -o nouns.txt
//go:generate sort -u prefix.txt -o prefix.txt
//go:generate sort -u suffix.txt -o suffix.txt
//go:generate sort -u blocked.txt -o blocked.txt

const (
	defaultPrefixThreshold = 0.2
	defaultSuffixThreshold = 0.2
	maxRetries             = 100
)

//go:embed prefix.txt
var prefixData string

//go:embed modifiers.txt
var modifierData string

//go:embed nouns.txt
var nounData string

//go:embed suffix.txt
var suffixData string

//go:embed blocked.txt
var blockedData string

var (
	prefixes  []string
	modifiers []string
	nouns     []string
	suffixes  []string
	blocked   []string
)

func init() {
	prefixes = strings.Split(strings.TrimSpace(prefixData), "\n")
	modifiers = strings.Split(strings.TrimSpace(modifierData), "\n")
	nouns = strings.Split(strings.TrimSpace(nounData), "\n")
	suffixes = strings.Split(strings.TrimSpace(suffixData), "\n")
	blocked = strings.Split(strings.TrimSpace(blockedData), "\n")
}

func isSafe(s string) bool {
	for _, blockedWord := range blocked {
		if strings.Contains(s, blockedWord) {
			return false
		}
	}
	return true
}

// generate returns random strings.
func generate(opts Options) []string {
	if opts.Repeat < 1 {
		opts.Repeat = 1
	}
	opts.PrefixThreshold = ordered.Clamp(opts.PrefixThreshold, 0, 1)
	opts.SuffixThreshold = ordered.Clamp(opts.SuffixThreshold, 0, 1)

	if opts.Sep == "" {
		opts.Sep = "-"
	}

	r := make([]string, opts.Repeat)
	src := rand.New(rand.NewSource(opts.Seed))

	for i := range r {
		for j := 0; j < maxRetries; j++ {
			var (
				prefix = ""
				suffix = ""
			)

			if opts.PrefixThreshold > 0 && src.Float64() < opts.PrefixThreshold {
				prefix = prefixes[src.Intn(len(prefixes))] + " "
			}
			if opts.SuffixThreshold > 0 && src.Float64() < opts.SuffixThreshold {
				suffix = " " + suffixes[src.Intn(len(suffixes))]
			}

			mod := modifiers[src.Intn(len(modifiers))]
			noun := nouns[src.Intn(len(nouns))]

			var builder strings.Builder
			builder.WriteString(prefix)
			builder.WriteString(mod)
			builder.WriteString(" ")
			builder.WriteString(noun)
			builder.WriteString(suffix)

			output := strings.ToLower(strings.ReplaceAll(builder.String(), " ", opts.Sep))
			if isSafe(output) {
				r[i] = output
				break
			}
		}
	}

	return r
}

// Options are options to customize output.
type Options struct {
	// Whether to show occasional additional prefix and suffix content. This
	// increases possibilities but can make strings longer.
	PrefixThreshold float64
	SuffixThreshold float64

	// Number of strings to generate.
	Repeat int

	// Separator to use between words.
	Sep string

	// Seed for the random number generator.
	Seed int64
}

// Generate returns a random string.
func Generate() string {
	return generate(Options{
		PrefixThreshold: defaultPrefixThreshold,
		SuffixThreshold: defaultSuffixThreshold,
		Repeat:          1,
		Seed:            rand.Int63(),
	})[0]
}

// GenerateN returns a given number of random strings.
func GenerateN(n int) []string {
	return generate(Options{
		PrefixThreshold: defaultPrefixThreshold,
		SuffixThreshold: defaultSuffixThreshold,
		Repeat:          n,
		Seed:            rand.Int63(),
	})
}

// Possibilities returns the number of possible strings produced.
func Possibilities() (int, int) {
	low := len(modifiers) * len(nouns)
	high := low * len(prefixes) * len(suffixes)
	return low, high
}

// GenerateWithOptions generates results against the given options.
func GenerateWithOptions(o Options) []string {
	return generate(o)
}

// PossibilitiesWithOptions returns the number of possible strings produced
// against the given options.
func PossibilitiesWithOptions(o Options) (int, int) {
	low := len(modifiers) * len(nouns)
	high := low

	o.PrefixThreshold = ordered.Clamp(o.PrefixThreshold, 0, 1)
	o.SuffixThreshold = ordered.Clamp(o.SuffixThreshold, 0, 1)

	if o.PrefixThreshold >= 1 {
		low *= len(prefixes)
	}
	if o.SuffixThreshold >= 1 {
		low *= len(suffixes)
	}

	if o.PrefixThreshold > 0 {
		high *= len(prefixes)
	}
	if o.SuffixThreshold > 0 {
		high *= len(suffixes)
	}

	return low, high
}
