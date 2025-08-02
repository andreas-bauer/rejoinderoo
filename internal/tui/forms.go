package tui

import (
	"fmt"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/templates"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type FormData struct {
	Filename         string
	AvailableHeaders []string
	SelectedHeaders  []string
	Template         string
}

func RunFilePicker() string {
	var file string
	huh.NewForm(
		huh.NewGroup(
			huh.NewFilePicker().
				Picking(true).
				Title("Input file").
				Description("Select a .csv, .xlsx, or .xls file").
				AllowedTypes([]string{".csv", ".xlsx", ".xls"}).
				Value(&file),
		),
	).WithShowHelp(true).Run()
	return file
}

func RunForm(fd *FormData) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().Title("Select Columns").
				Description("Select the columns you want to include in the rejoinder").
				Options(huh.NewOptions(fd.AvailableHeaders...)...).
				Validate(func(t []string) error {
					if len(t) < 3 {
						return fmt.Errorf("at least three columns need to be selected")
					}
					return nil

				}).
				Value(&fd.SelectedHeaders),
		),
		huh.NewGroup(
			huh.NewSelect[string]().Title("Template").
				Description("Select the output template for the rejoinder").
				Options(huh.NewOptions(templates.Available()...)...).
				Value(&fd.Template),
		),
		huh.NewGroup(
			huh.NewInput().Title("Filename").
				Description("The file name of the generated rejoinder").
				Prompt("> ").
				Placeholder("output.tex").
				Value(&fd.Filename),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	return nil
}

func PrintSummary(fd *FormData) {
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n\nTempate: %s\nFilename: %s\n\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Render("✅ Rejoinder created"),
		keyword(fd.Template),
		keyword(fd.Filename),
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
