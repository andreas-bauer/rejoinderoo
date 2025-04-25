package reader

import (
	"encoding/csv"
	"io"
)

// CSVReader implements FileReader for CSV files
type CSVReader struct {
	file    io.Reader
	headers []string
	records [][]string
}

// NewCSVReader creates a new CSVReader from an io.Reader
func NewCSVReader(file io.Reader) (*CSVReader, error) {
	r := &CSVReader{
		file:    file,
		headers: []string{},
		records: [][]string{},
	}

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	r.headers = headers

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	r.records = records

	return r, nil
}

// HasData checks if the CSV file has any data rows beyond the header
func (r *CSVReader) HasData() bool {
	return len(r.records) > 0
}

// Headers returns the column headers from the CSV file
func (r *CSVReader) Headers() []string {
	return r.headers
}

// Records returns the data rows from the CSV file
func (r *CSVReader) Records() [][]string {
	return r.records
}
