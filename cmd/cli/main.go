package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	inFile := flag.String("i", "", "file path to input file (CSV or Excel)")
	outFile := flag.String("o", "out.tex", "file path to output LaTeX file")
	flag.Parse()

	isCSV := strings.HasSuffix(*inFile, ".csv")
	isExcel := strings.HasSuffix(*inFile, ".xlsx") || strings.HasSuffix(*inFile, ".xls")

	if !isCSV && !isExcel {
		msg := "Error: Input file must be a CSV or Excel file, but got: " + *inFile
		panic(msg)
	}

	file, err := os.Open("input.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	columns := records[0]
	fmt.Println("Columns:")
	for _, column := range columns {
		fmt.Printf("%s\n", column)
	}

	fmt.Println("Rejoinderoo")
	fmt.Println("----------------------")
	fmt.Printf("input: %s\n", *inFile)
	fmt.Printf("output: %s\n", *outFile)
}
