package templates

import (
	"github.com/andreas-bauer/rejoinderoo/internal/reader"
)

// Template is an interface for templates.
type Template interface {
	Render(td reader.TabularData) (string, error)
}

// TemplateType defines the types of templates available.
type TemplateType string

const (
	LatexTemplate TemplateType = "latex"
	TypstTemplate TemplateType = "typst"
)
