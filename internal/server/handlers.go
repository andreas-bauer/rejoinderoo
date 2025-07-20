package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
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

	file, handler, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		http.Error(w, "No file submitted", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if !isAllowedContentType(handler.Header.Get("Content-Type")) {
		http.Error(w, "File type is not supported", http.StatusInternalServerError)
		return
	}

	csv := reader.CSVReader{}
	tb, err := csv.Read(file)
	if err != nil {
		fmt.Println("Error reading CSV")
	}
	fmt.Println(tb.Headers)

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

	file, handler, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		http.Error(w, "No file submitted", http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}

	if !isAllowedContentType(handler.Header.Get("Content-Type")) {
		http.Error(w, "File type is not supported", http.StatusInternalServerError)
		return
	}

	fmt.Println("Generate file:")
	fmt.Println(handler.Filename)

	selectedHeaders := []string{}
	for key, values := range r.Form {
		if strings.HasPrefix(key, "header-") {
			name := strings.TrimPrefix(key, "header-")
			selectedHeaders = append(selectedHeaders, name)
		}
		for _, value := range values {
			fmt.Printf("Form key: %s, value: %s\n", key, value)
		}
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

	var tmpl templates.Template = latex.NewLatexTemplate()
	out, err := tmpl.Render(*td)
	if err != nil {
		fmt.Println("Error rendering template:", err)
	}

	// Set headers to prompt file download
	w.Header().Set("Content-Disposition", "attachment; filename=\"output.txt\"")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(out)))
	io.Copy(w, strings.NewReader(out))
}

func isAllowedContentType(ct string) bool {
	return ct == "text/csv" || ct == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}
