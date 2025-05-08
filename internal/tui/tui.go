package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7"))                                 // Mocha mauve
	helpStyle         = lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1).Foreground(lipgloss.Color("#7f849c")) // Mocha Overlay 1
)

type Model struct {
	cursor  int
	Choices []item
}

type item struct {
	text    string
	checked bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewModel(headers []string) (*Model, error) {
	return &Model{
		Choices: asItems(headers),
	}, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
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

func (m Model) GetSelected() []string {
	selected := make([]string, 0, len(m.Choices))
	for _, item := range m.Choices {
		if item.checked {
			selected = append(selected, item.text)
		}
	}
	return selected
}
