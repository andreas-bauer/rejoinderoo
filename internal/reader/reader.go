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

// Keep filters the headers and records by removing all headers and the corresponding records
// that are not in the given list of headers to keep.
func (td *TabularData) Keep(headers []string) {
	keepMap := make(map[string]bool)
	for _, h := range headers {
		keepMap[h] = true
	}

	var indicesToKeep []int
	for i, header := range td.Headers {
		if keepMap[header] {
			indicesToKeep = append(indicesToKeep, i)
		}
	}

	newHeaders := make([]string, 0, len(indicesToKeep))
	for _, idx := range indicesToKeep {
		newHeaders = append(newHeaders, td.Headers[idx])
	}

	newRecords := make([][]string, len(td.Records))
	for i, record := range td.Records {
		newRecord := make([]string, 0, len(indicesToKeep))
		for _, idx := range indicesToKeep {
			if idx < len(record) {
				newRecord = append(newRecord, record[idx])
			}
		}
		newRecords[i] = newRecord
	}

	td.Headers = newHeaders
	td.Records = newRecords
}
