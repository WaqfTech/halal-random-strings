package halalrandomstrings

import (
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestIsSafe(t *testing.T) {
	if isSafe("this-is-a-test-with-wine") {
		t.Fatal("isSafe returned true for a string containing a blocked word")
	}
	if !isSafe("this-is-a-safe-string") {
		t.Fatal("isSafe returned false for a safe string")
	}
}

func TestWordsLoad(t *testing.T) {
	if len(words.Categories) == 0 {
		t.Fatal("categories in words.json is empty")
	}
	if len(words.Rules) == 0 {
		t.Fatal("rules in words.json is empty")
	}
	if len(words.Blocked) == 0 {
		t.Fatal("blocked in words.json is empty")
	}
	categories := GetCategories()
	if len(categories) == 0 {
		t.Fatal("GetCategories() returned no categories")
	}
}

func TestWordCount(t *testing.T) {
	opts := Options{
		Repeat:              10, // Test multiple strings
		MinWords:            5,
		MaxWords:            8,
		IncludeRandomNumber: true,
	}
	results := GenerateWithOptions(opts)
	if len(results) != opts.Repeat {
		t.Fatalf("expected %d results, got %d", opts.Repeat, len(results))
	}
	for _, result := range results {
		parts := strings.Split(result, "-")
		if len(parts) <= 1 {
			t.Fatalf("generated string has too few parts: %q", result)
		}
		// The last part is the random number, so we subtract 1 from the total parts
		wordCount := len(parts) - 1

		if wordCount < opts.MinWords || wordCount > opts.MaxWords {
			t.Fatalf("generated string has %d words, but expected between %d and %d words: %q", wordCount, opts.MinWords, opts.MaxWords, result)
		}

		// Check if the last part is a number
		_, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			t.Fatalf("last part of the generated string is not a number: %q", result)
		}
	}
}

func TestIncludeRandomNumber(t *testing.T) {
	// Test with IncludeRandomNumber = true and uninitialized Sep (verifies default "-" and avoids out-of-range panic)
	optsTrue := Options{
		Repeat:              1,
		IncludeRandomNumber: true,
	}
	resultsTrue := GenerateWithOptions(optsTrue)
	if len(resultsTrue) == 0 || resultsTrue[0] == "" {
		t.Fatal("expected non-empty string for IncludeRandomNumber=true")
	}
	partsTrue := strings.Split(resultsTrue[0], "-")
	if len(partsTrue) == 0 {
		t.Fatalf("parts slice is empty for %q", resultsTrue[0])
	}
	if _, err := strconv.Atoi(partsTrue[len(partsTrue)-1]); err != nil {
		t.Fatalf("expected random number suffix, but got error: %v in %q", err, resultsTrue[0])
	}

	// Test with IncludeRandomNumber = false and uninitialized Sep
	optsFalse := Options{
		Repeat:              1,
		IncludeRandomNumber: false,
	}
	resultsFalse := GenerateWithOptions(optsFalse)
	if len(resultsFalse) == 0 || resultsFalse[0] == "" {
		t.Fatal("expected non-empty string for IncludeRandomNumber=false")
	}
	partsFalse := strings.Split(resultsFalse[0], "-")
	if len(partsFalse) == 0 {
		t.Fatalf("parts slice is empty for %q", resultsFalse[0])
	}
	if _, err := strconv.Atoi(partsFalse[len(partsFalse)-1]); err == nil {
		t.Fatalf("did not expect random number suffix, but found one in %q", resultsFalse[0])
	}
}

func TestCategories(t *testing.T) {
	// Test with a specific category (sahaba)
	optsSahaba := Options{
		Repeat:              10,
		Categories:          []string{"sahaba"},
		Sep:                 "-",
		MinWords:            1,
		MaxWords:            5, // Sahaba names can be multi-word
		IncludeRandomNumber: false,
	}
	resultsSahaba, err := GenerateWithOptionsE(optsSahaba)
	if err != nil || len(resultsSahaba) != optsSahaba.Repeat {
		t.Fatalf("sahaba generation: count=%d, error=%v", len(resultsSahaba), err)
	}
	for _, result := range resultsSahaba {
		if !fitsCategory(result, "sahaba") {
			t.Fatalf("not composed of whole sahaba names: %q", result)
		}
	}

	// Test with multiple categories (colors_arabic, nouns_places)
	optsMulti := Options{
		Repeat:              10,
		Categories:          []string{"colors_arabic", "nouns_places"},
		Sep:                 "-",
		MinWords:            2,
		MaxWords:            2,
		IncludeRandomNumber: false,
	}
	resultsMulti, err := GenerateWithOptionsE(optsMulti)
	if err != nil || len(resultsMulti) != optsMulti.Repeat {
		t.Fatalf("multi-category generation: %d, %v", len(resultsMulti), err)
	}
	for _, result := range resultsMulti {
		parts := strings.Split(result, optsMulti.Sep)
		if len(parts) != 2 {
			t.Fatalf("expected 2 words for multi-category test, got %d in %q", len(parts), result)
		}

		// Check if first word is an adjective and second is a place
		adjFound := false
		for _, adj := range words.Categories["colors_arabic"] {
			if normalizeWord(adj, "-") == parts[0] {
				adjFound = true
				break
			}
		}
		if !adjFound {
			t.Fatalf("first word %q not found in colors_arabic category for %q", parts[0], result)
		}

		placeFound := false
		for _, place := range words.Categories["nouns_places"] {
			if normalizeWord(place, "-") == parts[1] {
				placeFound = true
				break
			}
		}
		if !placeFound {
			t.Fatalf("second word %q not found in nouns_places category for %q", parts[1], result)
		}
	}
}

func TestEngineConcurrent(t *testing.T) {
	engine, err := NewDefaultEngine()
	if err != nil {
		t.Fatalf("failed to create default engine: %v", err)
	}

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = engine.Generate()
				_ = engine.GetCategories()
				_ = engine.IsSafe("safe-string")
			}
		}()
	}
	wg.Wait()
}

func TestScunthorpeAndZeroFalsePositives(t *testing.T) {
	venerableNames := []string{
		"Muhammad",
		"Ahmad",
		"Aisha",
		"Khadijah",
		"Uthman ibn Affan",
		"Abu Hurairah",
		"Compassion",
		"Steadfastness",
		"sunnah",
		"obedient",
	}

	for _, name := range venerableNames {
		if !isSafe(name) {
			t.Errorf("expected %q to be considered safe, but was blocked (false positive)", name)
		}
	}

	blockedInputs := []string{
		"wine",
		"this-is-a-test-with-wine",
		"show-off",
		"user-with-pork-dish",
	}

	for _, input := range blockedInputs {
		if isSafe(input) {
			t.Errorf("expected %q to be blocked, but was considered safe", input)
		}
	}
}

func TestCustomSeparatorAndBoundaries(t *testing.T) {
	// Test custom separator "_"
	optsSep := Options{
		Repeat:              10,
		Sep:                 "_",
		MinWords:            3,
		MaxWords:            5,
		IncludeRandomNumber: true,
	}
	resultsSep := GenerateWithOptions(optsSep)
	for _, res := range resultsSep {
		if strings.Contains(res, "-") {
			t.Errorf("expected pure underscore delimiters for sep='_', got %q", res)
		}
		if strings.Contains(res, " ") {
			t.Errorf("expected no spaces in generated string, got %q", res)
		}
		parts := strings.Split(res, "_")
		wordCount := len(parts) - 1
		if wordCount < 3 || wordCount > 5 {
			t.Errorf("expected between 3 and 5 words, got %d in %q", wordCount, res)
		}

		num, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			t.Fatalf("failed to parse numeric suffix in %q: %v", res, err)
		}
		if num < 1000 || num > 9999 {
			t.Errorf("expected 4-digit numeric suffix in [1000, 9999], got %d in %q", num, res)
		}
	}

	// Test boundary constraints MinWords == MaxWords == 4
	optsExact := Options{
		Repeat:              10,
		Sep:                 "-",
		MinWords:            4,
		MaxWords:            4,
		IncludeRandomNumber: false,
	}
	resultsExact := GenerateWithOptions(optsExact)
	for _, res := range resultsExact {
		parts := strings.Split(res, "-")
		if len(parts) != 4 {
			t.Errorf("expected exactly 4 words, got %d in %q", len(parts), res)
		}
	}
}

func TestCategoryFilteringAndFallback(t *testing.T) {
	// 1. Single category without compound rules (e.g. fruits) should not fail or hang
	optsSingle := Options{
		Repeat:              5,
		Categories:          []string{"fruits"},
		MinWords:            2,
		MaxWords:            3,
		IncludeRandomNumber: false,
	}
	results, err := GenerateWithOptionsE(optsSingle)
	if err != nil {
		t.Fatalf("unexpected error for single category 'fruits': %v", err)
	}
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
	for _, res := range results {
		if res == "" {
			t.Fatal("expected non-empty string for 'fruits' fallback")
		}
	}

	// 2. Unknown category should return an explicit error
	optsInvalid := Options{
		Repeat:     1,
		Categories: []string{"invalid_nonexistent_category_xyz"},
	}
	_, errInvalid := GenerateWithOptionsE(optsInvalid)
	if errInvalid == nil {
		t.Fatal("expected error for nonexistent category, got nil")
	}
	if !strings.Contains(errInvalid.Error(), "unknown category") {
		t.Errorf("expected 'unknown category' error message, got %v", errInvalid)
	}
}

func TestHolySanctuariesAndAnimals(t *testing.T) {
	// 1. Sanctuaries category exists and contains sacred sites
	sanctuaries, exists := words.Categories["holy_sanctuaries"]
	if !exists || len(sanctuaries) == 0 {
		t.Fatal("expected holy_sanctuaries category to exist and have entries")
	}

	sanctuarySet := make(map[string]bool)
	for _, s := range sanctuaries {
		sanctuarySet[strings.ToLower(s)] = true
	}
	for _, expected := range []string{"kaaba", "masjid", "al-aqsa", "mecca", "medina"} {
		if !sanctuarySet[expected] {
			t.Errorf("expected %q in holy_sanctuaries", expected)
		}
	}

	// 2. Ensure holy sanctuaries are not in nouns_places
	for _, place := range words.Categories["nouns_places"] {
		norm := strings.ToLower(place)
		if norm == "kaaba" || norm == "masjid" || norm == "al-aqsa" {
			t.Errorf("sacred site %q should not be in nouns_places", place)
		}
	}

	// 3. Ensure impure animals are purged from animals category
	for _, animal := range words.Categories["animals"] {
		norm := strings.ToLower(animal)
		if norm == "warthog" || norm == "boar" || norm == "wild-dog" {
			t.Errorf("impure animal %q should not be in animals", animal)
		}
	}

	// 4. Test animal + geographic features generation
	opts := Options{
		Repeat:              5,
		Categories:          []string{"animals", "geographic_features"},
		MinWords:            2,
		MaxWords:            2,
		IncludeRandomNumber: false,
	}
	results, err := GenerateWithOptionsE(opts)
	if err != nil {
		t.Fatalf("failed to generate with animals and geographic_features: %v", err)
	}
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
}

func TestAsmaAllahSegregation(t *testing.T) {
	// 1. asma_allah category exists
	asma, exists := words.Categories["asma_allah"]
	if !exists || len(asma) == 0 {
		t.Fatal("expected asma_allah category to exist and have entries")
	}

	asmaSet := make(map[string]bool)
	for _, a := range asma {
		asmaSet[strings.ToLower(a)] = true
	}
	for _, expected := range []string{"rahman", "khaliq", "quddus", "razzaq", "samad"} {
		if !asmaSet[expected] {
			t.Errorf("expected %q in asma_allah", expected)
		}
	}

	// 2. Ensure exclusive Asma' Allah are not in generic adjectives
	for _, adj := range words.Categories["adjectives"] {
		norm := strings.ToLower(adj)
		if norm == "rahman" || norm == "khaliq" || norm == "qahhar" || norm == "razzaq" || norm == "quddus" {
			t.Errorf("divine attribute %q should not be in generic adjectives", adj)
		}
	}

	// The approved policy permits only self-mixing, without standalone prefixes.
	results, err := GenerateWithOptionsE(Options{Repeat: 5, Categories: []string{"asma_allah"}, MinWords: 2, MaxWords: 2})
	if err != nil || len(results) != 5 {
		t.Fatalf("Asma-only generation: %d, %v", len(results), err)
	}
	for _, result := range results {
		if !fitsCategory(result, "asma_allah") {
			t.Fatalf("non-Asma term: %s", result)
		}
	}
	if _, err := GenerateWithOptionsE(Options{Categories: []string{"servant_prefixes", "asma_allah"}}); err == nil {
		t.Fatal("standalone prefixes must be rejected")
	}

}

func BenchmarkGenerate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Generate()
	}
}

func BenchmarkGenerateWithOptions(b *testing.B) {
	opts := Options{
		Repeat:              1,
		MinWords:            5,
		MaxWords:            8,
		IncludeRandomNumber: true,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenerateWithOptions(opts)
	}
}

func TestTheologicalQuarantineAndSensitivityGuarantees(t *testing.T) {
	// 1. Verify that culturally insulting, demonic, or derogatory words are rejected by isSafe
	prohibitedWords := []string{
		"donkey", "mule", "ape", "monkey", "baboon", "chimpanzee", "gorilla",
		"hyena", "jackal", "vulture", "rat", "mouse", "toad", "snake", "cobra", "viper",
		"shoe", "sandal", "boot", "slipper", "shoelace", "dustpan", "mop", "broom", "vacuum-cleaner",
		"iblis", "jahannam", "shirk", "kufr", "nifaq", "bid-ah", "awrah",
		"death-of-the-prophet", "su-al-zann", "dine",
		"ahmad-donkey-food", "ismail-dine-desert", "aisha-shoe", "kaaba-baboon",
	}
	for _, pw := range prohibitedWords {
		if isSafe(pw) {
			t.Errorf("critical sensitivity failure: prohibited term %q must be rejected by isSafe", pw)
		}
	}

	// 2. Verify that incompatible category requests return explicit validation error
	incompatiblePairs := [][]string{
		{"muslim_names_male", "animals"},
		{"sahaba", "arabic_food"},
		{"asma_allah", "vegetables"},
		{"holy_sanctuaries", "fruits"},
		{"islamic_virtues", "spices"},
	}
	for _, pair := range incompatiblePairs {
		opts := Options{
			Repeat:     1,
			Categories: pair,
		}
		_, err := GenerateWithOptionsE(opts)
		if err == nil {
			t.Errorf("expected error for incompatible categories %v, got nil", pair)
		} else if !strings.Contains(err.Error(), "incompatible categories") {
			t.Errorf("expected 'incompatible categories' error for %v, got: %v", pair, err)
		}
	}

	// 3. Verify zero cross-domain violations across 1000 generated strings
	opts := Options{
		Repeat:              5000,
		MinWords:            5,
		MaxWords:            8,
		IncludeRandomNumber: true,
	}
	results, err := GenerateWithOptionsE(opts)
	if err != nil {
		t.Fatalf("batch generation failed: %v", err)
	}

	if len(results) != opts.Repeat {
		t.Fatalf("batch length=%d", len(results))
	}
	for _, result := range results {
		if err := ValidateString(result, "-"); err != nil {
			t.Fatalf("policy failure %q: %v", result, err)
		}
	}
}
