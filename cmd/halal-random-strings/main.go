package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaqfTech/halal-random-strings"
	"github.com/charmbracelet/lipgloss"
	flag "github.com/spf13/pflag"
)

func printUsage(out io.Writer) {
	appName := filepath.Base(os.Args[0])
	useColor := os.Getenv("NO_COLOR") == "" && os.Getenv("TERM") != "dumb"

	// Define accessible styles readable across dark and light palettes
	var (
		titleStyle          lipgloss.Style
		descriptionStyle    lipgloss.Style
		sectionStyle        lipgloss.Style
		exampleCommandStyle lipgloss.Style
		commentStyle        lipgloss.Style
		sampleOutputStyle   lipgloss.Style
		noteStyle           lipgloss.Style
	)

	if useColor {
		titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2E7D32")). // Deep Forest Green
			Bold(true).
			PaddingBottom(1)

		descriptionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			PaddingBottom(1)

		sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00838F")). // Deep Cyan
			Bold(true).
			PaddingTop(1).
			PaddingBottom(1)

		exampleCommandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C58500")). // Dark Amber / Gold
			Bold(true)

		commentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")).
			Italic(true)

		sampleOutputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2E7D32")).
			Faint(true)

		noteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D84315")). // Dark Orange
			Italic(true).
			PaddingTop(1)
	}

	// Print Usage and Description
	fmt.Fprintln(out, titleStyle.Render(fmt.Sprintf("Usage: %s [FLAGS]", appName)))
	fmt.Fprintln(out, descriptionStyle.Render("A wholesome, family-friendly random string generator."))
	fmt.Fprintln(out, descriptionStyle.Render("Each identifier uses one approved mixing group. Outputs are not secret tokens."))

	// Print Options
	fmt.Fprintln(out, sectionStyle.Render("Options:"))
	flag.CommandLine.SetOutput(out)
	flag.PrintDefaults()

	// Print Examples
	fmt.Fprintln(out, sectionStyle.Render("Examples:"))

	// Helper function to print examples
	printExample := func(cmd, desc string, outputs ...string) {
		fmt.Fprintln(out, exampleCommandStyle.Render(fmt.Sprintf("  %s", cmd)))
		fmt.Fprintln(out, commentStyle.Render(fmt.Sprintf("    # %s", desc)))
		if len(outputs) > 0 {
			fmt.Fprintln(out, sampleOutputStyle.Render("    # Sample Output:"))
			for _, sample := range outputs {
				fmt.Fprintln(out, sampleOutputStyle.Render(fmt.Sprintf("    #   %s", sample)))
			}
		}
		fmt.Fprintln(out, "")
	}

	printExample(
		appName,
		"Generates a single random string with a default length of 5-8 words and a numerical ending.",
		"olive-apple-pear-plum-peach-1234",
	)

	printExample(
		fmt.Sprintf("%s -r 3", appName),
		"Generates 3 random strings.",
		"aisha-maher-amira-fatima-layla-9876",
		"rahman-rahim-quddus-salam-samad-1234",
		"ahmad-musa-isa-nuh-hud-5678",
	)

	printExample(
		fmt.Sprintf("%s --sep _", appName),
		"Uses an underscore as the separator between words instead of a hyphen.",
		"olive_apple_pear_plum_peach_4567",
	)

	printExample(
		fmt.Sprintf("%s --seed 12345", appName),
		"Generates a string using a specific seed, useful for reproducing results.",
		"(Output will be consistent for seed 12345)",
	)

	printExample(
		fmt.Sprintf("%s --min-words 3 --max-words 5", appName),
		"Generates a string with 3 to 5 words.",
		"olive-apple-pear-7890",
	)

	printExample(
		fmt.Sprintf("%s --no-random-number", appName),
		"Generates a string without the numerical ending.",
		"olive-apple-pear-plum-peach",
	)

	printExample(
		fmt.Sprintf("%s --categories sahaba,islamic_months --min-words 5 --max-words 5", appName),
		"Combines only reviewed religious categories.",
		"umar-ibn-al-khattab-ramadan-9876",
	)

	printExample(
		fmt.Sprintf("%s --categories colors_arabic,nouns_places", appName),
		"Combines general categories; every requested category contributes.",
		"akhdar-city-town-house-azraq-3456",
	)

	fmt.Fprintln(out, noteStyle.Render("Note: Sample outputs are illustrative and actual outputs may vary due to randomness."))
	fmt.Fprintln(out, "Mixing groups: Asma Allah only; prophets only; personal names only; Islamic categories only; general categories only.")
	fmt.Fprintln(out, "Custom dictionaries and standalone servant_prefixes are disabled. See docs/MIXING_POLICY.md for the complete category table.")

	// Print Available Categories
	fmt.Fprintln(out, sectionStyle.Render("Available Categories:"))
	categories := halalrandomstrings.GetCategories()
	for _, category := range categories {
		fmt.Fprintln(out, commentStyle.Render(fmt.Sprintf("  - %s", category)))
	}
}

func usage() {
	printUsage(os.Stderr)
}

func main() {
	const (
		defaultResults = 1
	)

	var (
		showHelp       bool
		opts           halalrandomstrings.Options
		categoriesStr  string
		noRandomNumber bool
	)

	flag.BoolVarP(&showHelp, "help", "h", false, "Show this help and exit")
	flag.IntVarP(&opts.Repeat, "repeat", "r", defaultResults, "Number of strings to generate")
	flag.StringVar(&opts.Sep, "sep", "-", "Separator: -, _, . or / only")
	flag.Int64Var(&opts.Seed, "seed", 0, "Optional int64 seed for reproducibility")
	flag.IntVar(&opts.MinWords, "min-words", 0, "Minimum number of words in the generated string (default: 5)")
	flag.IntVar(&opts.MaxWords, "max-words", 0, "Maximum number of words in the generated string (default: 8)")
	flag.StringVar(&categoriesStr, "categories", "", "Comma-separated reviewed categories from one mixing group (e.g., colors_arabic,nouns_places)")
	flag.BoolVar(&noRandomNumber, "no-random-number", false, "Do not append a random number to the end of the string")

	flag.CommandLine.SortFlags = false
	flag.Usage = usage
	flag.Parse()

	if showHelp {
		printUsage(os.Stdout)
		os.Exit(0)
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

	r, err := halalrandomstrings.GenerateWithOptionsE(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i])
	}
}
