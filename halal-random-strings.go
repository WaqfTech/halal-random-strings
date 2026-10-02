// Package halalrandomstrings provides a human-readable random string generator.
package halalrandomstrings

import (
	crand "crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/x/exp/ordered"
)

const (
	defaultPrefixThreshold = 0.2
	defaultSuffixThreshold = 0.2
	maxRetries             = 1000 // Increased max retries
)

//go:embed words.json
var defaultWordsData []byte

type Words struct {
	Categories map[string][]string `json:"categories"`
	Rules      []Rule              `json:"rules"`	
	Blocked    []string            `json:"blocked"`
}

type Rule struct {
	Pattern  []string `json:"pattern"`
	Template string   `json:"template"`
}

type wordToken struct {
	text   string
	tokens int
}

// Engine provides a thread-safe random string generator engine.
type Engine struct {
	mu         sync.RWMutex
	words      Words
	blockedSet map[string]struct{}
	normWords  map[string][]wordToken
}

var randPool = sync.Pool{
	New: func() any {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

func normalizeTokenString(s string) string {
	var sb strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	norm := sb.String()
	for strings.Contains(norm, "--") {
		norm = strings.ReplaceAll(norm, "--", "-")
	}
	return strings.Trim(norm, "-")
}

func (e *Engine) rebuildIndexesLocked() {
	e.blockedSet = make(map[string]struct{}, len(e.words.Blocked))
	for _, b := range e.words.Blocked {
		norm := normalizeTokenString(b)
		if norm != "" {
			e.blockedSet[norm] = struct{}{}
		}
	}

	e.normWords = make(map[string][]wordToken, len(e.words.Categories))
	for cat, list := range e.words.Categories {
		tokens := make([]wordToken, 0, len(list))
		for _, w := range list {
			norm := normalizeWord(w, "-")
			if norm != "" {
				count := strings.Count(norm, "-") + 1
				tokens = append(tokens, wordToken{
					text:   norm,
					tokens: count,
				})
			}
		}
		e.normWords[cat] = tokens
	}
}

// NewEngine creates an Engine using the given Words dataset.
func NewEngine(w Words) *Engine {
	e := &Engine{words: w}
	e.rebuildIndexesLocked()
	return e
}

// NewDefaultEngine creates an Engine populated from embedded words.json.
func NewDefaultEngine() (*Engine, error) {
	var w Words
	if len(defaultWordsData) > 0 {
		if err := json.Unmarshal(defaultWordsData, &w); err != nil {
			return nil, fmt.Errorf("failed to parse embedded words.json: %w", err)
		}
	}
	e := &Engine{words: w}
	e.rebuildIndexesLocked()
	return e, nil
}


var (
	defaultEngine *Engine
	words         Words // Populated by default from embedded words.json, overridable by LoadWords
)

func init() {
	var err error
	defaultEngine, err = NewDefaultEngine()
	if err != nil {
		panic(fmt.Sprintf("failed to parse embedded words.json: %v", err))
	}
	words = defaultEngine.words
}

// DefaultEngine returns the package-level default engine.
func DefaultEngine() *Engine {
	return defaultEngine
}

// LoadWords reads words.json from the specified path and populates the engine words.
func (e *Engine) LoadWords(filePath string) error {
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open words.json: %w", err)
	}

	var newWords Words
	if err := json.Unmarshal(byteValue, &newWords); err != nil {
		return fmt.Errorf("failed to parse words.json: %w", err)
	}

	e.mu.Lock()
	e.words = newWords
	e.rebuildIndexesLocked()
	e.mu.Unlock()
	return nil
}

// LoadWords reads words.json from the specified path and updates defaultEngine and package-level words.
func LoadWords(filePath string) error {
	if err := defaultEngine.LoadWords(filePath); err != nil {
		return err
	}
	defaultEngine.mu.RLock()
	words = defaultEngine.words
	defaultEngine.mu.RUnlock()
	return nil
}

func (e *Engine) isSafeLocked(s string) bool {
	norm := normalizeTokenString(s)
	if norm == "" {
		return true
	}

	tokens := strings.Split(norm, "-")
	n := len(tokens)
	// Fast-path: single token O(1) checks without any string building
	for i := 0; i < n; i++ {
		if _, blocked := e.blockedSet[tokens[i]]; blocked {
			return false
		}
	}
	// Multi-token phrase checks
	for i := 0; i < n; i++ {
		var sb strings.Builder
		sb.WriteString(tokens[i])
		for j := i + 1; j < n && j < i+4; j++ {
			sb.WriteString("-")
			sb.WriteString(tokens[j])
			phrase := sb.String()
			if _, blocked := e.blockedSet[phrase]; blocked {
				return false
			}
		}
	}
	return true
}

// IsSafe checks whether a string is safe against blocked words.
func (e *Engine) IsSafe(s string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isSafeLocked(s)
}

func isSafe(s string) bool {
	return defaultEngine.IsSafe(s)
}

func normalizeWord(s, sep string) string {
	s = strings.ReplaceAll(s, "'", sep)
	if sep != "-" {
		s = strings.ReplaceAll(s, "-", sep)
	}
	s = strings.ToLower(strings.ReplaceAll(s, " ", sep))
	for strings.Contains(s, sep+sep) {
		s = strings.ReplaceAll(s, sep+sep, sep)
	}
	return strings.Trim(s, sep)
}

func (e *Engine) getCategoriesLocked() []string {
	keys := make([]string, 0, len(e.words.Categories))
	for k := range e.words.Categories {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// GenerateWithOptionsE returns random strings generated by the engine with explicit error reporting.
func (e *Engine) GenerateWithOptionsE(opts Options) ([]string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	engineWords := e.words

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
		if opts.MinWords > 8 {
			opts.MaxWords = opts.MinWords
		} else {
			opts.MaxWords = 8 // Default to 8 words
		}
	}

	// Validate categories once before generation loops
	var filteredRules []Rule
	if len(opts.Categories) > 0 {
		for _, cat := range opts.Categories {
			cat = strings.TrimSpace(cat)
			if _, exists := engineWords.Categories[cat]; !exists {
				return nil, fmt.Errorf("unknown category: %q (available: %s)", cat, strings.Join(e.getCategoriesLocked(), ", "))
			}
		}

		// Compound rule matching
		var matchingRules []Rule
		for _, rule := range engineWords.Rules {
			allCategoriesMatch := true
			for _, patternCategory := range rule.Pattern {
				found := false
				for _, selectedCategory := range opts.Categories {
					if patternCategory == strings.TrimSpace(selectedCategory) {
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
				matchingRules = append(matchingRules, rule)
			}
		}

		if len(opts.Categories) > 1 {
			var fullMatchRules []Rule
			for _, rule := range matchingRules {
				hasAll := true
				for _, selCat := range opts.Categories {
					catFound := false
					for _, patternCat := range rule.Pattern {
						if patternCat == strings.TrimSpace(selCat) {
							catFound = true
							break
						}
					}
					if !catFound {
						hasAll = false
						break
					}
				}
				if hasAll {
					fullMatchRules = append(fullMatchRules, rule)
				}
			}
			if len(fullMatchRules) > 0 {
				filteredRules = fullMatchRules
			} else {
				filteredRules = matchingRules
			}
		} else {
			filteredRules = matchingRules
		}


		// Dynamic single-category rule fallback if no compound rules match
		if len(filteredRules) == 0 {
			for _, cat := range opts.Categories {
				cat = strings.TrimSpace(cat)
				filteredRules = append(filteredRules, Rule{
					Pattern:  []string{cat},
					Template: fmt.Sprintf("{%s}", cat),
				})
			}
		}
	} else {
		filteredRules = engineWords.Rules
	}

	if len(filteredRules) == 0 {
		return nil, fmt.Errorf("no rules available for specified options")
	}

	r := make([]string, opts.Repeat)
	
	// Use math/rand for word selection (faster)
	var src *rand.Rand
	if opts.Seed == 0 {
		src = randPool.Get().(*rand.Rand)
		defer randPool.Put(src)
	} else {
		src = rand.New(rand.NewSource(opts.Seed))
	}

	for i := range r {
		for j := 0; j < maxRetries; j++ {
			var currentBuilder strings.Builder
			var currentWordCount int

			// Build the string iteratively
			for attempt := 0; attempt < maxRetries; attempt++ {
				// Pick a random rule from filtered rules
				rule := filteredRules[src.Intn(len(filteredRules))]

				var ruleBuilder strings.Builder
				var ruleWordCount int
				for _, category := range rule.Pattern {
					categoryWords := e.normWords[category]
					if len(categoryWords) == 0 {
						continue
					}
					wt := categoryWords[src.Intn(len(categoryWords))]
					var normWord string
					if opts.Sep == "-" {
						normWord = wt.text
					} else {
						normWord = strings.ReplaceAll(wt.text, "-", opts.Sep)
					}
					if normWord == "" {
						continue
					}
					if ruleBuilder.Len() > 0 {
						ruleBuilder.WriteString(opts.Sep)
					}
					ruleBuilder.WriteString(normWord)
					ruleWordCount += wt.tokens
				}
				if ruleBuilder.Len() == 0 {
					continue
				}

				// Check if adding this rule would exceed MaxWords
				if currentWordCount+ruleWordCount > opts.MaxWords {
					// If we have enough words already, finalize this string
					if currentWordCount >= opts.MinWords {
						break // Break from inner building loop
					} else {
						// Not enough words yet, and this rule exceeds MaxWords. Reset and try again.
						currentBuilder.Reset()
						currentWordCount = 0
						continue // Continue to next attempt in inner building loop
					}
				}

				// Append the generated part
				if currentBuilder.Len() > 0 {
					currentBuilder.WriteString(opts.Sep)
				}
				currentBuilder.WriteString(ruleBuilder.String())
				currentWordCount += ruleWordCount

				// If we have enough words, and haven't exceeded MaxWords, we can break
				if currentWordCount >= opts.MinWords && currentWordCount <= opts.MaxWords {
					break // Break from inner building loop
				}
			}

			// Finalize the output if it meets criteria
			if currentWordCount >= opts.MinWords && currentWordCount <= opts.MaxWords {
				candidate := currentBuilder.String()
				if e.isSafeLocked(candidate) {
					if opts.IncludeRandomNumber {
						var val int64
						if opts.Seed != 0 {
							val = int64(src.Intn(9000)) + 1000
						} else {
							var b [2]byte
							_, err := crand.Read(b[:])
							if err != nil {
								val = int64(src.Intn(9000)) + 1000
							} else {
								val = int64((uint16(b[0])<<8|uint16(b[1]))%9000) + 1000
							}
						}
						currentBuilder.WriteString(opts.Sep)
						currentBuilder.WriteString(strconv.FormatInt(val, 10))
						candidate = currentBuilder.String()
					}

					r[i] = candidate
					break // Break from maxRetries loop
				}
			}
		}
	}

	for i, s := range r {
		if s == "" {
			return nil, fmt.Errorf("failed to generate valid string for index %d within retry limit", i)
		}
	}

	return r, nil
}

// GenerateWithOptions returns random strings generated by the engine.
func (e *Engine) GenerateWithOptions(opts Options) []string {
	r, _ := e.GenerateWithOptionsE(opts)
	return r
}

// Options are options to customize output.
type Options struct {
	// PrefixThreshold is reserved for upstream compatibility and currently unused.
	PrefixThreshold float64
	// SuffixThreshold is reserved for upstream compatibility and currently unused.
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
func (e *Engine) Generate() string {
	res := e.GenerateWithOptions(Options{
		Repeat:              1,
		IncludeRandomNumber: true,
	})
	if len(res) == 0 {
		return ""
	}
	return res[0]
}

// GenerateN returns n random strings.
func (e *Engine) GenerateN(n int) []string {
	return e.GenerateWithOptions(Options{
		Repeat:              n,
		IncludeRandomNumber: true,
	})
}

// GetCategories returns a sorted list of available categories.
func (e *Engine) GetCategories() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.getCategoriesLocked()
}

// Generate returns a random string using the default engine.
func Generate() string {
	return defaultEngine.Generate()
}

// GenerateN returns a given number of random strings using the default engine.
func GenerateN(n int) []string {
	return defaultEngine.GenerateN(n)
}

// GenerateWithOptionsE generates results against the given options using the default engine with explicit error reporting.
func GenerateWithOptionsE(o Options) ([]string, error) {
	return defaultEngine.GenerateWithOptionsE(o)
}

// GenerateWithOptions generates results against the given options using the default engine.
func GenerateWithOptions(o Options) []string {
	return defaultEngine.GenerateWithOptions(o)
}

// GetCategories returns a sorted list of available categories.
func GetCategories() []string {
	return defaultEngine.GetCategories()
}