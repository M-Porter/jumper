package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/rivo/tview"
	"time"
)

type TUI struct {
	App       *Application
	screen    *tview.Application
	cursorPos int
	ticker    *time.Ticker
	done      chan struct{}
}

func NewTUI(app *Application) *TUI {
	tui := &TUI{
		App:       app,
		screen:    tview.NewApplication(),
		cursorPos: 0,
		ticker:    time.NewTicker(tickerTimeInterval),
		done:      make(chan struct{}),
	}

	defer tui.screen.Stop()

	return tui
}

func (t *TUI) Run() error {
	resultsView := tview.NewFlex()
	go t.resultsViewUpdater(resultsView)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(inputView(), 1, 1, true)
	flex.AddItem(resultsView, 0, 1, false)

	t.screen.SetInputCapture(t.tuiKeyCapture)

	// see https://github.com/rivo/tview/issues/270#issuecomment-485083503
	t.screen.SetBeforeDrawFunc(t.beforeDrawFunc)

	return t.screen.SetRoot(flex, true).EnableMouse(false).Run()
}

func (t *TUI) Stop() {
	t.ticker.Stop()
	t.screen.Stop()
}

func (*TUI) beforeDrawFunc(screen tcell.Screen) bool {
	screen.Clear()
	return false
}

func (t *TUI) tuiKeyCapture(event *tcell.EventKey) *tcell.EventKey {
	// tab to flip between list styles
	if event.Key() == tcell.KeyTab {
		selectedListStyle = 1 - selectedListStyle
	}

	// exit out
	if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape {
		fmt.Print(".")
		close(t.done)
	}

	// print out selected row on enter press
	if event.Key() == tcell.KeyEnter {
		fmt.Print(app.Directories[t.cursorPos])
		close(t.done)
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
	if t.cursorPos <= 0 {
		t.cursorPos = 0
	} else {
		t.cursorPos--
	}
}

func (t *TUI) moveCursorPosDown() {
	dirCount := len(t.App.Directories) - 1
	if t.cursorPos >= dirCount {
		t.cursorPos = dirCount
	} else {
		if t.cursorPos >= resultsListMaxH {
			t.cursorPos = resultsListMaxH - 1
		} else {
			t.cursorPos++
		}
	}
}

func (t *TUI) resultsViewUpdater(view *tview.Flex) {
	view.SetDirection(tview.FlexRow)
	view.SetBackgroundColor(tcell.ColorReset)

	for {
		select {
		case <-t.done:
			t.Stop()
			return
		case <-t.ticker.C:
			_, _, _, resultsListMaxH = view.GetInnerRect()
			t.screen.QueueUpdateDraw(func() {
				view.Clear()
				t.addResults(view)
			})
		}
	}
}

func (t *TUI) addResults(view *tview.Flex) {
	results := filterDirectories(app.Directories, searchVal)

	//for i, dir := range app.Directories {
	for i, result := range results {
		line := tview.NewTextView()
		line.SetBackgroundColor(tcell.ColorReset)
		line.SetTextColor(tcell.ColorReset)
		line.SetDynamicColors(true)

		//label := dir.Label
		//if selectedListStyle == listStyleLong {
		//	label = dir.Path
		//}
		label := result.Str

		space := " "
		if i == t.cursorPos {
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

func inputView() *tview.InputField {
	in := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorReset).
		SetLabelColor(tcell.ColorBlue).
		SetChangedFunc(func(text string) {
			searchVal = text
		})
	in.SetBackgroundColor(tcell.ColorReset)
	in.SetFieldTextColor(tcell.ColorReset)
	return in
}

func colorize(view *tview.TextView, text string) {
	w := tview.ANSIWriter(view)
	_, _ = w.Write([]byte(text))
}
