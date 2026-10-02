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
	resultsSahaba := GenerateWithOptions(optsSahaba)
	for _, result := range resultsSahaba {
		found := false
		// Normalize the generated result for comparison
		normalizedResult := strings.Join(strings.Fields(strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(result, optsSahaba.Sep, " "), "-", " "))), " ")
		for _, sahabi := range words.Categories["sahaba"] {
			// Normalize the sahabi name for comparison
			normalizedSahabi := strings.Join(strings.Fields(strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(sahabi, "'", " "), "-", " "))), " ")
			if normalizedSahabi == normalizedResult {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("generated string %q not found in sahaba category (normalized: %q)", result, normalizedResult)
		}
	}

	// Test with multiple categories (adjectives, nouns_places)
	optsMulti := Options{
		Repeat:              10,
		Categories:          []string{"adjectives", "nouns_places"},
		Sep:                 "-",
		MinWords:            2, 
		MaxWords:            2,
		IncludeRandomNumber: false,
	}
	resultsMulti := GenerateWithOptions(optsMulti)
	for _, result := range resultsMulti {
		parts := strings.Split(result, optsMulti.Sep)
		if len(parts) != 2 {
			t.Fatalf("expected 2 words for multi-category test, got %d in %q", len(parts), result)
		}
		
		// Check if first word is an adjective and second is a place
		adjFound := false
		for _, adj := range words.Categories["adjectives"] {
			if strings.ToLower(adj) == parts[0] {
				adjFound = true
				break
			}
		}
		if !adjFound {
			t.Fatalf("first word %q not found in adjectives category for %q", parts[0], result)
		}

		placeFound := false
		for _, place := range words.Categories["nouns_places"] {
			if strings.ToLower(place) == parts[1] {
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