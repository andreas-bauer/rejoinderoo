package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
)

// Tmpl is an interface for templates.
type Tmpl interface {
	File() string
	Escape(input string) string
}

// TemplateType defines the types of templates available.
type TemplateType string

const (
	LatexTemplate TemplateType = "latex"
	TypstTemplate TemplateType = "typst"
)

type Header struct {
	Name string
	Idx  int
}

type Record struct {
	Header string
	Text   string
}

type Response struct {
	ReviewerID string
	Records    []Record
}

type Document struct {
	ReviewerIDs []string
	LenHeaders  int
	Headers     []Header
	Responses   []Response
}

func Render(td *reader.TabularData, tmplType TemplateType) {
	var tt Tmpl
	switch tmplType {
	case LatexTemplate:
		tt = NewLatexTemplate()
	case TypstTemplate:
		tt = NewTypstTemplate()
	default:
		panic("Unknown template type")
	}

	// escapeAllStrings(td, tt)
	doc := createDoc(td)

	tmpl, err := template.ParseFiles(tt.File())

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, doc)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func escapeAllStrings(td *reader.TabularData, tt Tmpl) {
	for i, h := range td.Headers {
		td.Headers[i] = tt.Escape(h)
	}

	for i, rec := range td.Records {
		for j, r := range rec {
			td.Records[i][j] = tt.Escape(r)
		}
	}
}

func createDoc(td *reader.TabularData) Document {
	allRevIDs := extractReviewers(td.Records)
	headers := asDocHeaders(td.Headers)
	responses := asDocResponses(td.Headers, td.Records)

	return Document{
		ReviewerIDs: allRevIDs,
		LenHeaders:  len(td.Headers),
		Headers:     headers,
		Responses:   responses,
	}
}

// asDocHeaders converts a slice of strings to a slice of Header structs
func asDocHeaders(headers []string) []Header {
	var res = make([]Header, len(headers))
	for idx, h := range headers {
		header := &Header{
			Name: h,
			Idx:  idx + 1,
		}
		res[idx] = *header
	}
	return res
}

// asDocResponses converts a slice of strings to a slice of Response structs
func asDocResponses(headers []string, records [][]string) []Response {
	var res = make([]Response, len(records))
	for idx, rec := range records {
		if len(rec) < 1 {
			continue
		}
		response := &Response{
			ReviewerID: rec[0],
			Records:    make([]Record, len(headers)),
		}
		for i, h := range headers {
			var text string
			if i < len(rec) {
				text = rec[i]
			} 
			record := &Record{
				Header: h,
				Text:   text,
			}
			response.Records[i] = *record
		}
		res[idx] = *response
	}
	return res

}

func extractReviewers(records [][]string) []string {
	var allRevIDs []string
	for _, rec := range records {
		if len(rec) < 1 {
			continue
		}
		refID := extractReviewerID(rec[0])
		if searchSlice(allRevIDs, refID) == -1 {
			allRevIDs = append(allRevIDs, refID)
		}
	}
	return allRevIDs
}

func extractReviewerID(fullID string) string {
	return strings.Split(strings.Split(strings.Split(strings.Split(fullID, ".")[0], "-")[0], ":")[0], " ")[0]
}

// searchSlice is a generic linear search function that works for any slice type
func searchSlice[T comparable](slice []T, element T) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}
