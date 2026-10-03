package halalrandomstrings

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

// This oracle checks membership in one requested category, independently of
// the production generation plan and group validator.
func fitsCategory(value, category string) bool {
	phrases := make(map[string]bool)
	for _, word := range words.Categories[category] {
		phrases[normalizeWord(word, "-")] = true
	}
	tokens := strings.Split(value, "-")
	reach := make([]bool, len(tokens)+1)
	reach[0] = true
	for i := range tokens {
		if !reach[i] {
			continue
		}
		for j := i + 1; j <= len(tokens); j++ {
			if phrases[strings.Join(tokens[i:j], "-")] {
				reach[j] = true
			}
		}
	}
	return reach[len(tokens)]
}

// Accept only a segmentation that uses both requested categories and no others.
// An expression shared by both categories cannot satisfy both at once.
func fitsRequestedPair(value, left, right string) bool {
	leftSet, rightSet := make(map[string]bool), make(map[string]bool)
	for _, word := range words.Categories[left] {
		leftSet[normalizeWord(word, "-")] = true
	}
	for _, word := range words.Categories[right] {
		rightSet[normalizeWord(word, "-")] = true
	}
	tokens := strings.Split(value, "-")
	reach := make([][4]bool, len(tokens)+1)
	reach[0][0] = true
	for i := range tokens {
		for j := i + 1; j <= len(tokens); j++ {
			phrase := strings.Join(tokens[i:j], "-")
			for mask := 0; mask < 4; mask++ {
				if !reach[i][mask] {
					continue
				}
				if leftSet[phrase] {
					reach[j][mask|1] = true
				}
				if rightSet[phrase] {
					reach[j][mask|2] = true
				}
			}
		}
	}
	if left == right {
		return reach[len(tokens)][1] || reach[len(tokens)][2]
	}
	return reach[len(tokens)][3]
}

func TestApprovedCategoryMatrix(t *testing.T) {
	// Intentionally explicit: policy metadata cannot make this test silently
	// accept a protected category moved into the general group.
	expected := map[string][]string{
		"divine":   {"asma_allah"},
		"prophets": {"prophets"},
		"names":    {"muslim_names_male", "muslim_names_female"},
		"islamic":  {"adjectives", "sahaba", "islamic_art_forms", "islamic_events", "islamic_golden_age_scholars", "islamic_inventions", "islamic_months", "islamic_virtues", "scholarly_terms", "adab_terms", "holy_sanctuaries", "muslim_empires", "nouns_concepts"},
		"general":  {"animals", "arabic_food", "architectural_elements", "birds", "colors_arabic", "days_of_week_arabic", "flowers", "fruits", "gems_minerals", "geographic_features", "jordanian_food", "nouns_objects", "nouns_places", "saudi_food", "spices", "suffixes", "trees", "vegetables", "yemeni_food"},
	}
	groups := map[string]string{}
	for group, cats := range expected {
		for _, cat := range cats {
			groups[cat] = group
			if words.Policy.CategoryGroups[cat] != group {
				t.Fatalf("%s must belong to %s", cat, group)
			}
		}
	}
	if len(groups) != len(GetCategories()) {
		t.Fatal("unreviewed category added")
	}
	for _, left := range GetCategories() {
		for _, right := range GetCategories() {
			t.Run(left+"/"+right, func(t *testing.T) {
				result, err := GenerateWithOptionsE(Options{Repeat: 3, Seed: 42, Categories: []string{left, right}, MinWords: 2, MaxWords: 16})
				if groups[left] != groups[right] {
					if err == nil || result != nil {
						t.Fatalf("forbidden pair accepted: %v", result)
					}
					return
				}
				if err != nil || len(result) != 3 {
					t.Fatalf("approved pair: count=%d, err=%v", len(result), err)
				}
				for _, value := range result {
					if !fitsRequestedPair(value, left, right) {
						t.Fatalf("category omitted or unselected category used: %q", value)
					}
					if err := ValidateString(value, "-"); err != nil {
						t.Fatalf("%s: %v", value, err)
					}
				}
			})
		}
	}
}

func TestCustomDictionariesDisabledAndUnaliased(t *testing.T) {
	w := Words{Categories: map[string][]string{"prophets": {"ahmad"}, "animals": {"donkey"}, "arabic_food": {"food"}}, Rules: []Rule{{Pattern: []string{"prophets", "animals", "arabic_food"}}}}
	e := NewEngine(w)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			w.Rules[0].Pattern[0] = "animals"
			w.Categories["prophets"][0] = "donkey"
		}
	}()
	for i := 0; i < 1000; i++ {
		if result, err := e.GenerateWithOptionsE(Options{}); !errors.Is(err, ErrCustomDictionariesDisabled) || result != nil {
			t.Fatalf("custom dictionary accepted: %v, %v", result, err)
		}
	}
	wg.Wait()
	for _, load := range []func(string) error{LoadWords, e.LoadWords, DefaultEngine().LoadWords} {
		if err := load("/this-file-must-not-be-read"); !errors.Is(err, ErrCustomDictionariesDisabled) {
			t.Fatal(err)
		}
	}
	if e.Generate() != "" || e.IsSafe("safe") || len(e.GetCategories()) != 0 {
		t.Fatal("disabled engine must fail closed")
	}
	if DefaultEngine().Generate() == "" {
		t.Fatal("rejected override changed default engine")
	}
}

func TestIdentifierValidationRegressions(t *testing.T) {
	cases := []struct {
		value string
		valid bool
	}{
		{"ahmad-donkey-food", false}, {"ahmad-cow-rice", false}, {"ahmad-olive", false},
		{"taha-ramadan", false}, {"muzammil-ramadan", false}, {"kaffir-lime-leaves", false}, {"aisha-ramadan", false}, {"rahman-ramadan", false}, {"muhammad-aisha", false},
		{"abd-rahman", false}, {"mihrab-cow", false}, {"shrine-holy-silver", false},
		{"sad-ibn-abi-waqqas-cow", false}, {"hamza-ibn-abd-al-muttalib-rice", false},
		{"shrine-holy-dwelling-green-peppercorn-falafel", false},
		{"ramadan-olive", false}, {"unknownword", false}, {"olive--apple", false},
		{"olive-123", false}, {"olive-0123", false}, {"olive-12345", false},
		{"olive\napple", false}, {" olive", false}, {"Olive", false}, {"", false},
		{"makrut-lime-leaves-olive", true}, {"olive-apple", true}, {"cow-olive", true}, {"aisha-maher", true},
		{"abdul-hadi-maher", true}, {"rahman-rahim", true}, {"ahmad-musa", true},
		{"mihrab-ramadan", true}, {"sad-ibn-abi-waqqas-ramadan", true},
		{"bird-of-paradise-olive", true}, {"olive-1234", true},
	}
	for _, tc := range cases {
		t.Run(tc.value, func(t *testing.T) {
			for _, sep := range []string{"-", "_", ".", "/"} {
				value := strings.ReplaceAll(tc.value, "-", sep)
				err := ValidateString(value, sep)
				if (err == nil) != tc.valid {
					t.Fatalf("ValidateString(%q,%q)=%v; valid=%v", value, sep, err, tc.valid)
				}
			}
		})
	}
	for _, value := range []string{"s'mores", "smores", "wine", "prefix-phrase-of-five-blocked-words-suffix"} {
		e := &Engine{}
		if err := e.buildBlocked([]string{"smores", "wine", "phrase-of-five-blocked-words"}); err != nil {
			t.Fatal(err)
		}
		if e.passesBlocklist(value) {
			t.Fatalf("normalization or long phrase bypass: %s", value)
		}
	}
}

func TestRequestedCategoriesAndOptions(t *testing.T) {
	for _, sep := range []string{"donkey", "wine", "x", "\n", " ", "--", "🐴"} {
		if result, err := GenerateWithOptionsE(Options{Sep: sep, Categories: []string{"prophets"}}); err == nil || result != nil {
			t.Fatalf("injected separator accepted: %q", sep)
		}
	}
	for _, opts := range []Options{{Repeat: -1}, {Repeat: 1000001}, {MinWords: -1}, {MinWords: 10, MaxWords: 2}, {MaxWords: 65}, {MinWords: 1, MaxWords: 1, Categories: []string{"sahaba"}}, {MinWords: 1, MaxWords: 1, Categories: []string{"colors_arabic", "fruits"}}, {Categories: []string{"servant_prefixes"}}, {Categories: []string{"fruits", ""}}} {
		if result, err := GenerateWithOptionsE(opts); err == nil || result != nil {
			t.Fatalf("invalid options accepted: %+v", opts)
		}
	}
	result, err := GenerateWithOptionsE(Options{Repeat: 100, Seed: 42, Categories: []string{" colors_arabic ", "fruits", "fruits"}, MinWords: 2, MaxWords: 2})
	if err != nil || len(result) != 100 {
		t.Fatalf("requested categories: %d, %v", len(result), err)
	}
	for _, value := range result {
		tokens := strings.Split(value, "-")
		if len(tokens) != 2 || !fitsCategory(tokens[0], "colors_arabic") || !fitsCategory(tokens[1], "fruits") {
			t.Fatalf("category silently omitted: %s", value)
		}
	}
}

func TestEmbeddedPolicyFailsClosed(t *testing.T) {
	cases := map[string]func(*Words){
		"missing policy":  func(w *Words) { w.Policy = MixingPolicy{} },
		"unclassified":    func(w *Words) { delete(w.Policy.CategoryGroups, "animals") },
		"protected group": func(w *Words) { w.Policy.CategoryGroups["asma_allah"] = "general" },
		"unknown group":   func(w *Words) { w.Policy.CategoryGroups["animals"] = "neutral" },
		"ambiguous":       func(w *Words) { w.Categories["animals"] = append(w.Categories["animals"], "Ahmad") },
		"empty":           func(w *Words) { w.Categories["animals"] = nil },
		"blocked":         func(w *Words) { w.Categories["animals"] = append(w.Categories["animals"], "donkey") },
		"unicode":         func(w *Words) { w.Categories["animals"] = append(w.Categories["animals"], "hammām") },
		"mixed rule": func(w *Words) {
			w.Rules = append(w.Rules, Rule{Pattern: []string{"prophets", "animals"}, Template: "{prophets}-{animals}"})
		},
		"bad template": func(w *Words) { w.Rules[0].Template = "injected-{adjectives}" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			var w Words
			if err := json.Unmarshal(defaultWordsData, &w); err != nil {
				t.Fatal(err)
			}
			mutate(&w)
			e := &Engine{words: w}
			if err := e.buildIndexes(); err == nil {
				t.Fatal("invalid embedded policy accepted")
			}
		})
	}
}

func TestDefaultCorpusPolicy(t *testing.T) {
	for _, sep := range []string{"-", "_", ".", "/"} {
		t.Run(fmt.Sprintf("separator_%s", sep), func(t *testing.T) {
			result, err := GenerateWithOptionsE(Options{Repeat: 10000, Seed: 424242, Sep: sep, IncludeRandomNumber: true})
			if err != nil || len(result) != 10000 {
				t.Fatalf("corpus: %d, %v", len(result), err)
			}
			for _, value := range result {
				if err := ValidateString(value, sep); err != nil {
					t.Fatalf("corpus violation %q: %v", value, err)
				}
			}
		})
	}
}
