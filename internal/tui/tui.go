package tui

import (
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/m-porter/jumper/internal/lib"
	"github.com/m-porter/jumper/internal/logger"
	"github.com/m-porter/jumper/internal/theme"
	"go.uber.org/zap"

	tea "charm.land/bubbletea/v2"
	"github.com/m-porter/jumper/internal/core"
)

var selectedPath = "."

type Options struct {
	StartingQuery string
}

type searchResultsMsg struct {
	Items     []string
	Timestamp int64
}

type cacheUpdatedEvent struct{}

type windowSize struct {
	Height int
	Width  int
}

type model struct {
	App               *core.Application
	CursorPos         int
	ListItems         []string
	ListLastUpdatedAt int64
	InputValue        string
	WindowSize        *windowSize
	TruncatePaths     bool
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
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

	case tea.KeyPressMsg:
		return m.handleKeyMsg(msg)
	}

	return m, nil
}

func (m *model) handleKeyMsg(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.moveCursorUp()

	case "down":
		m.moveCursorDown()

	case "enter":
		selectedPath = m.ListItems[m.CursorPos]
		return m, tea.Quit

	case "tab":
		m.togglePathTruncation()

	case "delete", "ctrl+h":
		m.InputValue = ""
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	case "backspace":
		if len(m.InputValue) > 0 {
			m.InputValue = m.InputValue[:len(m.InputValue)-1]
		}
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	case "space":
		m.InputValue += " "
		return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())

	default:
		if len(msg.Text) > 0 {
			m.InputValue += msg.Text
			return m, searchCmd(m.App.Directories, m.InputValue, time.Now().UnixNano())
		}
	}

	return m, nil
}

func (m *model) View() tea.View {
	var output []string

	width := 0
	if m.WindowSize != nil {
		width = m.WindowSize.Width
	}

	searchBox := SearchBoxComponent(m.InputValue, len(m.ListItems), len(m.App.Directories), width)
	output = append(output, searchBox)

	statusBar := StatusBarComponent(StatusBarParams{
		LeftContents: []string{
			KeyHelpComponent(theme.ArrowUp+theme.ArrowDown, "move"),
			KeyHelpComponent(theme.Enter, "select"),
		},
		RightContents: []string{
			KeyHelpComponent("esc", "exit"),
		},
	}, width)

	// keep track of the height of everything not in the main render view so we can fill the list view
	// without overflowing
	componentsHeight := lipgloss.Height(searchBox) + lipgloss.Height(statusBar)

	// only print stuff if we know the window size or rendering gets messed up
	if m.WindowSize != nil {
		listMaxHeight := m.WindowSize.Height - componentsHeight
		for i, item := range m.ListItems {
			if i < listMaxHeight {
				line := ProjectRowComponent(item, m.CursorPos == i, m.TruncatePaths)
				line = lipgloss.NewStyle().PaddingLeft(2).Render(line)
				output = append(output, line)
			}
		}
	}

	output = append(output, statusBar)

	v := tea.NewView(strings.Join(output, "\n"))
	v.AltScreen = true
	return v
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

func (m *model) togglePathTruncation() {
	m.TruncatePaths = !m.TruncatePaths
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
			Items:     results,
			Timestamp: now,
		}
	}
}

func Run(opts Options) (string, error) {
	app := core.NewApp()

	m := &model{
		App:           app,
		InputValue:    opts.StartingQuery,
		TruncatePaths: true,
	}

	program := tea.NewProgram(m)

	app.SetCacheUpdateCallback(func() {
		program.Send(cacheUpdatedEvent{})
	})

	_, err := program.Run()
	return selectedPath, err
}
