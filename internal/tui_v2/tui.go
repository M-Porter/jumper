package tui_v2

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const lineIndicator = "❯"

type listItem struct {
	Path string
	Base string
	Dir  string
}

type model struct {
	CursorPos         int
	ListStyle         listStyle
	ListItems         []listItem
	ListLastUpdatedAt int64
	InputValue        string
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.DisableMouse)
}

func (m *model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp:
			m.MoveCursorUp()

		case tea.KeyDown:
			m.MoveCursorDown()

		case tea.KeyEnter:
			// todo

		case tea.KeyBackspace:
			if len(m.InputValue) > 0 {
				m.InputValue = m.InputValue[:len(m.InputValue)-1]
			}

		case tea.KeyRunes:
			m.InputValue = fmt.Sprintf("%s%s", m.InputValue, msg.String())
		}
	}

	return m, nil
}

func (m *model) View() string {
	var output string

	// output the input value line
	{
		inputArrowStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))
		inputLine := fmt.Sprintf("%s %s", inputArrowStyle.Render(lineIndicator), m.InputValue)
		output += inputLine
	}

	return output
}

// MoveCursorUp decrements the cursor pos value
func (m *model) MoveCursorUp() {
	if m.CursorPos <= 0 {
		m.CursorPos = 0
	} else {
		m.CursorPos--
	}
}

// MoveCursorDown increments the cursor pos value
func (m *model) MoveCursorDown() {
	listLen := len(m.ListItems) - 1

	if m.CursorPos >= listLen {
		m.CursorPos = listLen
	} else {
		//if m.CursorPos >= m.ResultsListMaxH {
		//	m.CursorPos = m.ResultsListMaxH - 1
		//} else {
		//	m.CursorPos++
		//}
	}
}

func Run() error {
	m := &model{}

	p := tea.NewProgram(m, tea.WithAltScreen())

	return p.Start()
}
