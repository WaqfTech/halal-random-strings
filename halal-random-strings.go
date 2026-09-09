// Package halalrandomstrings provides a human-readable random string generator.
package halalrandomstrings

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil" // Added
	"math/big"
	"math/rand"
	"os"        // Added
	"strings"
	"time"
	sort "sort"

	"github.com/charmbracelet/x/exp/ordered"
)

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

var words Words // Keep as package-level, but will be populated by LoadWords

// LoadWords reads words.json from the specified path and populates the words variable.
func LoadWords(filePath string) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open words.json: %w", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &words); err != nil {
		return fmt.Errorf("failed to parse words.json: %w", err)
	}
	return nil
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
			
			// Filter rules based on selected categories
			filteredRules := []Rule{}
			if len(opts.Categories) > 0 {
				for _, rule := range words.Rules {
					// Check if all categories in the rule's pattern are in the selected categories
					allCategoriesMatch := true
					for _, patternCategory := range rule.Pattern {
						found := false
						for _, selectedCategory := range opts.Categories {
							if patternCategory == selectedCategory {
								found = true
								break
							}
						}
						if !found {
							allCategoriesMatch = false
							break
						}
					}
					if allCategoriesMatch {
						filteredRules = append(filteredRules, rule)
					}
				}
			} else {
				filteredRules = words.Rules // Use all rules if no categories are specified
			}

			// If no rules match the selected categories, skip this iteration
			if len(filteredRules) == 0 {
				continue
			}

			// Build the string iteratively
			for attempt := 0; attempt < maxRetries; attempt++ {
				// Pick a random rule from filtered rules
				rule := filteredRules[src.Intn(len(filteredRules))]

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
				if currentWordCount + generatedWordCount > opts.MaxWords {
					// If we have enough words already, finalize this string
					if currentWordCount >= opts.MinWords {
						break // Break from inner building loop
					} else {
						// Not enough words yet, and this rule exceeds MaxWords. Reset and try again.
						currentOutput = ""
						currentWordCount = 0
						continue // Continue to next attempt in inner building loop
					}
				}

				// Append the generated part
				if currentOutput == "" {
					currentOutput = generated
				} else {
					currentOutput = currentOutput + opts.Sep + generated
				}
				currentWordCount = len(strings.Split(currentOutput, opts.Sep))

				// If we have enough words, and haven't exceeded MaxWords, we can break
				// This condition is crucial to ensure we don't keep adding words unnecessarily
				// This is the problematic line, it should be `currentWordCount >= opts.MinWords`
				// and then check `currentWordCount <= opts.MaxWords` after the loop
				if currentWordCount >= opts.MinWords && currentWordCount <= opts.MaxWords {
					break // Break from inner building loop
				}
			}

			// Finalize the output if it meets criteria
			if currentWordCount >= opts.MinWords && currentWordCount <= opts.MaxWords {
				// Before returning, ensure the string is safe
				// Normalize the output before checking safety and returning
				finalOutput := strings.ReplaceAll(currentOutput, "'", "-")
				finalOutput = strings.ToLower(strings.ReplaceAll(finalOutput, " ", opts.Sep))

				if isSafe(finalOutput) {
					// Add a random number for extra uniqueness using crypto/rand
					output := finalOutput
					if opts.IncludeRandomNumber {
						randomNum, err := crand.Int(crand.Reader, big.NewInt(1000000000000000000))
						if err != nil {
							// Fallback to time-based if crypto/rand fails
							randomNum = big.NewInt(time.Now().UnixNano() % 1000000000000000000)
						}
						output = fmt.Sprintf("%s%s%d", finalOutput, opts.Sep, randomNum.Int64())
					}

					r[i] = output
					break // Break from maxRetries loop
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

	// Whether to include a random number at the end of the string.
	IncludeRandomNumber bool

	// List of categories to use for string generation.
	Categories []string
}

// Generate returns a random string.
func Generate() string {
	return generate(Options{
		Repeat: 1,
		IncludeRandomNumber: true, // Default to true
	})[0]
}

// GenerateN returns a given number of random strings.
func GenerateN(n int) []string {
	return generate(Options{
		Repeat: n,
		IncludeRandomNumber: true, // Default to true
	})
}

// GenerateWithOptions generates results against the given options.
func GenerateWithOptions(o Options) []string {
	return generate(o)
}

// GetCategories returns a sorted list of available categories.
func GetCategories() []string {
	keys := make([]string, 0, len(words.Categories))
	for k := range words.Categories {
		keys = append(keys, k)
	}
	// Sort keys for consistent output
	sort.Strings(keys)
	return keys
}