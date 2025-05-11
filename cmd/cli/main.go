package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/tui"
)

func main() {
	inFile := flag.String("i", "./tests/small.csv", "file path to input file (CSV or Excel)")
	outFile := flag.String("o", "out.tex", "file path to output LaTeX file")
	flag.Parse()

	fmt.Printf("input: %s\n", *inFile)
	fmt.Printf("output: %s\n", *outFile)

	td, err := reader.NewReader(*inFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		fmt.Println("Unable to proceed.")

		os.Exit(1)
	}

	formResult, err := tui.RunForm(td.Headers)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running TUI form:", err)
		os.Exit(1)
	}

	tui.PrintSummary(formResult)
}
