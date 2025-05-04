package reader

import (
	"encoding/csv"
	"io"
)

type CSVReader struct{}

func (r CSVReader) Read(file io.Reader) (*TabularData, error) {
	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &TabularData{
			Headers: []string{},
			Records: [][]string{},
		}, nil
	}

	return &TabularData{
		Headers: rows[0],
		Records: rows[1:],
	}, nil
}
