package reader

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

// ExcelReader implements FileReader for Excel files
type ExcelReader struct {
	file    io.Reader
	Headers []string
	Records [][]string
}

func NewExcelReader(file io.Reader) (*ExcelReader, error) {
	r := &ExcelReader{
		file:    file,
		Headers: []string{},
		Records: [][]string{},
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
		r.Headers = rows[0]
	}
	r.Records = rows[1:]

	return r, nil
}
