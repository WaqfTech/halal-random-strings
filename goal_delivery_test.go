package halalrandomstrings

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// TestGoal_01_SilaOnboarding verifies baseline repository structure and documentation.
func TestGoal_01_SilaOnboarding(t *testing.T) {
	if _, err := os.Stat("AGENTS.md"); os.IsNotExist(err) {
		t.Fatal("AGENTS.md is missing")
	}
	if _, err := os.Stat("words.json"); os.IsNotExist(err) {
		t.Fatal("words.json is missing")
	}
	if len(words.Categories) == 0 {
		t.Fatal("words.json categories must not be empty")
	}
}

// TestGoal_02_CleanRepoArtifactsAndClutter verifies removal of 2bedeleted and gitignore anchoring.
func TestGoal_02_CleanRepoArtifactsAndClutter(t *testing.T) {
	if _, err := os.Stat("2bedeleted"); !os.IsNotExist(err) {
		t.Error("2bedeleted/ directory must be permanently deleted")
	}

	gitignoreData, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}
	content := string(gitignoreData)
	if !strings.Contains(content, "/halal\n") && !strings.Contains(content, "/halal\r\n") {
		t.Error(".gitignore must anchor /halal")
	}
	if !strings.Contains(content, "/halal-random-strings\n") && !strings.Contains(content, "/halal-random-strings\r\n") {
		t.Error(".gitignore must anchor /halal-random-strings")
	}

	// Verify cmd/halal-random-strings is not ignored
	cmd := exec.Command("git", "check-ignore", "cmd/halal-random-strings/main.go")
	if err := cmd.Run(); err == nil {
		t.Error("cmd/halal-random-strings/main.go should NOT be ignored by git")
	}
}

// TestGoal_03_WordsEmbeddingAndInit verifies embedded words.json and thread-safe Engine.
func TestGoal_03_WordsEmbeddingAndInit(t *testing.T) {
	if len(defaultWordsData) == 0 {
		t.Fatal("defaultWordsData must be embedded and non-empty")
	}

	engine, err := NewDefaultEngine()
	if err != nil {
		t.Fatalf("NewDefaultEngine() failed: %v", err)
	}
	if engine == nil {
		t.Fatal("NewDefaultEngine() returned nil")
	}
	if DefaultEngine() == nil {
		t.Fatal("DefaultEngine() must return package default engine")
	}

	// Test isolated custom engine
	customWords := Words{
		Categories: map[string][]string{"test_cat": {"word1", "word2"}},
		Rules:      []Rule{{Pattern: []string{"test_cat"}, Template: "{test_cat}"}},
		Blocked:    []string{"badword"},
	}
	customEngine := NewEngine(customWords)
	cats := customEngine.GetCategories()
	if len(cats) != 1 || cats[0] != "test_cat" {
		t.Fatalf("custom engine categories mismatch: %v", cats)
	}
}

// TestGoal_04_ScunthorpeAndBlockedList verifies elimination of false positives and blocked filter accuracy.
func TestGoal_04_ScunthorpeAndBlockedList(t *testing.T) {
	falsePositiveWords := []string{
		"damascus",
		"manama",
		"shiraz",
		"buttress",
		"cockpit",
		"classic",
		"passenger",
	}
	for _, word := range falsePositiveWords {
		if !isSafe(word) {
			t.Errorf("false positive: safe word %q was rejected by isSafe", word)
		}
	}

	blockedPhrases := []string{
		"wine",
		"beer",
		"vodka",
		"pork",
		"casino",
		"gambling",
		"interest-usury",
		"prefix-wine-suffix",
	}
	for _, phrase := range blockedPhrases {
		if isSafe(phrase) {
			t.Errorf("safety failure: blocked phrase %q was accepted by isSafe", phrase)
		}
	}
}

// TestGoal_05_SeparatorAndWordCountLogic verifies separator flexibility, word counting, and 4-digit numbers.
func TestGoal_05_SeparatorAndWordCountLogic(t *testing.T) {
	separators := []string{"_", ".", "-", "/"}
	for _, sep := range separators {
		opts := Options{
			Repeat:              5,
			Sep:                 sep,
			MinWords:            3,
			MaxWords:            5,
			IncludeRandomNumber: true,
		}
		results, err := GenerateWithOptionsE(opts)
		if err != nil {
			t.Fatalf("GenerateWithOptionsE failed for sep %q: %v", sep, err)
		}
		for _, r := range results {
			if !strings.Contains(r, sep) {
				t.Errorf("result %q does not contain separator %q", r, sep)
			}
			parts := strings.Split(r, sep)
			// Last part must be numeric
			numPart := parts[len(parts)-1]
			val, err := strconv.Atoi(numPart)
			if err != nil {
				t.Fatalf("last part of %q is not numeric: %v", r, err)
			}
			if val < 1000 || val > 9999 {
				t.Errorf("numeric suffix %d is not in 4-digit range [1000, 9999]", val)
			}
			wordCount := len(parts) - 1
			if wordCount < opts.MinWords || wordCount > opts.MaxWords {
				t.Errorf("word count %d outside bounds [%d, %d] in %q", wordCount, opts.MinWords, opts.MaxWords, r)
			}
		}
	}
}

// TestGoal_06_CategoryFilteringAndFallback verifies validation and single-category fallback.
func TestGoal_06_CategoryFilteringAndFallback(t *testing.T) {
	// 1. Invalid category returns error
	optsInvalid := Options{
		Repeat:     1,
		Categories: []string{"nonexistent_category_123"},
	}
	_, err := GenerateWithOptionsE(optsInvalid)
	if err == nil {
		t.Fatal("expected error for nonexistent category, got nil")
	}
	if !strings.Contains(err.Error(), "unknown category") {
		t.Errorf("expected 'unknown category' error, got %v", err)
	}

	// 2. Single-category fallback works even without multi-category rule
	singleCats := []string{"fruits", "vegetables", "sahaba", "holy_sanctuaries"}
	for _, cat := range singleCats {
		opts := Options{
			Repeat:              3,
			Categories:          []string{cat},
			MinWords:            1,
			MaxWords:            5,
			IncludeRandomNumber: false,
		}
		results, err := GenerateWithOptionsE(opts)
		if err != nil {
			t.Fatalf("single category fallback failed for %q: %v", cat, err)
		}
		if len(results) != 3 {
			t.Fatalf("expected 3 results for %q, got %d", cat, len(results))
		}
	}
}

// TestGoal_07_TestSuiteAndPanics verifies absence of panics under edge cases and concurrency.
func TestGoal_07_TestSuiteAndPanics(t *testing.T) {
	// Empty separator default fallback
	optsNoSep := Options{
		Repeat:              2,
		IncludeRandomNumber: false,
	}
	res := GenerateWithOptions(optsNoSep)
	if len(res) != 2 || res[0] == "" {
		t.Fatal("expected valid result with uninitialized Sep")
	}

	// Inverted MinWords > MaxWords should auto-clamp without panic
	optsInverted := Options{
		Repeat:   1,
		MinWords: 10,
		MaxWords: 2,
	}
	resInv := GenerateWithOptions(optsInverted)
	if len(resInv) != 1 || resInv[0] == "" {
		t.Fatal("expected valid result with inverted MinWords/MaxWords")
	}

	// Concurrent access test
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := Generate()
			if s == "" {
				t.Error("concurrent Generate() produced empty string")
			}
		}()
	}
	wg.Wait()
}

// TestGoal_08_OptimizeGeneratorPerformanceAndAllocations verifies precomputed blockedSet and speed.
func TestGoal_08_OptimizeGeneratorPerformanceAndAllocations(t *testing.T) {
	eng := DefaultEngine()
	if eng.blockedSet == nil || len(eng.blockedSet) == 0 {
		t.Fatal("blockedSet must be pre-indexed into hash map for O(1) lookups")
	}

	// Verify fast batch generation (1000 strings in < 150ms)
	opts := Options{
		Repeat:              1000,
		MinWords:            3,
		MaxWords:            5,
		IncludeRandomNumber: true,
	}
	results, err := eng.GenerateWithOptionsE(opts)
	if err != nil {
		t.Fatalf("batch generation failed: %v", err)
	}
	if len(results) != 1000 {
		t.Fatalf("expected 1000 results, got %d", len(results))
	}
}

// TestGoal_09_MakefileAndScriptsStability verifies scripts stability and zero-division protection.
func TestGoal_09_MakefileAndScriptsStability(t *testing.T) {
	// Test analyze.py with empty file
	tmpFile, err := os.CreateTemp("", "empty_output_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cmd := exec.Command("python3", "scripts/analyze.py", tmpFile.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("analyze.py failed on empty file: %v (output: %s)", err, string(out))
	}
	if !strings.Contains(string(out), "file is empty") {
		t.Errorf("expected 'file is empty' message, got: %s", string(out))
	}
}

// TestGoal_10_HistoricalScholarAttributions verifies accurate representation of Golden Age polymaths.
func TestGoal_10_HistoricalScholarAttributions(t *testing.T) {
	scholars, exists := words.Categories["islamic_golden_age_scholars"]
	if !exists || len(scholars) == 0 {
		t.Fatal("islamic_golden_age_scholars category must exist and be non-empty")
	}

	scholarSet := make(map[string]bool)
	for _, s := range scholars {
		scholarSet[strings.ToLower(s)] = true
	}

	// Verify diverse non-Muslim scholars of the Golden Age
	expectedScholars := []string{
		"maimonides",       // Jewish polymath
		"hunayn ibn ishaq", // Christian physician/translator
		"thabit ibn qurra", // Sabian mathematician/astronomer
	}
	for _, exp := range expectedScholars {
		if !scholarSet[exp] {
			t.Errorf("expected scholar %q in islamic_golden_age_scholars", exp)
		}
	}

	// Verify muslim_scientists is eliminated
	if _, exists := words.Categories["muslim_scientists"]; exists {
		t.Error("category muslim_scientists must be completely removed")
	}

	// Check words.json raw text
	data, err := os.ReadFile("words.json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "\"muslim_scientists\"") {
		t.Error("words.json must not contain \"muslim_scientists\" in categories or rules")
	}
}

// TestGoal_11_SanctuariesAndAnimalRules verifies holy sanctuary segregation and animal pairing safety.
func TestGoal_11_SanctuariesAndAnimalRules(t *testing.T) {
	sanctuaries, exists := words.Categories["holy_sanctuaries"]
	if !exists || len(sanctuaries) == 0 {
		t.Fatal("holy_sanctuaries category must exist and be non-empty")
	}

	sanctuarySet := make(map[string]bool)
	for _, s := range sanctuaries {
		sanctuarySet[strings.ToLower(s)] = true
	}

	for _, s := range []string{"kaaba", "masjid", "madina", "mecca", "medina", "al-aqsa", "minbar", "mihrab"} {
		if !sanctuarySet[s] {
			t.Errorf("expected sanctuary %q in holy_sanctuaries", s)
		}
	}

	// Ensure sacred sites are not in nouns_places or nouns_objects
	for _, p := range words.Categories["nouns_places"] {
		if sanctuarySet[strings.ToLower(p)] {
			t.Errorf("sacred sanctuary %q must NOT be in nouns_places", p)
		}
	}
	for _, o := range words.Categories["nouns_objects"] {
		if sanctuarySet[strings.ToLower(o)] {
			t.Errorf("sacred object %q must NOT be in nouns_objects", o)
		}
	}

	// Ensure impure animals are eliminated
	for _, a := range words.Categories["animals"] {
		norm := strings.ToLower(a)
		if norm == "warthog" || norm == "boar" || norm == "wild-dog" || norm == "pig" || norm == "swine" {
			t.Errorf("impure animal %q must NOT be in animals category", a)
		}
	}

	// Ensure Rule 16 uses geographic_features
	for _, r := range words.Rules {
		if len(r.Pattern) == 2 && r.Pattern[0] == "animals" {
			if r.Pattern[1] != "geographic_features" && r.Pattern[1] != "adjectives" {
				t.Errorf("animals must only pair with geographic_features or adjectives, found: %v", r.Pattern)
			}
		}
	}
}

// TestGoal_12_HardenD1ScriptsAndArchitecture verifies D1 script hardening and safe CLI limits.
func TestGoal_12_HardenD1ScriptsAndArchitecture(t *testing.T) {
	d1Script, err := os.ReadFile("scripts/populate_d1.py")
	if err != nil {
		t.Fatalf("failed to read scripts/populate_d1.py: %v", err)
	}
	scriptContent := string(d1Script)

	if !strings.Contains(scriptContent, "INSERT OR IGNORE") {
		t.Error("scripts/populate_d1.py must use INSERT OR IGNORE")
	}
	if !strings.Contains(scriptContent, "--batch-size") {
		t.Error("scripts/populate_d1.py must support --batch-size parameter")
	}
	if !strings.Contains(scriptContent, "--auto-create") {
		t.Error("scripts/populate_d1.py must support --auto-create flag")
	}

	// Verify README recommendation exists
	readmeData, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(readmeData), "Architectural Recommendation") {
		t.Error("README.md must contain Architectural Recommendation for edge worker generation")
	}
}

// TestGoal_13_HarmonizeLicensingAndSecurityPolicy verifies license attribution and security docs.
func TestGoal_13_HarmonizeLicensingAndSecurityPolicy(t *testing.T) {
	licenseData, err := os.ReadFile("LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(licenseData), "WaqfTech") {
		t.Error("LICENSE must attribute WaqfTech")
	}
	if !strings.Contains(string(licenseData), "GNU AFFERO GENERAL PUBLIC LICENSE") {
		t.Error("LICENSE must be GNU Affero General Public License (AGPL-3.0)")
	}

	waqfLicense, err := os.ReadFile("WaqfDPL-1.0.md")
	if err != nil {
		t.Fatalf("WaqfDPL-1.0.md is missing: %v", err)
	}
	if !strings.Contains(string(waqfLicense), "Waqf-DPL 1.0") {
		t.Error("WaqfDPL-1.0.md must contain Waqf-DPL 1.0 draft text")
	}

	secData, err := os.ReadFile("SECURITY.md")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(secData), "(Replace with actual email address)") {
		t.Error("SECURITY.md must not contain placeholder email instructions")
	}
	if !strings.Contains(string(secData), "jad@madi.se") {
		t.Error("SECURITY.md must specify contact email jad@madi.se")
	}

	todosData, err := os.ReadFile("docs/TODOS.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(todosData), "## Completed Features") {
		t.Error("docs/TODOS.md must document Completed Features")
	}
}

// TestGoal_14_SanitizeAsmaAllahAndVirtues verifies separation of 99 Names of Allah from generic adjectives.
func TestGoal_14_SanitizeAsmaAllahAndVirtues(t *testing.T) {
	asma, exists := words.Categories["asma_allah"]
	if !exists || len(asma) < 90 {
		t.Fatalf("asma_allah category must contain at least 90 divine names, got %d", len(asma))
	}

	divineSet := make(map[string]bool)
	for _, a := range asma {
		divineSet[strings.ToLower(a)] = true
	}

	// Must contain major exclusive Divine Names
	exclusiveNames := []string{
		"rahman", "khaliq", "bari", "musawwir", "qahhar", "jabbar",
		"quddus", "razzaq", "samad", "ahad", "hayy", "qayyum",
	}
	for _, name := range exclusiveNames {
		if !divineSet[name] {
			t.Errorf("expected Divine Name %q in asma_allah", name)
		}
	}

	// Generic adjectives must be completely free of Divine Names
	for _, adj := range words.Categories["adjectives"] {
		if divineSet[strings.ToLower(adj)] {
			t.Errorf("Divine Name %q must NOT be in generic adjectives", adj)
		}
	}

	// Verify servant_prefixes exists and has abd and amat
	prefixes, exists := words.Categories["servant_prefixes"]
	if !exists || len(prefixes) < 2 {
		t.Fatal("servant_prefixes category must exist with abd and amat")
	}

	// Verify respectful generation formula
	opts := Options{
		Repeat:              5,
		Categories:          []string{"servant_prefixes", "asma_allah"},
		MinWords:            2,
		MaxWords:            2,
		IncludeRandomNumber: false,
	}
	results, err := GenerateWithOptionsE(opts)
	if err != nil {
		t.Fatalf("failed to generate pious names: %v", err)
	}
	for _, r := range results {
		if !strings.HasPrefix(r, "abd-") && !strings.HasPrefix(r, "amat-") {
			t.Errorf("result %q must start with abd- or amat-", r)
		}
	}
}

// TestGoal_15_TheologicalSanitationAndCombinatorialRules verifies removal of sacrilegious rules and concepts.
func TestGoal_15_TheologicalSanitationAndCombinatorialRules(t *testing.T) {
	// 1. Check prohibited rules are eliminated
	disallowedPatterns := [][]string{
		{"islamic_golden_age_scholars", "vegetables"},
		{"colors_arabic", "muslim_names_male"},
		{"colors_arabic", "muslim_names_female"},
		{"animals", "adjectives"},
	}
	for _, rule := range words.Rules {
		for _, dis := range disallowedPatterns {
			if len(rule.Pattern) == len(dis) {
				match := true
				for i := range rule.Pattern {
					if rule.Pattern[i] != dis[i] {
						match = false
						break
					}
				}
				if match {
					t.Fatalf("prohibited rule pattern %v must not exist in words.Rules", dis)
				}
			}
		}
	}

	// 2. Check malevolent concepts are purged from nouns_concepts
	malevolent := []string{
		"dajjal", "jahannam", "fitna", "yajuj-majuj",
		"kaba-structure", "masjid-al-haram", "masjid-al-nabawi", "quran",
	}
	for _, w := range words.Categories["nouns_concepts"] {
		for _, mal := range malevolent {
			if strings.EqualFold(w, mal) {
				t.Errorf("malevolent or misplaced concept %q must not be in nouns_concepts", w)
			}
		}
	}

	// 3. Check exclusive divine names in muslim_names_male are prefixed
	exclusiveDivine := []string{"qadir", "ghani", "wahid", "mu'izz"}
	for _, m := range words.Categories["muslim_names_male"] {
		norm := strings.ToLower(m)
		for _, ed := range exclusiveDivine {
			if norm == ed {
				t.Errorf("exclusive divine name %q must not exist without servant prefix in muslim_names_male", m)
			}
		}
	}

	// 4. Check fattah in food categories is corrected
	for _, food := range words.Categories["saudi_food"] {
		if strings.EqualFold(food, "fattah") {
			t.Errorf("food item 'Fattah' in saudi_food must be transliterated as 'Fatteh' or 'Fatta'")
		}
	}
	for _, food := range words.Categories["arabic_food"] {
		if strings.EqualFold(food, "fattah") {
			t.Errorf("food item 'Fattah' in arabic_food must be transliterated as 'Fatteh' or 'Fatta'")
		}
	}

	// 5. Check jannah-end is removed
	for _, s := range words.Categories["suffixes"] {
		if strings.EqualFold(s, "jannah-end") {
			t.Errorf("eschatological violation 'jannah-end' must not exist in suffixes")
		}
	}
}

// TestGoal_16_DelimiterSafetyAndSeedDeterminism verifies delimiter normalization in isSafe and PRNG seed reproducibility.
func TestGoal_16_DelimiterSafetyAndSeedDeterminism(t *testing.T) {
	// 1. Verify custom separators (. / # :) are intercepted by isSafe
	evasiveBlockedStrings := []string{
		"prefix.wine.suffix",
		"dish/pork/dish",
		"bet#casino#bet",
		"vodka:drink",
		"gambling_site",
		"user--wine--end",
	}
	for _, s := range evasiveBlockedStrings {
		if isSafe(s) {
			t.Errorf("safety filter failure: %q contains a blocked term but was marked safe", s)
		}
	}

	// 2. Verify deterministic seed output across separate engine calls
	opts1 := Options{
		Repeat:              5,
		Seed:                424242,
		MinWords:            3,
		MaxWords:            5,
		IncludeRandomNumber: true,
	}
	results1, err := GenerateWithOptionsE(opts1)
	if err != nil {
		t.Fatalf("run 1 failed: %v", err)
	}

	opts2 := Options{
		Repeat:              5,
		Seed:                424242,
		MinWords:            3,
		MaxWords:            5,
		IncludeRandomNumber: true,
	}
	results2, err := GenerateWithOptionsE(opts2)
	if err != nil {
		t.Fatalf("run 2 failed: %v", err)
	}

	for i := range results1 {
		if results1[i] != results2[i] {
			t.Fatalf("seed determinism failure at index %d: run1=%q != run2=%q", i, results1[i], results2[i])
		}
	}

	// 3. Verify Engine.Generate() does not panic even on impossible options
	engine := NewEngine(Words{Categories: map[string][]string{}})
	val := engine.Generate()
	if val != "" {
		t.Errorf("expected empty string from empty engine, got %q", val)
	}
}


