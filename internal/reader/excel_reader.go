package reader

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

// ExcelReader implements FileReader for Excel files
type ExcelReader struct {
	headers []string
	records [][]string
}

// NewExcelReader creates a new ExcelReader instance
func NewExcelReader(file io.Reader) (*ExcelReader, error) {
	r := &ExcelReader{
		headers: []string{},
		records: [][]string{},
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		r.headers = rows[0]
	}
	r.records = rows[1:]

	return r, nil
}

// HasData checks if the Excel file has any data rows beyond the header
func (r *ExcelReader) HasData() bool {
	return len(r.records) > 0
}

// Headers returns the column headers from the Excel file
func (r *ExcelReader) Headers() []string {
	return r.headers
}

// Records returns the data rows from the Excel file
func (r *ExcelReader) Records() [][]string {
	return r.records
}
