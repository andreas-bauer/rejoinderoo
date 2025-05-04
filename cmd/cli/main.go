package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andreas-bauer/rejoinderoo/internal/reader"
	"github.com/andreas-bauer/rejoinderoo/internal/templates"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))                                 // Mocha mauve
	helpStyle         = lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1).Foreground(lipgloss.Color("#7f849c")) // Mocha Overlay 1
)

type model struct {
	cursor  int
	Choices []item
}

type item struct {
	text    string
	checked bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case " ":
			// Toggle the checked state of the current item.
			m.Choices[m.cursor].checked = !m.Choices[m.cursor].checked

		case "enter":
			// Send the choice on the channel and exit.
			return m, tea.Quit

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.Choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.Choices) - 1
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("Which columns should be considered?\n")
	s.WriteString("Select at least three fields (ID, reviewer comment, author response)\n\n")

	for i, item := range m.Choices {
		if m.cursor == i {
			s.WriteString(">")

		} else {
			s.WriteString(" ")
		}
		if item.checked {
			s.WriteString(selectedItemStyle.Render("(•) " + item.text))
		} else {
			s.WriteString("( ) " + item.text)
		}
		s.WriteString("\n")
	}
	s.WriteString(helpStyle.Render("\n(press space to select item)\n(press q to quit)\n"))

	return s.String()
}

func asItems(headers []string) []item {
	items := make([]item, len(headers))
	for i, header := range headers {
		items[i] = item{
			text:    header,
			checked: false,
		}
	}
	return items
}

func selectColumns(headers []string) []string {
	tm := model{
		Choices: asItems(headers),
	}

	p := tea.NewProgram(tm)

	// Run returns the model as a tea.Model.
	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.
	m, ok := m.(model)
	if !ok {
		return []string{}
	}

	selected := make([]string, 0, len(m.(model).Choices))
	for _, item := range m.(model).Choices {
		if item.checked {
			selected = append(selected, item.text)
		}
	}

	return selected
}

func main() {
	inFile := flag.String("i", "./tests/small.csv", "file path to input file (CSV or Excel)")
	outFile := flag.String("o", "out.tex", "file path to output LaTeX file")
	flag.Parse()

	fmt.Println("Rejoinderoo")
	fmt.Println("----------------------")
	fmt.Printf("input: %s\n", *inFile)
	fmt.Printf("output: %s\n", *outFile)
	fmt.Println("----------------------")

	td, err := reader.NewReader(*inFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		fmt.Println("Unable to proceed.")

		os.Exit(1)
	}

	selected := selectColumns(td.Headers)

	if len(selected) < 3 {
		fmt.Println("Unable to proceed. Requires at least three columns (ID, reviewer comment, author response) to proceed.")
		os.Exit(1)
	}

	msg := fmt.Sprintf("You selected %d/%d columns:", len(selected), cap(selected))
	fmt.Println(msg)
	for _, col := range selected {
		fmt.Printf("  - %s\n", col)
	}
	td.Keep(selected)
	templates.TemplateTest(td)

	fmt.Println("\n\n")

	fmt.Println("⭐️ If you enjoy this project, please consider giving it a star on GitHub")
	fmt.Println("└─ https://github.com/andreas-bauer/rejoinderoo")
}
