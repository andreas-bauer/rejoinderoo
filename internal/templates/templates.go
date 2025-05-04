package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
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

func TemplateTest(td *reader.TabularData) {

	doc := createDoc(td)

	var tmplFile = "./internal/templates/latex.tmpl"

	tmpl, err := template.ParseFiles(tmplFile)

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
			record := &Record{
				Header: h,
				Text:   rec[i],
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
