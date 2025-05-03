package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
)

type Response struct {
	Comment  string
	Response string
	Action   string
	Where    string
	More     []string
}

type Header struct {
	Name string
	Idx  int
}

type Document struct {
	ReviewerIDs []string
	LenHeaders  int
	Headers     []Header
	Response    []Response
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
	data := []Response{{
		Comment:  "More data is required",
		Response: "You said not more data",
		Action:   "do nothing",
		Where:    "Section 4.1",
	}, {
		Comment:  "More data is required",
		Response: "Ok more data",
		Action:   "crawl GitHub",
		Where:    "Section 4.w",
	}}

	allRevIDs := extractReviewers(r.Records())
	headers := asDocHeaders(selected)

	return Document{
		ReviewerIDs: allRevIDs,
		LenHeaders:  len(selected),
		Headers:     headers,
		Response:    data,
	}
}

func asDocHeaders(headers []string) []Header {
	var res = make([]Header, len(headers))
	for idx, h := range headers {
		header := &Header{
			Name: h,
			Idx:  idx,
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
