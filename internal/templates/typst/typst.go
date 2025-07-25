package typst

import (
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates/common"
)

// Typst handles escaping special characters for Typst templates.
type Typst struct{}

type record struct {
	Header string
	Text   string
}

type response struct {
	ID         string
	ReviewerID string
	Records    []record
}

type document struct {
	ReviewerIDs []string
	Responses   []response
}

//go:embed typst.tmpl
var file string

func NewTypstTemplate() *Typst {
	return &Typst{}
}

// FileExtension returns the file extension for Typst templates.
func (t *Typst) FileExtension() string {
	return ".typ"
}

// Render processes the Typst template with the provided tabular data.
func (t *Typst) Render(td reader.TabularData) (string, error) {

	escapeAllStrings(&td)
	doc := createDoc(&td)

	tmpl, err := template.New("typst").Parse(file)

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	var result strings.Builder
	err = tmpl.Execute(&result, doc)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return result.String(), nil
}

func createDoc(td *reader.TabularData) document {
	allRevIDs := common.ExtractReviewers(td.Records)
	responses := asDocResponses(td.Headers, td.Records)

	return document{
		ReviewerIDs: allRevIDs,
		Responses:   responses,
	}
}

// asDocResponses converts a slice of strings to a slice of Response structs
func asDocResponses(headers []string, records [][]string) []response {
	var res = make([]response, len(records))
	for idx, rec := range records {
		if len(rec) < 1 {
			continue
		}
		response := &response{
			ID:         rec[0],
			ReviewerID: common.ExtractReviewerID(rec[0]),
			Records:    make([]record, len(headers)-1),
		}
		for i, h := range headers {
			if i == 0 {
				continue // skip ID
			}
			var text string
			if i < len(rec) {
				text = rec[i]
			}
			record := &record{
				Header: h,
				Text:   text,
			}
			response.Records[i-1] = *record
		}
		res[idx] = *response
	}
	return res

}

func escapeAllStrings(td *reader.TabularData) {
	for i, h := range td.Headers {
		td.Headers[i] = escape(h)
	}

	for i, rec := range td.Records {
		for j, r := range rec {
			td.Records[i][j] = escape(r)
		}
	}
}

// Escape escapes special characters for Typst.
func escape(input string) string {
	// Replace Typst special characters with their escaped versions.
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"#", "\\#",
	)
	return replacer.Replace(input)
}
