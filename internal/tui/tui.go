package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/m-porter/jumper/internal/config"
	"github.com/m-porter/jumper/internal/lib"
	"github.com/m-porter/jumper/internal/logger"
	"go.uber.org/zap"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/m-porter/jumper/internal/core"
)

var selectedPath = "."

type Options struct {
	StartingQuery string
}

type searchResultsMsg struct {
	Items     []listItem
	Timestamp int64
}

type cacheUpdatedEvent struct{}

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
	return tea.Batch(
		tea.EnterAltScreen,
		tea.DisableMouse,
		searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano()),
		func() tea.Msg {
			m.App.Setup()
			return cacheUpdatedEvent{}
		},
	)
}

func (m *model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case searchResultsMsg:
		m.handleSearchResultsMsg(msg)
		return m, nil

	case cacheUpdatedEvent:
		logger.Log("cacheUpdatedEvent received")
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	case tea.WindowSizeMsg:
		size := message.(tea.WindowSizeMsg)
		m.WindowSize = &windowSize{
			Height: size.Height,
			Width:  size.Width,
		}

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

func (m *model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	case tea.KeyBackspace:
		if len(m.InputValue) > 0 {
			m.InputValue = m.InputValue[:len(m.InputValue)-1]
		}
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	case tea.KeySpace:
		m.InputValue += " "
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	case tea.KeyRunes:
		m.InputValue = fmt.Sprintf("%s%s", m.InputValue, msg.String())
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())
	}

	return m, nil
}

func (m *model) View() string {
	var output []string

	width := 0
	if m.WindowSize != nil {
		width = m.WindowSize.Width
	}

	searchBox := SearchBoxComponent(m.InputValue, len(m.ListItems), len(m.App.Directories), width)
	output = append(output, searchBox)

	// only print stuff if we know the window size or rendering gets messed up
	if m.WindowSize != nil {
		for i, item := range m.ListItems {
			if i < m.WindowSize.Height-2 {
				//line := m.ListStyle.format(item, m.CursorPos == i)
				line := ProjectRowComponent(item.Path, m.CursorPos == i)
				output = append(output, line)
			}
		}
	}

	return strings.Join(output, "\n")
}

func (m *model) handleSearchResultsMsg(msg searchResultsMsg) {
	if msg.Timestamp > m.ListLastUpdatedAt {
		logger.Log("searchResultsMsg received", zap.Int64("last_ts", m.ListLastUpdatedAt), zap.Int64("ts", msg.Timestamp))
		m.ListLastUpdatedAt = msg.Timestamp
		m.ListItems = msg.Items
		m.CursorPos = 0
	} else {
		logger.Log("out of order searchResultsMsg received", zap.Int64("last_ts", m.ListLastUpdatedAt), zap.Int64("ts", msg.Timestamp))
	}
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

func (m *model) toggleListStyle() {
	next := int(m.ListStyle) + 1
	if next < len(listStyles) {
		m.ListStyle = listStyles[next]
	} else {
		m.ListStyle = listStyles[0]
	}
}

func searchCmd(dirs []string, term string, now int64) tea.Cmd {
	return func() tea.Msg {
		var results []string
		if term == "" {
			results = dirs
		} else {
			results = lib.FuzzySearchSlice(dirs, term)
		}

		return searchResultsMsg{
			Items:     pathsToListItems(results),
			Timestamp: now,
		}
	}
}

func Run(opts Options) (string, error) {
	initIndicators(config.Get().LineIndicator)

	app := core.NewApp()

	m := &model{
		App:        app,
		InputValue: opts.StartingQuery,
	}

	program := tea.NewProgram(m, tea.WithAltScreen())

	app.SetCacheUpdateCallback(func() {
		program.Send(cacheUpdatedEvent{})
	})

	_, err := program.Run()
	return selectedPath, err
}
