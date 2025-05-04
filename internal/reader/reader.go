package reader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TabularReader defines an interface for reading tabular files (CSV, Excel, etc.)
type TabularReader interface {
	Read(file io.Reader) (*TabularData, error)
}

// TabularData represents the structure of a spreadsheet file with headers and records.
type TabularData struct {
	Headers []string
	Records [][]string
}

// NewReader creates an appropriate TabularReader based on the file extension.
// Supported formats: CSV, XLSX, XLS
func NewReader(filename string) (*TabularData, error) {
	filename = strings.ToLower(filename)

	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var r TabularReader
	// Create appropriate reader based on file extension
	fileExt := filepath.Ext(filename)
	switch fileExt {
	case ".csv":
		r = CSVReader{}
	case ".xlsx", ".xls":
		r = ExcelReader{}
	default:
		return nil, fmt.Errorf("Unsupported file extension: %s", fileExt)
	}

	return r.Read(file)
}
