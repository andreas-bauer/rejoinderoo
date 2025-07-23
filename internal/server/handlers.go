package server

import (
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates"
	"github.com/andreas-bauer/rejoinderoo/internal/templates/latex"
)

var size10MB int64 = 10 << 20

// Handler struct for handling HTTP requests
type Handler struct {
	tmpl *template.Template
}

// NewHandler creates a new Handler instance
func NewHandler(html *template.Template) *Handler {
	return &Handler{
		tmpl: html,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	h.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(size10MB)

	file, _, err := h.assertFormFile(r)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "error", err.Error())
		return
	}
	defer file.Close()

	csv := reader.CSVReader{}
	tb, err := csv.Read(file)
	if err != nil {
		fmt.Println("Error reading CSV")
	}

	tmplArgs := struct {
		Headers   []string
		Templates []string
	}{
		Headers:   tb.Headers,
		Templates: []string{string(templates.LatexTemplate), string(templates.TypstTemplate)},
	}

	if err := h.tmpl.ExecuteTemplate(w, "select-column-form", tmplArgs); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(size10MB)

	file, filename, err := h.assertFormFile(r)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "error", err.Error())
		return
	}
	defer file.Close()

	selectedHeaders := getFormValuesWithPrefix(r.Form, "header-")

	if len(selectedHeaders) < 3 {
		h.tmpl.ExecuteTemplate(w, "error", "select at least 3 columns")
		return
	}

	csv := reader.CSVReader{}
	td, err := csv.Read(file)
	if err != nil {
		fmt.Println("Error reading CSV")
	}

	td.Keep(selectedHeaders)

	genTemplate := r.FormValue("gen-template")
	if genTemplate == "" {
		genTemplate = "latex"
	}
	genTemplateExt := "tex"

	var tmpl templates.Template = latex.NewLatexTemplate()
	out, err := tmpl.Render(*td)
	if err != nil {
		fmt.Println("Error rendering template:", err)
	}

	filenameWithoutExt := fileNameWithoutExtension(filename)

	doc := struct {
		Content   string
		Filename  string
		Extension string
	}{
		Content:   out,
		Filename:  filenameWithoutExt,
		Extension: genTemplateExt,
	}

	if err := h.tmpl.ExecuteTemplate(w, "result", doc); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) assertFormFile(r *http.Request) (multipart.File, string, error) {
	file, handler, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		return nil, handler.Filename, errors.New("no file submitted")
	}
	if err != nil {
		return nil, handler.Filename, errors.New("error retrieving the file")
	}
	defer file.Close()

	if !isAllowedContentType(handler.Header.Get("Content-Type")) {
		return nil, handler.Filename, errors.New("file type is not supported")
	}

	return file, handler.Filename, nil
}

func getFormValuesWithPrefix(formValues url.Values, prefix string) []string {
	var values []string
	for key, vals := range formValues {
		if strings.HasPrefix(key, prefix) {
			values = append(values, vals...)
		}
	}
	return values
}

func isAllowedContentType(ct string) bool {
	return ct == "text/csv" || ct == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
