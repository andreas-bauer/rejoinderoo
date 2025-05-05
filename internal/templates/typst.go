package templates

import "strings"

// Typst handles escaping special characters for Typst templates.
type Typst struct{}

func NewTypstTemplate() *Typst {
	return &Typst{}
}

func (t *Typst) File() string {
	return "./internal/templates/typst.tmpl"
}

// Escape escapes special characters for Typst.
func (t *Typst) Escape(input string) string {
	// Replace Typst special characters with their escaped versions.
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"#", "\\#",
	)
	return replacer.Replace(input)
}
