package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
)

func main() {
	inFile := flag.String("i", "./tests/small.csv", "file path to input file (CSV or Excel)")
	outFile := flag.String("o", "out.tex", "file path to output LaTeX file")
	flag.Parse()

	isCSV := strings.HasSuffix(*inFile, ".csv")
	isExcel := strings.HasSuffix(*inFile, ".xlsx") || strings.HasSuffix(*inFile, ".xls")

	if !isCSV && !isExcel {
		msg := "Error: Input file must be a CSV or Excel file, but got: " + *inFile
		panic(msg)
	}

	file, err := os.Open(*inFile)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("Rejoinderoo")
	fmt.Println("----------------------")
	fmt.Printf("input: %s\n", *inFile)
	fmt.Printf("output: %s\n", *outFile)
	fmt.Println("----------------------")

	r, err := reader.NewCSVReader(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	recs := r.Records()

	fmt.Println("Amount of records:", len(recs))

}
