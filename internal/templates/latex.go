package templates

import "strings"

// Latex handles escaping special characters for LaTeX templates.
type Latex struct{}

func NewLatexTemplate() *Latex {
	return &Latex{}
}

func (l *Latex) File() string {
	return "./internal/templates/latex.tmpl"
}

// Escape escapes special characters for LaTeX.
func (l *Latex) Escape(input string) string {
	// Replace LaTeX special characters with their escaped versions.
	replacer := strings.NewReplacer(
		"\\", "\\textbackslash{}",
		"{", "\\{",
		"}", "\\}",
		"$", "\\$",
		"&", "\\&",
		"#", "\\#",
		"_", "\\_",
		"%", "\\%",
		"~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}",
	)
	return replacer.Replace(input)
}
