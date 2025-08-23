// Package halalrandomstrings provides a human-readable random string generator.
package halalrandomstrings

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/x/exp/ordered"
)

//go:embed words.json
var wordsData []byte

const (
	defaultPrefixThreshold = 0.2
	defaultSuffixThreshold = 0.2
	maxRetries             = 100
)

type Words struct {
	Categories map[string][]string `json:"categories"`
	Rules      []Rule              `json:"rules"`
	Blocked    []string            `json:"blocked"`
}

type Rule struct {
	Pattern  []string `json:"pattern"`
	Template string   `json:"template"`
}

var words Words

func init() {
	if err := json.Unmarshal(wordsData, &words); err != nil {
		panic(fmt.Sprintf("failed to parse words.json: %v", err))
	}
}

func isSafe(s string) bool {
	for _, blockedWord := range words.Blocked {
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

	if opts.MinWords < 1 {
		opts.MinWords = 4
	}

	if opts.MaxWords < opts.MinWords {
		opts.MaxWords = opts.MinWords
	}

	r := make([]string, opts.Repeat)
	src := rand.New(rand.NewSource(opts.Seed))

	for i := range r {
		for j := 0; j < maxRetries; j++ {
			var output string
			var wordCount int

			for wordCount < opts.MinWords {
				// Pick a random rule
				rule := words.Rules[src.Intn(len(words.Rules))]

				// Generate the string based on the rule
				var replacerArgs []string
				for _, category := range rule.Pattern {
					word := words.Categories[category][src.Intn(len(words.Categories[category]))]
					replacerArgs = append(replacerArgs, fmt.Sprintf("{%s}", category), word)
				}

				replacer := strings.NewReplacer(replacerArgs...)
				generated := replacer.Replace(rule.Template)

				if output == "" {
					output = generated
				} else {
					output = output + opts.Sep + generated
				}

				wordCount = len(strings.Split(output, opts.Sep))

				if wordCount > opts.MaxWords {
					output = ""
					wordCount = 0
				}
			}

			if isSafe(output) {
				output = strings.ReplaceAll(output, "'", "-")
				r[i] = strings.ToLower(strings.ReplaceAll(output, " ", opts.Sep))
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

	// Minimum number of words in the generated string.
	MinWords int

	// Maximum number of words in the generated string.
	MaxWords int
}

// Generate returns a random string.
func Generate() string {
	return generate(Options{
		Repeat: 1,
		Seed:   rand.Int63(),
	})[0]
}

// GenerateN returns a given number of random strings.
func GenerateN(n int) []string {
	return generate(Options{
		Repeat: n,
		Seed:   rand.Int63(),
	})
}

// GenerateWithOptions generates results against the given options.
func GenerateWithOptions(o Options) []string {
	return generate(o)
}
