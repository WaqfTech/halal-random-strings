package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaqfTech/halal-random-strings"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/ordered"
	flag "github.com/spf13/pflag"
)

	func usage() {
	appName := filepath.Base(os.Args[0])

	// Define styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")). // Green
		Bold(true).
		PaddingBottom(1)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA")). // Light Gray
		PaddingBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")). // Cyan
		Bold(true).
		PaddingTop(1).
		PaddingBottom(1)

	exampleCommandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")). // Yellow
		Bold(true)

	commentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")). // Dark Gray
		Italic(true)

	sampleOutputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")). // Green
		Faint(true)

	noteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")). // Orange
		Italic(true).
		PaddingTop(1)

	// Print Usage and Description
	fmt.Fprintln(os.Stderr, titleStyle.Render(fmt.Sprintf("Usage: %s [FLAGS]", appName)))
	fmt.Fprintln(os.Stderr, descriptionStyle.Render("A wholesome, family-friendly random string generator."))
	fmt.Fprintln(os.Stderr, descriptionStyle.Render("Generates unique and meaningful identifiers using Islamic terms and concepts."))

	// Print Options
	fmt.Fprintln(os.Stderr, sectionStyle.Render("Options:"))
	flag.PrintDefaults()

	// Print Examples
	fmt.Fprintln(os.Stderr, sectionStyle.Render("Examples:"))

	// Helper function to print examples
	printExample := func(cmd, desc string, outputs ...string) {
		fmt.Fprintln(os.Stderr, exampleCommandStyle.Render(fmt.Sprintf("  %s", cmd)))
		fmt.Fprintln(os.Stderr, commentStyle.Render(fmt.Sprintf("    # %s", desc)))
		if len(outputs) > 0 {
			fmt.Fprintln(os.Stderr, sampleOutputStyle.Render("    # Sample Output:"))
			for _, out := range outputs {
				fmt.Fprintln(os.Stderr, sampleOutputStyle.Render(fmt.Sprintf("    #   %s", out)))
			}
		}
		fmt.Fprintln(os.Stderr, "") // Add a blank line for spacing
	}

	printExample(
		appName,
		"Generates a single random string with a default length of 5-8 words and a numerical ending.",
		"wisdom-qalam-light-station-123456789",
	)

	printExample(
		fmt.Sprintf("%s -r 3", appName),
		"Generates 3 random strings.",
		"salam-peace-mosque-987654321",
		"iman-truth-garden-123456789",
		"barakah-blessing-river-567890123",
	)

	printExample(
		fmt.Sprintf("%s --sep _", appName),
		"Uses an underscore as the separator between words instead of a hyphen.",
		"subhanallah_glory_mountain_456789012",
	)

	printExample(
		fmt.Sprintf("%s --seed 12345", appName),
		"Generates a string using a specific seed, useful for reproducing results.",
		"(Output will be consistent for seed 12345)",
	)

	printExample(
		fmt.Sprintf("%s --min-words 3 --max-words 5", appName),
		"Generates a string with 3 to 5 words.",
		"jannah-paradise-tree-789012345",
	)

	printExample(
		fmt.Sprintf("%s --no-random-number", appName),
		"Generates a string without the numerical ending.",
		"masjid-prayer-carpet",
	)

	printExample(
		fmt.Sprintf("%s --categories sahaba", appName),
		"Generates a string using only Sahaba names.",
		"umar-ibn-al-khattab-987654321",
	)

	printExample(
		fmt.Sprintf("%s --categories adjectives,nouns_places", appName),
		"Combines words from 'adjectives' and 'nouns_places' categories.",
		"beautiful-kaaba-city-345678901",
	)

	fmt.Fprintln(os.Stderr, noteStyle.Render("Note: Sample outputs are illustrative and actual outputs may vary due to randomness."))

	// Print Available Categories
	fmt.Fprintln(os.Stderr, sectionStyle.Render("Available Categories:"))
	categories := halalrandomstrings.GetCategories()
	for _, category := range categories {
		fmt.Fprintln(os.Stderr, commentStyle.Render(fmt.Sprintf("  - %s", category)))
	}
}

func main() {
	const (
		minResults     = 1
		maxResults     = 1000
		defaultResults = 1
	)

	var (
		showHelp bool
		opts     halalrandomstrings.Options
		categoriesStr string
		noRandomNumber bool
	)

	flag.BoolVarP(&showHelp, "help", "h", false, "Show this help and exit")
	flag.IntVarP(&opts.Repeat, "repeat", "r", defaultResults, "Number of strings to generate")
	flag.StringVar(&opts.Sep, "sep", "-", "Separator to use between words")
	flag.Int64Var(&opts.Seed, "seed", 0, "Optional int64 seed for reproducibility")
	flag.IntVar(&opts.MinWords, "min-words", 0, "Minimum number of words in the generated string (default: 5)")
	flag.IntVar(&opts.MaxWords, "max-words", 0, "Maximum number of words in the generated string (default: 8)")
	flag.StringVar(&categoriesStr, "categories", "", "Comma-separated list of categories to use (e.g., adjectives,nouns_places)")
	flag.BoolVar(&noRandomNumber, "no-random-number", false, "Do not append a random number to the end of the string")

	flag.CommandLine.SortFlags = false
	flag.Usage = usage // Set custom usage function
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0) // Exit with 0 for help
	}

	// Apply defaults if not set by flags
	if opts.MinWords == 0 {
		opts.MinWords = 5
	}
	if opts.MaxWords == 0 {
		opts.MaxWords = 8
	}

	// Parse categories string
	if categoriesStr != "" {
		opts.Categories = strings.Split(categoriesStr, ",")
	}

	// Set random number option
	opts.IncludeRandomNumber = !noRandomNumber

	opts.Repeat = ordered.Clamp(opts.Repeat, minResults, maxResults)

	r, err := halalrandomstrings.GenerateWithOptionsE(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i])
	}
}