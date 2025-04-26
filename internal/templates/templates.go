package templates

import (
	"fmt"
	"os"
	"text/template"
)

type Response struct {
	Comment  string
	Response string
	Action   string
	Where    string
	More     []string
}

type Document struct {
	Response []Response
	Style    string
	Name     string
}

func TemplateTest() {
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

	doc := Document{
		Response: data,
		Style:    "green",
		Name:     "Andi",
	}

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
