package main

import (
	"fmt"

	"github.com/WaqfTech/halal-random-strings"
	"github.com/charmbracelet/x/exp/ordered"
	flag "github.com/spf13/pflag"
)

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

	flag.Parse()

	if showHelp {
		flag.Usage()
	}

	opts.Repeat = ordered.Clamp(opts.Repeat, minResults, maxResults)

	r := halalrandomstrings.GenerateWithOptions(opts)
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i])
	}
}
