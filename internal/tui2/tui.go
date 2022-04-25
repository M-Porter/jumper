package tui2

import (
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/m-porter/jumper/internal/core"
	"github.com/m-porter/jumper/internal/lib"
	"github.com/rivo/tview"
)

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

type state struct {
	CursorPos         int
	ListStyle         listStyle
	ListItems         []listItem
	ListLastUpdatedAt int64
	InputValue        string
	WindowHeight      int
}

func (s *state) cycleListStyle() {
	next := int(s.ListStyle) + 1
	if next < len(listStyles) {
		s.ListStyle = listStyles[next]
	} else {
		s.ListStyle = listStyles[0]
	}
}

func (m *state) moveCursorDown() {
	listLen := len(m.ListItems) - 1
	if m.CursorPos >= listLen {
		m.CursorPos = listLen
	} else if m.CursorPos >= m.WindowHeight {
		m.CursorPos = m.WindowHeight - 1
	} else {
		m.CursorPos++
	}
}

func (m *state) moveCursorUp() {
	if m.CursorPos <= 0 {
		m.CursorPos = 0
	} else {
		m.CursorPos--
	}
}

func (m *state) resetCursorPos() {
	m.CursorPos = 0
}

type TUI struct {
	App    *core.Application
	Screen *tview.Application
	Events lib.Events
	State  *state
}

var selectedPath string

type listItem struct {
	Path string
	Base string
	Dir  string
}

func Run(debug bool, startingQuery string) (string, error) {
	m := &state{
		ListStyle:  listStyleShort,
		InputValue: startingQuery,
	}

	tui := &TUI{
		App:    core.NewApp(debug),
		Screen: tview.NewApplication(),
		Events: lib.EventsStream(),
		State:  m,
	}

	err := tui.Start()
	return selectedPath, err
}

func (t *TUI) Start() error {
	// from https://github.com/M-Porter/jumper/blob/d00abeed173a91e88b5e0d071a0d7bc9c47e6bcb/internal/tui/tui.go

	go func() {
		t.App.Setup()
		t.doSearch("")
	}()

	resultsView := tview.NewFlex()
	go t.resultsViewUpdater(resultsView)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(t.inputView(), 1, 1, true)
	flex.AddItem(resultsView, 0, 1, false)

	t.Screen.SetInputCapture(t.inputCaptureFunc)

	// see https://github.com/rivo/tview/issues/270#issuecomment-485083503
	t.Screen.SetBeforeDrawFunc(t.beforeDrawFunc)

	defer t.Screen.Stop()
	return t.Screen.SetRoot(flex, true).EnableMouse(false).Run()
}

func (t *TUI) beforeDrawFunc(screen tcell.Screen) bool {
	screen.Clear()
	return false
}

func (t *TUI) inputCaptureFunc(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape {
		t.Events.Done()
		selectedPath = "."
	}

	if event.Key() == tcell.KeyEnter {
		t.Events.Done()
		selectedPath = t.State.ListItems[t.State.CursorPos].Path
	}

	if event.Key() == tcell.KeyTab {
		t.State.cycleListStyle()
		t.Events.Update()
	}

	if event.Key() == tcell.KeyDown {
		t.State.moveCursorDown()
		t.Events.Update()
	}

	if event.Key() == tcell.KeyUp {
		t.State.moveCursorUp()
		t.Events.Update()
	}

	return event
}

func (t *TUI) inputView() tview.Primitive {
	in := tview.NewInputField().
		SetLabel(lineIndicator + " ").
		SetFieldBackgroundColor(tcell.ColorReset).
		SetLabelColor(tcell.GetColor("#0EA5E9")).
		SetChangedFunc(func(text string) {
			t.State.resetCursorPos()
			go t.doSearch(text)
		})

	in.SetBackgroundColor(tcell.ColorReset)
	in.SetFieldTextColor(tcell.ColorReset)

	return in
}

func (t *TUI) doSearch(text string) {
	var results []string

	if text == "" {
		results = t.App.Directories
	} else {
		results = lib.FuzzySearchSlice(t.App.Directories, text)
	}

	now := time.Now().UnixNano()
	if now > t.State.ListLastUpdatedAt {
		t.State.ListItems = pathsToListItems(results)
		t.State.ListLastUpdatedAt = now
		t.Events.Update()
	}
}

func (t *TUI) resultsViewUpdater(view *tview.Flex) {
	view.SetDirection(tview.FlexRow)
	view.SetBackgroundColor(tcell.ColorReset)

	for {
		select {
		case evt := <-t.Events:
			switch evt {
			case lib.EventUpdate:
				_, _, _, height := view.GetInnerRect()
				t.State.WindowHeight = height
				t.Screen.QueueUpdateDraw(func() {
					view.Clear()
					t.addResults(view)
				})
			case lib.EventDone:
				t.Events.Close()
				t.Screen.Stop()
				return
			}
		}
	}
}

func (t *TUI) addResults(view *tview.Flex) {
	if t.State.ListLastUpdatedAt == 0 {
		t.State.ListItems = pathsToListItems(t.App.Directories)
		t.State.ListLastUpdatedAt = time.Now().UnixNano()
	}

	for i, item := range t.State.ListItems {
		line := tview.NewTextView()
		line.SetBackgroundColor(tcell.ColorReset)
		line.SetTextColor(tcell.ColorReset)
		line.SetDynamicColors(true)
		label := t.State.ListStyle.format(item, i == t.State.CursorPos)
		_, _ = tview.ANSIWriter(line).Write([]byte(label))
		view.AddItem(line, 1, 1, false)
	}
}
