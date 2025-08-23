package halalrandomstrings

import (
	"strconv"
	"strings"
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
		Repeat:   1,
		MinWords: 5, // Updated default
		MaxWords: 8, // Updated default
	}
	results := GenerateWithOptions(opts)
	for _, result := range results {
		parts := strings.Split(result, "-")
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