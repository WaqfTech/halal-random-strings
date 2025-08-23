// Package halalrandomstrings provides a human-readable random string generator.
package halalrandomstrings

import (
	crand "crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/x/exp/ordered"
)

//go:embed words.json
var wordsData []byte

const (
	defaultPrefixThreshold = 0.2
	defaultSuffixThreshold = 0.2
	maxRetries             = 1000 // Increased max retries
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
		opts.MinWords = 5 // Default to 5 words
	}

	if opts.MaxWords < opts.MinWords {
		opts.MaxWords = 8 // Default to 8 words
	}

	r := make([]string, opts.Repeat)
	
	// Use math/rand for word selection (faster)
	var src *rand.Rand
	if opts.Seed == 0 {
		src = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		src = rand.New(rand.NewSource(opts.Seed))
	}

	for i := range r {
		for j := 0; j < maxRetries; j++ {
			var currentOutput string
			var currentWordCount int

			// Build the string iteratively
			for {
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
				
				// Calculate word count of the generated part
				generatedWordCount := len(strings.Split(generated, opts.Sep))

				// Check if adding this rule would exceed MaxWords
				if currentWordCount + generatedWordCount > opts.MaxWords && currentWordCount > 0 {
					// If we have words already and adding this rule exceeds MaxWords, try to finalize
					break
				}

				// Append the generated part
				if currentOutput == "" {
					currentOutput = generated
				} else {
					currentOutput = currentOutput + opts.Sep + generated
				}
				currentWordCount = len(strings.Split(currentOutput, opts.Sep))

				// If we have enough words, break the inner loop
				if currentWordCount >= opts.MinWords {
					break
				}
			}

			// Finalize the output
			if currentWordCount >= opts.MinWords && currentWordCount <= opts.MaxWords {
				if isSafe(currentOutput) {
					output := strings.ReplaceAll(currentOutput, "'", "-")
					
					// Add a random number for extra uniqueness using crypto/rand
					randomNum, err := crand.Int(crand.Reader, big.NewInt(1000000000000000000))
					if err != nil {
						// Fallback to time-based if crypto/rand fails
						randomNum = big.NewInt(time.Now().UnixNano() % 1000000000000000000)
					}
					output = fmt.Sprintf("%s%s%d", output, opts.Sep, randomNum.Int64())

					r[i] = strings.ToLower(strings.ReplaceAll(output, " ", opts.Sep))
					break
				}
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
	})[0]
}

// GenerateN returns a given number of random strings.
func GenerateN(n int) []string {
	return generate(Options{
		Repeat: n,
	})
}

// GenerateWithOptions generates results against the given options.
func GenerateWithOptions(o Options) []string {
	return generate(o)
}
