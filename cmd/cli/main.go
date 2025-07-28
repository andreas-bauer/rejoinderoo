package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates"
	"github.com/andreas-bauer/rejoinderoo/internal/tui"
)

func main() {
	inFileFlag := flag.String("i", "", "file path to input file (CSV or Excel)")
	flag.Parse()

	inFile := *inFileFlag
	if *inFileFlag == "" {
		inFile = tui.RunFilePicker()
	}

	reader, err := reader.NewReader(inFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, errR1 := os.Open(inFile)
	td, errR2 := reader.Read(file)
	if errR1 != nil || errR2 != nil {
		fmt.Println("Error reading file:")
		fmt.Println("Unable to proceed.")
		os.Exit(1)
	}
	defer file.Close()

	fd := &tui.FormData{
		AvailableHeaders: td.Headers,
	}
	err = tui.RunForm(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running TUI form:", err)
		os.Exit(1)
	}

	td.Keep(fd.SelectedHeaders)

	tmpl := templates.NewTemplate(fd.Template)

	if strings.TrimSpace(fd.Filename) == "" {
		fd.Filename = "output"
	}
	fd.Filename = appendExtensionIfNotPresent(fd.Filename, tmpl.FileExtension())

	out, err := tmpl.Render(*td)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error rendering template:", err)
		os.Exit(1)
	}

	err = os.WriteFile(fd.Filename, []byte(out), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error saving output file:", err)
		os.Exit(1)
	}

	tui.PrintSummary(fd)
}

func appendExtensionIfNotPresent(filename, ext string) string {
	if !strings.HasSuffix(strings.ToLower(filename), strings.ToLower(ext)) {
		return filename + ext
	}
	return filename
}
