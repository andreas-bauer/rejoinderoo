package reader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileReader provides a common interface for reading tabular data from different file formats.
// It supports operations like getting headers and checking if the file has data.
type FileReader interface {
	// Headers returns the column headers from the file.
	Headers() []string

	// HasData checks if the file contains any data rows beyond the header.
	HasData() bool

	// Records returns the data rows from the file.
	Records() [][]string
}

// NewReader creates an appropriate FileReader based on the file extension.
// Supported formats: CSV, XLSX, XLS
func NewReader(filename string) (FileReader, error) {
	filename = strings.ToLower(filename)

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create appropriate reader based on file extension
	fileExt := filepath.Ext(filename)
	switch fileExt {
	case ".csv":
		return NewCSVReader(file)
	case ".xlsx", ".xls":
		return NewExcelReader(file)
	default:
		return nil, errors.New("Input file must be a CSV (*.csv) or Excel (*.xlsx or *.xls) file, but got: " + filename)
	}
}
