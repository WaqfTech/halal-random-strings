package halalrandomstrings

import (
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

func TestReproducibility(t *testing.T) {
	opts := Options{
		Repeat: 5,
		Seed:   42,
	}
	results1 := GenerateWithOptions(opts)
	results2 := GenerateWithOptions(opts)

	if len(results1) != len(results2) {
		t.Fatalf("expected %d results, got %d", len(results1), len(results2))
	}

	for i := range results1 {
		if results1[i] != results2[i] {
			t.Fatalf("results are not reproducible with the same seed. got %q and %q", results1[i], results2[i])
		}
	}
}

func TestWordlistsLoad(t *testing.T) {
	if len(prefixes) == 0 {
		t.Fatal("prefixes wordlist is empty")
	}
	if len(modifiers) == 0 {
		t.Fatal("modifiers wordlist is empty")
	}
	if len(nouns) == 0 {
		t.Fatal("nouns wordlist is empty")
	}
	if len(suffixes) == 0 {
		t.Fatal("suffixes wordlist is empty")
	}
	if len(blocked) == 0 {
		t.Fatal("blocked wordlist is empty")
	}
}
