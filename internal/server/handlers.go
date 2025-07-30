package server

import (
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates"
)

const (
	maxUploadSize      = 10 << 20 // 10 MB
	minSelectedColumns = 3
)

const (
	templateIndex        = "index.html"
	templateError        = "error"
	templateSelectColumn = "select-column-form"
	templateResult       = "result"
)

const (
	formFieldFile        = "file"
	formFieldGenTemplate = "gen-template"
	headerPrefix         = "header-"
)

// Handler struct for handling HTTP requests.
type Handler struct {
	tmpl *template.Template
}

// NewHandler creates a new Handler instance.
func NewHandler(html *template.Template) *Handler {
	return &Handler{
		tmpl: html,
	}
}

// Index serves the main page.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	h.tmpl.ExecuteTemplate(w, templateIndex, nil)
}

// ColSelectForm handles the file upload and displays the column selection form.
func (h *Handler) ColSelectForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "The uploaded file is too large.")
		return
	}

	file, handler, err := h.getFormFile(r)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, err.Error())
		return
	}
	defer file.Close()

	fileReader, err := reader.NewReader(handler.Filename)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Unsupported file type: "+err.Error())
		return
	}

	tableData, err := fileReader.Read(file)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Error reading file as table data: "+err.Error())
		return
	}

	tmplArgs := struct {
		Headers   []string
		Templates []string
	}{
		Headers:   tableData.Headers,
		Templates: templates.Available(),
	}

	if err := h.tmpl.ExecuteTemplate(w, templateSelectColumn, tmplArgs); err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Error rendering results: "+err.Error())
		return
	}
}

// Generate creates the output document based on the user's selection.
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "The uploaded file is too large.")
		return
	}

	file, handler, err := h.getFormFile(r)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, err.Error())
		return
	}
	defer file.Close()

	selectedHeaders := getFormValuesWithPrefix(r.Form, headerPrefix)
	if len(selectedHeaders) < minSelectedColumns {
		h.tmpl.ExecuteTemplate(w, templateError, fmt.Sprintf("Please select at least %d columns.", minSelectedColumns))
		return
	}

	fileReader, err := reader.NewReader(handler.Filename)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Unsupported file type: "+err.Error())
		return
	}

	tableData, err := fileReader.Read(file)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Error reading file as table data: "+err.Error())
		return
	}

	selectedHeaders = sortHeaders(selectedHeaders, tableData.Headers)

	tableData.Keep(selectedHeaders)

	templateName := r.FormValue(formFieldGenTemplate)
	genTmpl := templates.NewTemplate(templateName)

	out, err := genTmpl.Render(*tableData)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Error generating output: "+err.Error())
		return
	}

	doc := struct {
		Content   string
		Filename  string
		Extension string
	}{
		Content:   out,
		Filename:  fileNameWithoutExtension(handler.Filename),
		Extension: genTmpl.FileExtension(),
	}

	if err := h.tmpl.ExecuteTemplate(w, templateResult, doc); err != nil {
		h.tmpl.ExecuteTemplate(w, templateError, "Error rendering results: "+err.Error())
	}
}

func sortHeaders(selectedHeaders []string, originalOrder []string) []string {
	var ordered []string
	for _, header := range originalOrder {
		if slices.Contains(selectedHeaders, header) {
			ordered = append(ordered, header)
		}
	}
	return ordered
}

// getFormFile retrieves the uploaded file from the request, ensuring it's valid.
// The caller is responsible for closing the returned multipart.File.
func (h *Handler) getFormFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	file, handler, err := r.FormFile(formFieldFile)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, nil, errors.New("no file submitted")
		}
		return nil, nil, fmt.Errorf("error retrieving the file: %w", err)
	}

	if !isAllowedContentType(handler.Header.Get("Content-Type")) {
		return nil, nil, fmt.Errorf("file type '%s' is not supported", handler.Header.Get("Content-Type"))
	}

	return file, handler, nil
}

// getFormValuesWithPrefix extracts values from a form whose keys have a given prefix.
func getFormValuesWithPrefix(formValues url.Values, prefix string) []string {
	var values []string
	for key := range formValues {
		if strings.HasPrefix(key, prefix) {
			values = append(values, strings.TrimPrefix(key, prefix))
		}
	}
	return values
}

// isAllowedContentType checks if the content type is in the allowed list.
func isAllowedContentType(ct string) bool {
	return ct == "text/csv" ||
		ct == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" ||
		ct == "application/vnd.ms-excel"
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}
