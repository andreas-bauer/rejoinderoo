package reader

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelReader struct{}

func (r ExcelReader) Read(file io.Reader) (*TabularData, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
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
		Records: rows[1:], // skip headers
	}, nil
}
