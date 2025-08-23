package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/WaqfTech/halal-random-strings"
	"github.com/charmbracelet/x/exp/higherorder"
	"github.com/charmbracelet/x/exp/ordered"
	"github.com/dustin/go-humanize"
	"github.com/mattn/go-runewidth"
	flag "github.com/spf13/pflag"
)

func formatPossibilities(o halalrandomstrings.Options) string {
	low, high := halalrandomstrings.PossibilitiesWithOptions(o)
	return fmt.Sprintf(
		"Minimum combinations: %s\nMaximum combinations: %s",
		humanize.Comma(int64(low)),
		humanize.Comma(int64(high)),
	)
}

func usage(o halalrandomstrings.Options) {
	fmt.Fprintf(os.Stderr, "Usage: %s [FLAGS]", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "\n\n%s\n\n", formatPossibilities(o))
	fmt.Fprintln(os.Stderr, "Options:")

	var flags []string
	var descs []string
	flag.VisitAll(func(f *flag.Flag) {
		if f.Shorthand == "" {
			flags = append(flags, fmt.Sprintf("--%s", f.Name))
		} else {
			flags = append(flags, fmt.Sprintf("-%s --%s", f.Shorthand, f.Name))
		}
		descs = append(descs, f.Usage)
	})

	max := func(a, b int) int { return ordered.Max(a, b) }
	widestFlag := higherorder.Foldl(max, 0, higherorder.Map(runewidth.StringWidth, flags))

	const gap = 2
	for i := range flags {
		fmt.Fprintf(os.Stderr, "  %-*s %s\n", widestFlag+gap, flags[i], descs[i])
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
	)

	flag.BoolVarP(&showHelp, "help", "h", false, "Show this help and exit")
	flag.IntVarP(&opts.Repeat, "repeat", "r", defaultResults, "Number of strings to generate")
	flag.StringVar(&opts.Sep, "sep", "-", "Separator to use between words")
	flag.Int64Var(&opts.Seed, "seed", 0, "Optional int64 seed for reproducibility")
	flag.Float64VarP(&opts.PrefixThreshold, "prefix-threshold", "p", 0.2, "How often to include bonus prefixes (0.2)")
	flag.Float64VarP(&opts.SuffixThreshold, "suffix-threshold", "s", 0.2, "How often to include bonus suffixes (0.2)")

	flag.CommandLine.SortFlags = false
	flag.Usage = func() {
		usage(opts)
	}
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(1)
	}

	opts.Repeat = ordered.Clamp(opts.Repeat, minResults, maxResults)

	r := halalrandomstrings.GenerateWithOptions(opts)
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i])
	}
}
