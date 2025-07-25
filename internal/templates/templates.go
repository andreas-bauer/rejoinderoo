package templates

import (
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates/latex"
	"github.com/andreas-bauer/rejoinderoo/internal/templates/typst"
)

// Template is an interface for templates.
type Template interface {
	Render(td reader.TabularData) (string, error)
	FileExtension() string
}

// Available returns a list of available template names.
func Available() []string {
	return []string{
		"LaTeX",
		"Typst",
	}
}

// NewTemplate creates a new template based on the specified type.
// NewTemplate defaults to returning a LaTeX template, if a given name is not recognized.
func NewTemplate(name string) Template {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "typst":
		return typst.NewTypstTemplate()
	default:
		return latex.NewLatexTemplate()
	}
}
