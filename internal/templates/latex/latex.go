package latex

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates"
)

// Latex handles escaping special characters for LaTeX templates.
type Latex struct{}

type header struct {
	Name string
	Idx  int
}

type record struct {
	Header string
	Text   string
}

type response struct {
	ReviewerID string
	Records    []record
}

type document struct {
	ReviewerIDs []string
	LenHeaders  int
	Headers     []header
	Responses   []response
}

func NewLatexTemplate() *Latex {
	return &Latex{}
}

const (
	filename = "./internal/templates/latex/latex.tmpl"
)

func (l *Latex) Render(td reader.TabularData) (string, error) {

	escapeAllStrings(&td)
	doc := createDoc(&td)

	tmpl, err := template.ParseFiles(filename)

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
	allRevIDs := templates.ExtractReviewers(td.Records)
	headers := asDocHeaders(td.Headers)
	responses := asDocResponses(td.Headers, td.Records)

	return document{
		ReviewerIDs: allRevIDs,
		LenHeaders:  len(td.Headers),
		Headers:     headers,
		Responses:   responses,
	}
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

// escape escapes special characters for LaTeX.
func escape(input string) string {
	// Replace LaTeX special characters with their escaped versions.
	replacer := strings.NewReplacer(
		"\\", "\\textbackslash{}",
		"{", "\\{",
		"}", "\\}",
		"$", "\\$",
		"&", "\\&",
		"#", "\\#",
		"_", "\\_",
		"%", "\\%",
		"~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}",
	)
	return replacer.Replace(input)
}

// asDocHeaders converts a slice of strings to a slice of Header structs
func asDocHeaders(headers []string) []header {
	var res = make([]header, len(headers))
	for idx, h := range headers {
		header := &header{
			Name: h,
			Idx:  idx + 1,
		}
		res[idx] = *header
	}
	return res
}

// asDocResponses converts a slice of strings to a slice of Response structs
func asDocResponses(headers []string, records [][]string) []response {
	var res = make([]response, len(records))
	for idx, rec := range records {
		if len(rec) < 1 {
			continue
		}
		response := &response{
			ReviewerID: rec[0],
			Records:    make([]record, len(headers)),
		}
		for i, h := range headers {
			var text string
			if i < len(rec) {
				text = rec[i]
			}
			record := &record{
				Header: h,
				Text:   text,
			}
			response.Records[i] = *record
		}
		res[idx] = *response
	}
	return res

}

