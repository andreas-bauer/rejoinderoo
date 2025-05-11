package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/tui"
)

func main() {
	inFileFlag := flag.String("i", "", "file path to input file (CSV or Excel)")
	flag.Parse()

	inFile := *inFileFlag
	if *inFileFlag == "" {
		inFile = tui.RunFilePicker()
	}

	td, err := reader.NewReader(inFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		fmt.Println("Unable to proceed.")

		os.Exit(1)
	}

	fd := &tui.FormData{
		AvailableHeaders: td.Headers,
	}
	err = tui.RunForm(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running TUI form:", err)
		os.Exit(1)
	}

	tui.PrintSummary(fd)
}
