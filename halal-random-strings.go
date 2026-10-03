// Package halalrandomstrings generates identifiers from reviewed word lists.
// These identifiers are not passwords or a religious certification.
package halalrandomstrings

import (
	crand "crypto/rand"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed words.json
var defaultWordsData []byte

type Words struct {
	Categories map[string][]string `json:"categories"`
	Rules      []Rule              `json:"rules"`
	Blocked    []string            `json:"blocked"`
	Policy     MixingPolicy        `json:"policy"`
}

type Rule struct {
	Pattern  []string `json:"pattern"`
	Template string   `json:"template"`
}

// MixingPolicy is embedded with the reviewed dictionary, not user configurable.
type MixingPolicy struct {
	Version            int               `json:"version"`
	CategoryGroups     map[string]string `json:"category_groups"`
	AllowedSeparators  []string          `json:"allowed_separators"`
	DisabledCategories map[string]string `json:"disabled_categories"`
	MaxWords           int               `json:"max_words"`
	MaxRepeat          int               `json:"max_repeat"`
}

type wordToken struct {
	text   string
	tokens int
}
type wordBuckets [][]wordToken

// Engine is immutable after construction and safe for concurrent readers.
type Engine struct {
	words           Words
	blockedSet      map[string]struct{}
	blockedTrie     *blockedNode
	normWords       map[string][]wordToken
	catBuckets      map[string]wordBuckets
	groupBuckets    map[string]wordBuckets
	groupPhrases    map[string]map[string]struct{}
	maxPhraseTokens map[string]int
	defaultPlans    []generationPlan
	err             error
}

// ErrCustomDictionariesDisabled is returned for all custom dictionary requests.
var ErrCustomDictionariesDisabled = errors.New("custom dictionaries are disabled; use NewDefaultEngine and reviewed categories")

var randPool = sync.Pool{New: func() any { return rand.New(rand.NewSource(time.Now().UnixNano())) }}
var apostropheNormalizer = strings.NewReplacer("'", "", "’", "")

// NewEngine retains its old signature but returns a disabled engine.
// Deprecated: use NewDefaultEngine. GenerateWithOptionsE reports the rejection.
func NewEngine(_ Words) *Engine { return &Engine{err: ErrCustomDictionariesDisabled} }

// LoadWords rejects overrides without reading the file or changing the engine.
// Deprecated: custom dictionaries are disabled.
func (e *Engine) LoadWords(_ string) error { return ErrCustomDictionariesDisabled }

// LoadWords rejects overrides of the package default dictionary.
// Deprecated: custom dictionaries are disabled.
func LoadWords(_ string) error { return ErrCustomDictionariesDisabled }

// NewDefaultEngine loads and validates the embedded, reviewed dataset.
func NewDefaultEngine() (*Engine, error) {
	var w Words
	if err := json.Unmarshal(defaultWordsData, &w); err != nil {
		return nil, fmt.Errorf("parse embedded dictionary: %w", err)
	}
	e := &Engine{words: w}
	if err := e.buildIndexes(); err != nil {
		return nil, fmt.Errorf("validate embedded dictionary: %w", err)
	}
	var err error
	e.defaultPlans, err = e.plans(nil, 5, 8)
	if err != nil {
		return nil, err
	}
	return e, nil
}

var defaultEngine *Engine
var words Words // Private snapshot for dataset checks; not an override API.

func init() {
	var err error
	defaultEngine, err = NewDefaultEngine()
	if err != nil {
		panic(err)
	}
	// A separate copy prevents internal dataset checks from aliasing engine data.
	if err := json.Unmarshal(defaultWordsData, &words); err != nil {
		panic(err)
	}
}

func DefaultEngine() *Engine { return defaultEngine }

// Canonical apostrophe handling matches generation, blocking and corpus auditing.
func normalizeTokenString(s string) string {
	s = apostropheNormalizer.Replace(strings.ToLower(strings.TrimSpace(s)))
	var b strings.Builder
	lastDash := true
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			lastDash = false
		} else if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

func normalizeWord(s, sep string) string {
	return strings.ReplaceAll(normalizeTokenString(s), "-", sep)
}

func (e *Engine) passesBlocklist(s string) bool {
	tokens := strings.Split(normalizeTokenString(s), "-")
	for i := range tokens {
		node := e.blockedTrie
		for j := i; j < len(tokens) && node != nil; j++ {
			node = node.children[tokens[j]]
			if node != nil && node.terminal {
				return false
			}
		}
	}
	return true
}

// IsSafe checks blocked vocabulary only. Use ValidateString for the mixing policy.
func (e *Engine) IsSafe(s string) bool {
	return e != nil && e.err == nil && len(e.blockedSet) > 0 && e.passesBlocklist(s)
}
func isSafe(s string) bool { return defaultEngine.IsSafe(s) }

func (e *Engine) validSep(sep string) bool {
	for _, allowed := range e.words.Policy.AllowedSeparators {
		if sep == allowed {
			return true
		}
	}
	return false
}

// ValidateString checks formatting, blocked vocabulary and complete phrase
// segmentation within one approved group. A trailing 4 digit number is optional.
func (e *Engine) ValidateString(s, sep string) error {
	if e == nil {
		return errors.New("nil engine")
	}
	if e.err != nil {
		return e.err
	}
	if len(e.groupPhrases) == 0 {
		return errors.New("engine has no reviewed dictionary")
	}
	if sep == "" {
		sep = "-"
	}
	if !e.validSep(sep) {
		return fmt.Errorf("unsupported separator %q", sep)
	}
	tokens := strings.Split(s, sep)
	if len(tokens) > 1 {
		tail := tokens[len(tokens)-1]
		if n, err := strconv.Atoi(tail); err == nil && len(tail) == 4 && n >= 1000 && n <= 9999 {
			tokens = tokens[:len(tokens)-1]
		}
	}
	if len(tokens) == 0 || len(tokens) > e.words.Policy.MaxWords {
		return errors.New("invalid word count")
	}
	for _, token := range tokens {
		if token == "" {
			return errors.New("empty token")
		}
		for _, r := range token {
			if r < 'a' || r > 'z' {
				return errors.New("identifiers require lowercase ASCII words")
			}
		}
	}
	canonical := strings.Join(tokens, "-")
	if !e.passesBlocklist(canonical) {
		return errors.New("blocked vocabulary")
	}
	for group, phrases := range e.groupPhrases {
		reach := make([]bool, len(tokens)+1)
		reach[0] = true
		for i := range tokens {
			if !reach[i] {
				continue
			}
			for j := i + 1; j <= len(tokens) && j-i <= e.maxPhraseTokens[group]; j++ {
				if _, ok := phrases[strings.Join(tokens[i:j], "-")]; ok {
					reach[j] = true
				}
			}
		}
		if reach[len(tokens)] {
			return nil
		}
	}
	return errors.New("unknown words or incompatible mixing groups")
}

func ValidateString(s, sep string) error { return defaultEngine.ValidateString(s, sep) }

// generationPlan proves length feasibility before generation. Each requested
// category supplies a whole expression, followed by optional same-group filler.
type generationPlan struct {
	required []wordBuckets
	pool     wordBuckets
	reach    [][]bool
	totals   []int
}

func makePlan(required []wordBuckets, pool wordBuckets, min, max int) generationPlan {
	p := generationPlan{required: required, pool: pool, reach: make([][]bool, len(required)+1)}
	for i := range p.reach {
		p.reach[i] = make([]bool, max+1)
	}
	fill := p.reach[len(required)]
	fill[0] = true
	for n := 1; n <= max; n++ {
		for size := 1; size <= n; size++ {
			if len(pool[size]) > 0 && fill[n-size] {
				fill[n] = true
				break
			}
		}
	}
	for i := len(required) - 1; i >= 0; i-- {
		for n := 1; n <= max; n++ {
			for size := 1; size <= n; size++ {
				if len(required[i][size]) > 0 && p.reach[i+1][n-size] {
					p.reach[i][n] = true
					break
				}
			}
		}
	}
	for n := min; n <= max; n++ {
		if p.reach[0][n] {
			p.totals = append(p.totals, n)
		}
	}
	return p
}

func (e *Engine) plans(categories []string, min, max int) ([]generationPlan, error) {
	if len(categories) == 0 {
		var plans []generationPlan
		// Sort groups so seeded results never depend on map iteration order.
		for _, g := range []string{"divine", "prophets", "names", "islamic", "general"} {
			p := makePlan(nil, e.groupBuckets[g], min, max)
			if len(p.totals) > 0 {
				plans = append(plans, p)
			}
		}
		if len(plans) == 0 {
			return nil, errors.New("word count is infeasible")
		}
		return plans, nil
	}
	group := ""
	seen := make(map[string]bool)
	var required []wordBuckets
	pool := make(wordBuckets, e.words.Policy.MaxWords+1)
	for _, raw := range categories {
		cat := strings.TrimSpace(raw)
		if reason, disabled := e.words.Policy.DisabledCategories[cat]; disabled {
			return nil, fmt.Errorf("disabled category %q: %s", cat, reason)
		}
		g, exists := e.words.Policy.CategoryGroups[cat]
		if !exists {
			return nil, fmt.Errorf("unknown category: %q", cat)
		}
		if group != "" && group != g {
			return nil, fmt.Errorf("incompatible categories: %s belongs to %s, expected %s", cat, g, group)
		}
		group = g
		if seen[cat] {
			continue
		}
		seen[cat] = true
		required = append(required, e.catBuckets[cat])
		for size, ws := range e.catBuckets[cat] {
			pool[size] = append(pool[size], ws...)
		}
	}
	if len(required) > max {
		return nil, errors.New("too many categories for maximum word count")
	}
	p := makePlan(required, pool, min, max)
	if len(p.totals) == 0 {
		return nil, errors.New("word count is infeasible for selected categories")
	}
	return []generationPlan{p}, nil
}

func chooseWord(src *rand.Rand, buckets wordBuckets, remaining int, next []bool) wordToken {
	count := 0
	for size := 1; size <= remaining; size++ {
		if next[remaining-size] {
			count += len(buckets[size])
		}
	}
	pick := src.Intn(count)
	for size := 1; size <= remaining; size++ {
		if !next[remaining-size] {
			continue
		}
		if pick < len(buckets[size]) {
			return buckets[size][pick]
		}
		pick -= len(buckets[size])
	}
	panic("unreachable: validated generation plan has no word")
}

func (p generationPlan) generate(src *rand.Rand) string {
	remaining := p.totals[src.Intn(len(p.totals))]
	parts := make([]string, 0, remaining)
	for i, buckets := range p.required {
		wt := chooseWord(src, buckets, remaining, p.reach[i+1])
		parts = append(parts, wt.text)
		remaining -= wt.tokens
	}
	for remaining > 0 {
		wt := chooseWord(src, p.pool, remaining, p.reach[len(p.required)])
		parts = append(parts, wt.text)
		remaining -= wt.tokens
	}
	return strings.Join(parts, "-")
}

// GenerateWithOptionsE returns complete results, or an explicit error.
func (e *Engine) GenerateWithOptionsE(opts Options) ([]string, error) {
	if e == nil {
		return nil, errors.New("nil engine")
	}
	if e.err != nil {
		return nil, e.err
	}
	if len(e.normWords) == 0 {
		return nil, errors.New("engine has no reviewed dictionary")
	}
	if opts.Repeat == 0 {
		opts.Repeat = 1
	}
	if opts.Sep == "" {
		opts.Sep = "-"
	}
	if opts.MinWords == 0 {
		opts.MinWords = 5
	}
	if opts.MaxWords == 0 {
		opts.MaxWords = 8
	}
	if opts.Repeat < 1 || opts.Repeat > e.words.Policy.MaxRepeat {
		return nil, fmt.Errorf("repeat must be between 1 and %d", e.words.Policy.MaxRepeat)
	}
	if !e.validSep(opts.Sep) {
		return nil, fmt.Errorf("unsupported separator %q; use -, _, . or /", opts.Sep)
	}
	if opts.MinWords < 1 || opts.MaxWords < opts.MinWords || opts.MaxWords > e.words.Policy.MaxWords {
		return nil, fmt.Errorf("word bounds must satisfy 1 <= min <= max <= %d", e.words.Policy.MaxWords)
	}
	plans := e.defaultPlans
	if len(opts.Categories) > 0 || opts.MinWords != 5 || opts.MaxWords != 8 {
		var err error
		plans, err = e.plans(opts.Categories, opts.MinWords, opts.MaxWords)
		if err != nil {
			return nil, err
		}
	}
	var src *rand.Rand
	if opts.Seed == 0 {
		src = randPool.Get().(*rand.Rand)
		defer randPool.Put(src)
	} else {
		src = rand.New(rand.NewSource(opts.Seed))
	}
	results := make([]string, opts.Repeat)
	for i := range results {
		for attempt := 0; attempt < 256; attempt++ {
			candidate := plans[src.Intn(len(plans))].generate(src)
			if opts.IncludeRandomNumber {
				var n int
				if opts.Seed != 0 {
					n = src.Intn(9000) + 1000
				} else {
					var b [2]byte
					for {
						if _, err := crand.Read(b[:]); err != nil {
							return nil, fmt.Errorf("random suffix: %w", err)
						}
						v := int(b[0])<<8 | int(b[1])
						if v < 63000 {
							n = v%9000 + 1000
							break
						}
					}
				}
				candidate += "-" + strconv.Itoa(n)
			}
			if !e.passesBlocklist(candidate) {
				continue
			}
			results[i] = strings.ReplaceAll(candidate, "-", opts.Sep)
			break
		}
		if results[i] == "" {
			return nil, fmt.Errorf("failed to generate safe string at index %d", i)
		}
	}
	return results, nil
}

// Options controls length, output count and selection of reviewed categories.
type Options struct {
	PrefixThreshold     float64 // Deprecated: unused upstream compatibility field.
	SuffixThreshold     float64 // Deprecated: unused upstream compatibility field.
	Repeat              int     // 0 means 1; otherwise 1..1,000,000.
	Sep                 string  // Empty means "-"; supported: "-", "_", ".", "/".
	Seed                int64   // Nonzero seeds reproduce outputs; outputs are not secret tokens.
	MinWords            int     // Separator-delimited tokens; 0 means 5. Numeric suffix excluded.
	MaxWords            int     // 0 means 8; must be >= MinWords and <= 64.
	IncludeRandomNumber bool
	Categories          []string // Every distinct selected category contributes an expression.
}

// GenerateWithOptions is the compatibility wrapper; invalid options return nil.
func (e *Engine) GenerateWithOptions(opts Options) []string {
	r, _ := e.GenerateWithOptionsE(opts)
	return r
}
func (e *Engine) Generate() string {
	r := e.GenerateWithOptions(Options{IncludeRandomNumber: true})
	if len(r) == 0 {
		return ""
	}
	return r[0]
}
func (e *Engine) GenerateN(n int) []string {
	return e.GenerateWithOptions(Options{Repeat: n, IncludeRandomNumber: true})
}
func (e *Engine) GetCategories() []string {
	if e == nil {
		return nil
	}
	keys := make([]string, 0, len(e.words.Categories))
	for k := range e.words.Categories {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func Generate() string                                 { return defaultEngine.Generate() }
func GenerateN(n int) []string                         { return defaultEngine.GenerateN(n) }
func GenerateWithOptionsE(o Options) ([]string, error) { return defaultEngine.GenerateWithOptionsE(o) }
func GenerateWithOptions(o Options) []string           { return defaultEngine.GenerateWithOptions(o) }
func GetCategories() []string                          { return defaultEngine.GetCategories() }
func isSacredCategory(cat string) bool {
	g := defaultEngine.words.Policy.CategoryGroups[cat]
	return g != "" && g != "general"
}
func isMundaneCategory(cat string) bool {
	return defaultEngine.words.Policy.CategoryGroups[cat] == "general"
}
