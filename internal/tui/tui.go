package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/m-porter/jumper/internal/lib"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-porter/jumper/internal/core"
)

var selectedPath string

var program *tea.Program

type searchResultsUpdated struct{}

type listItem struct {
	Path string
	Base string
	Dir  string
}

type windowSize struct {
	Height int
	Width  int
}

type model struct {
	App               *core.Application
	CursorPos         int
	ListStyle         listStyle
	ListItems         []listItem
	ListLastUpdatedAt int64
	InputValue        string
	WindowSize        *windowSize
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
	go func() {
		m.App.Setup()
		m.search()
	}()
	return tea.Batch(tea.EnterAltScreen, tea.DisableMouse)
}

func (m *model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		size := message.(tea.WindowSizeMsg)
		m.WindowSize = &windowSize{
			Height: size.Height,
			Width:  size.Width,
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp:
			m.moveCursorUp()

		case tea.KeyDown:
			m.moveCursorDown()

		case tea.KeyEnter:
			selectedPath = m.ListItems[m.CursorPos].Path
			return m, tea.Quit

		case tea.KeyTab:
			m.toggleListStyle()

		case tea.KeyDelete, tea.KeyCtrlH:
			m.InputValue = ""
			go m.search()

		case tea.KeyBackspace:
			if len(m.InputValue) > 0 {
				m.InputValue = m.InputValue[:len(m.InputValue)-1]
			}
			go m.search()

		case tea.KeyRunes:
			m.InputValue = fmt.Sprintf("%s%s", m.InputValue, msg.String())
			go m.search()
		}
	}

	return m, nil
}

func (m *model) View() string {
	var output []string

	inputLine := fmt.Sprintf("%s %s", inputIndicatorPart, m.InputValue)
	output = append(output, inputLine)

	// only print stuff if we know the window size or rendering gets messed up
	if m.WindowSize != nil {
		for i, item := range m.ListItems {
			if i < m.WindowSize.Height-1 {
				line := m.ListStyle.format(item, m.CursorPos == i)
				output = append(output, line)
			}
		}
	}

	return strings.Join(output, "\n")
}

// moveCursorUp decrements the cursor pos value
func (m *model) moveCursorUp() {
	if m.CursorPos <= 0 {
		m.CursorPos = 0
	} else {
		m.CursorPos--
	}
}

// moveCursorDown increments the cursor pos value
func (m *model) moveCursorDown() {
	listLen := len(m.ListItems) - 1

	if m.CursorPos >= listLen {
		m.CursorPos = listLen
	} else {
		if m.CursorPos >= m.WindowSize.Height {
			m.CursorPos = m.WindowSize.Height - 1
		} else {
			m.CursorPos++
		}
	}
}

func (m *model) search() {
	var results []string

	if m.InputValue == "" {
		results = m.App.Directories
	} else {
		results = lib.FuzzySearchSlice(m.App.Directories, m.InputValue)
	}

	// prevents out-of-order updates
	now := time.Now().UnixNano()
	if now > m.ListLastUpdatedAt {
		m.ListItems = pathsToListItems(results)
		m.CursorPos = 0

		program.Send(searchResultsUpdated{})
	}
}

func (m *model) toggleListStyle() {
	next := int(m.ListStyle) + 1
	if next < len(listStyles) {
		m.ListStyle = listStyles[next]
	} else {
		m.ListStyle = listStyles[0]
	}
}

func Run(debug bool, startingQuery string) (string, error) {
	m := &model{
		App:        core.NewApp(debug),
		InputValue: startingQuery,
	}

	program = tea.NewProgram(m, tea.WithAltScreen())

	return selectedPath, program.Start()
}
