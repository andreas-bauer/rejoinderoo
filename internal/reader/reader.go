package reader

import (
	"fmt"
	"io"
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
func NewReader(filename string) (TabularReader, error) {
	filename = strings.ToLower(filename)
	fileExt := filepath.Ext(filename)
	switch fileExt {
	case ".csv":
		return &CSVReader{}, nil
	case ".xlsx", ".xls":
		return &ExcelReader{}, nil
	default:
		return nil, fmt.Errorf("file extension '%s' is not supported. Supported extensions are: %v", fileExt, SupportedFileExtensions())
	}
}

// SupportedFileExtensions returns a list of file extensions that are supported.
func SupportedFileExtensions() []string {
	return []string{".csv", ".xlsx", ".xls"}
}

// Keep filters the headers and records by removing all headers and the corresponding records
// that are not in the given list of headers to keep.
// It changes the order of the headers and records to match the order of headers to keep.
func (td *TabularData) Keep(headers []string) {
	// Map header to its index in the original headers
	headerIndex := make(map[string]int)
	for i, h := range td.Headers {
		headerIndex[h] = i
	}

	newHeaders := make([]string, 0, len(headers))
	indicesToKeep := make([]int, 0, len(headers))
	for _, h := range headers {
		if idx, ok := headerIndex[h]; ok {
			newHeaders = append(newHeaders, h)
			indicesToKeep = append(indicesToKeep, idx)
		}
	}

	newRecords := make([][]string, len(td.Records))
	for i, record := range td.Records {
		newRecord := make([]string, 0, len(indicesToKeep))
		for _, idx := range indicesToKeep {
			if idx < len(record) {
				newRecord = append(newRecord, record[idx])
			} else {
				newRecord = append(newRecord, "")
			}
		}
		newRecords[i] = newRecord
	}

	td.Headers = newHeaders
	td.Records = newRecords
}
