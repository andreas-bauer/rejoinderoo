package server

import (
	"html/template"
	"net/http"
)

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
