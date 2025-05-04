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

func TemplateTest(r reader.FileReader, selected []string) {

	doc := createDoc(r, selected)

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

func createDoc(r reader.FileReader, selected []string) Document {
	recs := []Record{{
		Header: "H1",
		Text:   "Some Text",
	}, {
		Header: "H2",
		Text:   "Other Text",
	}, {
		Header: "H3",
		Text:   "Different Text",
	}}
	data := []Response{{
		ReviewerID: "R1",
		Records:    recs,
	}, {
		ReviewerID: "R2",
		Records:    recs,
	}}

	allRevIDs := extractReviewers(r.Records())
	headers := asDocHeaders(selected)

	return Document{
		ReviewerIDs: allRevIDs,
		LenHeaders:  len(selected),
		Headers:     headers,
		Responses:   data,
	}
}

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
	return strings.Split(strings.Split(strings.Split(fullID, ".")[0], "-")[0], ":")[0]
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
