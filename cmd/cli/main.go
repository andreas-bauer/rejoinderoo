package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
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

	initialModel := tui.NewModel(td.Headers)
	p := tea.NewProgram(initialModel)

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running TUI:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	m, ok := m.(tui.Model)
	if !ok {
		fmt.Fprintln(os.Stderr, "Inconsistent model")
	}
}
