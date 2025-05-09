package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define different screens
type screen int

const (
	screenColumn screen = iota
	screenFormat
	screenFilename
	screenSummary
)

// Styles
var (
	docStyle                = lipgloss.NewStyle().Margin(1, 2)
	selectedStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	helpStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle             = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	errorStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	dimCheckedStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginRight(1)
	highlightedCheckedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).MarginRight(1)
)

type itemCheckable struct {
	text    string
	checked bool
}

type Model struct {
	screen screen

	// Screen 1: Column selection
	columnChoices  []itemCheckable
	columnCursor   int
	selectedColumn string

	// Screen 2: Format selection
	formatChoices  []string
	formatCursor   int
	selectedFormat string

	// Screen 3: Filename input
	filenameInput   textinput.Model
	enteredFilename string
	inputError      string // For potential filename validation

	// Screen 4: Summary
	quitting bool // To ensure summary is shown before exit
}

func NewModel(headers []string) *Model {
	ti := textinput.New()
	ti.Placeholder = "my_document"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30
	ti.Prompt = "❯ "
	ti.PromptStyle = selectedStyle
	ti.TextStyle = selectedStyle

	return &Model{
		screen:        screenColumn,
		columnChoices: asItems(headers),
		columnCursor:  0,
		formatChoices: []string{"LaTeX", "Typst"},
		formatCursor:  0,
		filenameInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink // Start the text input blinking
}

// Custom message for delayed quit
type quitMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			switch m.screen {
			case screenColumn:
				m.selectedColumn = m.columnChoices[m.columnCursor].text
				m.screen = screenFormat
				m.columnCursor = 0 // Reset cursor for next list if needed
			case screenFormat:
				m.selectedFormat = m.formatChoices[m.formatCursor]
				m.screen = screenFilename
				m.formatCursor = 0
				m.filenameInput.Focus() // Focus the text input
				return m, textinput.Blink
			case screenFilename:
				filename := strings.TrimSpace(m.filenameInput.Value())
				if filename == "" {
					m.inputError = "Filename cannot be empty."
				} else {
					m.enteredFilename = filename
					m.screen = screenSummary
					m.inputError = ""
					// Schedule a quit message after a short delay to show summary
					return m, tea.Tick(1500*time.Millisecond, func(t time.Time) tea.Msg {
						return quitMsg{}
					})
				}
			case screenSummary:
				// Already handled by tea.Tick, but an enter here could also quit
				m.quitting = true
				return m, tea.Quit
			}
		case "up", "k":
			if m.screen == screenColumn && m.columnCursor > 0 {
				m.columnCursor--
			}
			if m.screen == screenFormat && m.formatCursor > 0 {
				m.formatCursor--
			}
		case "down", "j":
			if m.screen == screenColumn && m.columnCursor < len(m.columnChoices)-1 {
				m.columnCursor++
			}
			if m.screen == screenFormat && m.formatCursor < len(m.formatChoices)-1 {
				m.formatCursor++
			}
		case " ":
			if m.screen == screenColumn {
				m.columnChoices[m.columnCursor].checked = !m.columnChoices[m.columnCursor].checked
			}
		}
	case quitMsg: // Handle our custom quit message
		m.quitting = true
		return m, tea.Quit
	}

	// Handle text input updates for the filename screen
	if m.screen == screenFilename {
		m.filenameInput, cmd = m.filenameInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting && m.screen != screenSummary { // if quitting from other screens
		return "Exiting...\n"
	}

	s := ""
	help := helpStyle.Render("\nUse ↑/↓ or j/k to navigate, Space to select, Enter to continue, Ctrl+C or Esc to quit.")

	switch m.screen {
	case screenColumn:
		s += headerStyle.Render("1. Choose the Columns to include:") + "\n"
		for i, choice := range m.columnChoices {
			cursor := "  " // Not selected
			checked := "•" // Not checked
			itemStyle := lipgloss.NewStyle()
			checkedStyle := dimCheckedStyle
			if m.columnCursor == i {
				cursor = selectedStyle.Render("❯ ")
			}
			if m.columnChoices[i].checked {
				itemStyle = selectedStyle
				checked = "✓"
				checkedStyle = highlightedCheckedStyle
			}
			s += fmt.Sprintf("%s%s%s\n", cursor, checkedStyle.Render(checked), itemStyle.Render(choice.text))
		}
		s += help

	case screenFormat:
		s += headerStyle.Render("2. Choose an Output Format:") + "\n"
		for i, choice := range m.formatChoices {
			cursor := "  "
			itemStyle := lipgloss.NewStyle()
			if m.formatCursor == i {
				cursor = selectedStyle.Render("❯ ")
				itemStyle = selectedStyle
			}
			s += fmt.Sprintf("%s%s\n", cursor, itemStyle.Render(choice))
		}
		s += help

	case screenFilename:
		s += headerStyle.Render("3. Enter Output Filename:") + "\n"
		s += m.filenameInput.View() + "\n"
		if m.inputError != "" {
			s += errorStyle.Render(m.inputError) + "\n"
		}
		s += helpStyle.Render("\nType filename and press Enter. No extension needed (e.g., 'report').")

	case screenSummary:
		s += headerStyle.Render("4. Summary:") + "\n"
		s += fmt.Sprintf("   Selected column: %s\n", selectedStyle.Render(m.selectedColumn))
		s += fmt.Sprintf("   Selected Format:   %s\n", selectedStyle.Render(m.selectedFormat))
		s += fmt.Sprintf("   Output Filename:   %s\n", selectedStyle.Render(m.enteredFilename))
		s += "\n" + selectedStyle.Render("File was generated!") + "\n\n"
		s += helpStyle.Render("Exiting automatically...")
	}

	return docStyle.Render(s)
}

func asItems(s []string) []itemCheckable {
	items := make([]itemCheckable, len(s))
	for i, text := range s {
		items[i] = itemCheckable{
			text:    text,
			checked: false,
		}
	}
	return items
}
