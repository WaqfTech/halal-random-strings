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
}

func TestWordCount(t *testing.T) {
	opts := Options{
		Repeat:              10, // Test multiple strings
		Sep:                 "-",
		MinWords:            5, 
		MaxWords:            8, 
		IncludeRandomNumber: true,
	}
	results := GenerateWithOptions(opts)
	for _, result := range results {
		parts := strings.Split(result, opts.Sep)
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
	// Test with IncludeRandomNumber = true (default)
	optsTrue := Options{
		Repeat:              1,
		Sep:                 "-",
		IncludeRandomNumber: true,
	}
	resultTrue := GenerateWithOptions(optsTrue)[0]
	partsTrue := strings.Split(resultTrue, optsTrue.Sep)
	if _, err := strconv.Atoi(partsTrue[len(partsTrue)-1]); err != nil {
		t.Fatalf("expected random number, but got error: %v in %q", err, resultTrue)
	}

	// Test with IncludeRandomNumber = false
	optsFalse := Options{
		Repeat:              1,
		Sep:                 "-",
		IncludeRandomNumber: false,
	}
	resultFalse := GenerateWithOptions(optsFalse)[0]
	partsFalse := strings.Split(resultFalse, optsFalse.Sep)
	if _, err := strconv.Atoi(partsFalse[len(partsFalse)-1]); err == nil {
		t.Fatalf("did not expect random number, but found one in %q", resultFalse)
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
		normalizedResult := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(result, optsSahaba.Sep, " "), "-", " "))
		for _, sahabi := range words.Categories["sahaba"] {
			// Normalize the sahabi name for comparison
			normalizedSahabi := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(sahabi, "'", " "), "-", " "))
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