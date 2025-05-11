package tui

import (
	"fmt"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/templates"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type FormResult struct {
	Filename        string
	SelectedColumns []string
	Template        templates.TemplateType
}

func RunForm(allHeaders []string) (FormResult, error) {
	var headerOpts = make([]huh.Option[string], len(allHeaders))
	for i, header := range allHeaders {
		headerOpts[i] = huh.NewOption(header, header)
	}

	var result FormResult

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("Select Columns").
				Description("Select the columns you want to include in the rejoinder").
				Options(headerOpts...).
				Validate(func(t []string) error {
					if len(t) < 3 {
						return fmt.Errorf("at least three columns need to be selected")
					}
					return nil

				}).
				Value(&result.SelectedColumns),
		),
		huh.NewGroup(
			huh.NewSelect[templates.TemplateType]().Title("Template").
				Description("Select the output template for the rejoinder").
				Options(
					huh.NewOption("LaTeX", templates.LatexTemplate).Selected(true),
					huh.NewOption("Typst", templates.TypstTemplate),
				).
				Value(&result.Template),
		),
		huh.NewGroup(
			huh.NewInput().Title("Filename").
				Description("The file name of the generated rejoinder").
				Prompt("> ").
				Placeholder("output.tex").
				Value(&result.Filename),
		),
	)

	err := form.Run()
	if err != nil {
		return result, err
	}

	return result, nil

}

func PrintSummary(result FormResult) {
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n\nTempate: %s\nFilename: %s\n\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Render("✅ Rejoinder created"),
		keyword(string(result.Template)),
		keyword(result.Filename),
		"⭐️ If you enjoy this project, please consider giving it a star on GitHub:",
		keyword("   https://github.com/andreas-bauer/rejoinderoo"),
	)
	fmt.Println(
		lipgloss.NewStyle().
			Width(80).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)
}
