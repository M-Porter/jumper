package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/rivo/tview"
	"time"
)

type listStyle int

var (
	searchVal = ""

	done               = make(chan struct{})
	tickerTimeInterval = time.Millisecond * 10

	cursorPos       = 0
	resultsListMaxH = 0

	listStyleShort    listStyle = 0
	listStyleLong     listStyle = 1
	selectedListStyle           = listStyleShort
)

func colorize(view *tview.TextView, text string) {
	w := tview.ANSIWriter(view)
	_, _ = w.Write([]byte(text))
}

func tui() error {
	app := tview.NewApplication()

	resultsView := tview.NewFlex()
	go resultsViewUpdater(app, resultsView)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(inputView(), 1, 1, true)
	flex.AddItem(resultsView, 0, 1, false)

	app.SetInputCapture(tuiKeyCapture)

	return app.SetRoot(flex, true).EnableMouse(false).Run()
}

func tuiKeyCapture(event *tcell.EventKey) *tcell.EventKey {
	// tab to flip between list styles
	if event.Key() == tcell.KeyTab {
		selectedListStyle = 1 - selectedListStyle
	}

	// exit out
	if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape {
		fmt.Print(".")
		close(done)
	}

	// print out selected row on enter press
	if event.Key() == tcell.KeyEnter {
		fmt.Print(rt.Directories[cursorPos].Path)
		close(done)
	}

	// move cursor around
	if event.Key() == tcell.KeyUp {
		moveCursorPosUp()
	}
	if event.Key() == tcell.KeyDown {
		moveCursorPosDown()
	}

	return event
}

func moveCursorPosUp() {
	if cursorPos <= 0 {
		cursorPos = 0
	} else {
		cursorPos--
	}
}

func moveCursorPosDown() {
	dirCount := len(rt.Directories) - 1
	if cursorPos >= dirCount {
		cursorPos = dirCount
	} else {
		if cursorPos >= resultsListMaxH {
			cursorPos = resultsListMaxH - 1
		} else {
			cursorPos++
		}
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

func resultsViewUpdater(app *tview.Application, view *tview.Flex) {
	view.SetDirection(tview.FlexRow)
	view.SetBackgroundColor(tcell.ColorReset)

	ticker := time.NewTicker(tickerTimeInterval)
	for {
		select {
		case <-done:
			ticker.Stop()
			app.Stop()
			return
		case <-ticker.C:
			_, _, _, resultsListMaxH = view.GetInnerRect()
			app.QueueUpdateDraw(func() {
				view.Clear()
				addResults(view)
			})
		}
	}
}

func addResults(view *tview.Flex) {
	//results := filterDirectories(rt.Directories, searchVal)

	for i, dir := range rt.Directories {
		//for i, result := range results {
		line := tview.NewTextView()
		line.SetBackgroundColor(tcell.ColorReset)
		line.SetTextColor(tcell.ColorReset)
		line.SetDynamicColors(true)

		label := dir.Label
		if selectedListStyle == listStyleLong {
			label = dir.Path
		}
		//label := result.Str

		space := " "
		if i == cursorPos {
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
