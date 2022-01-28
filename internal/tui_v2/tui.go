package tui_v2

import (
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/m-porter/jumper/internal/core"
)

const lineIndicator = "❯"

var program *tea.Program

type listItem struct {
	Path string
	Base string
	Dir  string
}

type model struct {
	App               *core.Application
	CursorPos         int
	ListStyle         listStyle
	ListItems         []listItem
	ListLastUpdatedAt int64
	InputValue        string
}

func pathsToListItems(paths []string) []listItem {
	var r []listItem
	for _, path := range paths {
		r = append(r, listItem{
			Path: path,
			Base: filepath.Base(path),
			Dir:  filepath.Dir(path),
		})
	}
	return r
}

func (m *model) Init() tea.Cmd {
	go m.App.Setup()
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
			go m.Search()
		}
	}

	return m, nil
}

func (m *model) View() string {
	var output string

	{
		inputArrowStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
		inputLine := fmt.Sprintf("%s %s", inputArrowStyle.Render(lineIndicator), m.InputValue)
		output += inputLine
		output += "\n"
	}

	for _, item := range m.ListItems {
		vBarStyle := lipgloss.NewStyle().Background(lipgloss.Color("#343433"))
		line := fmt.Sprintf("%s %s", vBarStyle.Render(" "), item.Base)
		output += line
		output += "\n"
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

func (m *model) Search() {
	m.ListItems = pathsToListItems(m.App.Directories)
	program.Send(nil)
}

func Run(debug bool) error {
	m := &model{
		App: core.NewApp(debug),
	}

	program = tea.NewProgram(m, tea.WithAltScreen())

	return program.Start()
}
