package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/rivo/tview"
	"io"
	"time"
)

type TUIState struct {
	SearchVal       string
	CursorPos       int
	ListStyle       listStyle
	ResultsListMaxH int
}

type TUI struct {
	App    *Application
	Screen *tview.Application
	Ticker *time.Ticker
	Done   chan struct{}
	State  *TUIState
}

func NewTUI(app *Application) *TUI {
	tui := &TUI{
		App:    app,
		Screen: tview.NewApplication(),

		Ticker: time.NewTicker(tickerTimeInterval),
		Done:   make(chan struct{}),

		State: &TUIState{
			SearchVal:       "",
			CursorPos:       0,
			ListStyle:       listStyleShort,
			ResultsListMaxH: 0,
		},
	}

	defer tui.Screen.Stop()

	return tui
}

func (t *TUI) Run() error {
	resultsView := tview.NewFlex()
	go t.resultsViewUpdater(resultsView)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(t.inputView(), 1, 1, true)
	flex.AddItem(resultsView, 0, 1, false)

	t.Screen.SetInputCapture(t.tuiKeyCapture)

	// see https://github.com/rivo/tview/issues/270#issuecomment-485083503
	t.Screen.SetBeforeDrawFunc(t.beforeDrawFunc)

	return t.Screen.SetRoot(flex, true).EnableMouse(false).Run()
}

func (t *TUI) Stop() {
	t.Ticker.Stop()
	t.Screen.Stop()
}

func (*TUI) beforeDrawFunc(screen tcell.Screen) bool {
	screen.Clear()
	return false
}

func (t *TUI) toggleListStyle() {
	t.State.ListStyle = 1 - t.State.ListStyle
}

func (t *TUI) tuiKeyCapture(event *tcell.EventKey) *tcell.EventKey {
	// tab to flip between list styles
	if event.Key() == tcell.KeyTab {
		t.toggleListStyle()
	}

	// exit out
	if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape {
		fmt.Print(".")
		close(t.Done)
	}

	// print out selected row on enter press
	if event.Key() == tcell.KeyEnter {
		fmt.Print(t.App.Directories[t.State.CursorPos])
		close(t.Done)
	}

	// move cursor around
	if event.Key() == tcell.KeyUp {
		t.moveCursorPosUp()
	}
	if event.Key() == tcell.KeyDown {
		t.moveCursorPosDown()
	}

	return event
}

func (t *TUI) moveCursorPosUp() {
	if t.State.CursorPos <= 0 {
		t.State.CursorPos = 0
	} else {
		t.State.CursorPos--
	}
}

func (t *TUI) moveCursorPosDown() {
	dirCount := len(t.App.Directories) - 1
	if t.State.CursorPos >= dirCount {
		t.State.CursorPos = dirCount
	} else {
		if t.State.CursorPos >= t.State.ResultsListMaxH {
			t.State.CursorPos = t.State.ResultsListMaxH - 1
		} else {
			t.State.CursorPos++
		}
	}
}

func (t *TUI) resultsViewUpdater(view *tview.Flex) {
	view.SetDirection(tview.FlexRow)
	view.SetBackgroundColor(tcell.ColorReset)

	for {
		select {
		case <-t.Done:
			t.Stop()
			return
		case <-t.Ticker.C:
			_, _, _, height := view.GetInnerRect()
			t.State.ResultsListMaxH = height

			t.Screen.QueueUpdateDraw(func() {
				view.Clear()
				t.addResults(view)
			})
		}
	}
}

func (t *TUI) addResults(view *tview.Flex) {
	results := filterDirectories(t.App.Directories, t.State.SearchVal)

	//for i, dir := range app.Directories {
	for i, result := range results {
		line := tview.NewTextView()
		line.SetBackgroundColor(tcell.ColorReset)
		line.SetTextColor(tcell.ColorReset)
		line.SetDynamicColors(true)

		//label := dir.Label
		//if ListStyle == listStyleLong {
		//	label = dir.Path
		//}
		label := result.Str

		space := " "
		if i == t.State.CursorPos {
			space = ">"
			label = color.HEX("#424242", true).Sprintf(" %s ", label)
		} else {
			label = fmt.Sprintf(" %s", label)
		}

		space = color.HEX("#424242", true).
			Sprintf("%s", color.HEX("#E53935").Sprint(space))

		colorize(
			line,
			color.BgDefault.Sprintf("%s%s", space, label),
		)

		view.AddItem(line, 1, 1, false)
	}
}

func (t *TUI) inputView() *tview.InputField {
	in := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorReset).
		SetLabelColor(tcell.ColorBlue).
		SetChangedFunc(func(text string) {
			t.State.SearchVal = text
		})

	in.SetBackgroundColor(tcell.ColorReset)
	in.SetFieldTextColor(tcell.ColorReset)

	return in
}

func colorize(v io.Writer, text string) {
	_, _ = tview.ANSIWriter(v).Write([]byte(text))
}
