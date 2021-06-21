package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/rivo/tview"
	"sync"
	"time"
)

var (
	searchVal = ""

	done               = make(chan struct{})
	tuiWg              sync.WaitGroup
	tickerTimeInterval = time.Millisecond * 10

	mainFrameID = "mainJumperFrame"

	cursorPos = 0

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

	in := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorReset).
		SetLabelColor(tcell.ColorRed).
		SetChangedFunc(func(text string) {
			searchVal = text
		})
	in.SetBackgroundColor(tcell.ColorReset)
	in.SetFieldTextColor(tcell.ColorReset)

	resultsView := tview.NewFlex()
	go mainViewUpdater(app, resultsView)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(in, 1, 1, true)
	flex.AddItem(resultsView, 0, 1, false)

	app.SetInputCapture(tuiKeyCapture)

	return app.SetRoot(flex, true).EnableMouse(false).Run()
}

func tuiKeyCapture(event *tcell.EventKey) *tcell.EventKey {
	// tab to flip between list styles
	if event.Key() == tcell.KeyTab {
		selectedListStyle = 1 - selectedListStyle
	}

	if event.Key() == tcell.KeyCtrlC {
		fmt.Print(".")
		return event
	}

	return event
}

func mainViewUpdater(app *tview.Application, view *tview.Flex) {
	view.SetDirection(tview.FlexRow)
	view.SetBackgroundColor(tcell.ColorReset)

	ticker := time.NewTicker(tickerTimeInterval)
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			view.Clear()
			addResults(view)
			app.Draw()
		}
	}
}

func addResults(view *tview.Flex) {
	for _, dir := range rt.Directories {
		line := tview.NewTextView()
		line.SetBackgroundColor(tcell.ColorReset)
		line.SetTextColor(tcell.ColorReset)
		line.SetDynamicColors(true)

		label := dir.Label
		if selectedListStyle == listStyleLong {
			label = dir.Path
		}

		colorize(
			line,
			color.BgDefault.Sprintf("%s %s", color.BgGray.Sprint(" "), label),
		)

		view.AddItem(line, 1, 1, false)
	}
}
