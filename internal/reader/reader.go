package reader

import "io"

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
func NewReader(file io.ReadSeeker, filename string) (FileReader, error) {
	// Implementation would determine the file type from extension
	// and return the appropriate reader implementation
	// return &CSVReader{...} or &ExcelReader{...}
	return nil, nil
}
