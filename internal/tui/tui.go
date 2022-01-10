package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/m-porter/jumper/internal/core"
	"github.com/m-porter/jumper/internal/lib"
	"github.com/m-porter/jumper/internal/logger"
	"github.com/rivo/tview"
	"go.uber.org/zap"
)

type State struct {
	CursorPos         int
	ListStyle         ListStyle
	ResultsListMaxH   int
	ListItems         []ListItem
	ListLastUpdatedAt int64
}

type TUI struct {
	App    *core.Application
	Screen *tview.Application
	Events lib.Events
	State  *State
}

func New(app *core.Application) *TUI {
	return &TUI{
		App:    app,
		Screen: tview.NewApplication(),
		Events: lib.EventsStream(),
		State: &State{
			CursorPos:         0,
			ListStyle:         ListStyleShort,
			ResultsListMaxH:   0,
			ListItems:         []ListItem{},
			ListLastUpdatedAt: 0,
		},
	}
}

func (t *TUI) Run() error {
	go t.Setup()

	resultsView := tview.NewFlex()
	go t.resultsViewUpdater(resultsView)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(t.inputView(), 1, 1, true)
	flex.AddItem(resultsView, 0, 1, false)

	t.Screen.SetInputCapture(t.tuiKeyCapture)

	// see https://github.com/rivo/tview/issues/270#issuecomment-485083503
	t.Screen.SetBeforeDrawFunc(t.beforeDrawFunc)

	defer t.Screen.Stop()
	return t.Screen.SetRoot(flex, true).EnableMouse(false).Run()
}

func (t *TUI) Setup() {
	t.App.Setup()  // setup the app
	t.doSearch("") // fetch initial results from cache
}

func (t *TUI) Stop() {
	t.Events.Close()
	t.Screen.Stop()
}

func (t *TUI) ExitWithNoChange() {
	t.Events.Done()
	fmt.Print(".")
}

func (t *TUI) ExitWithSelected() {
	t.Events.Done()
	fmt.Print(t.State.ListItems[t.State.CursorPos].Path)
}

func (*TUI) beforeDrawFunc(screen tcell.Screen) bool {
	screen.Clear()
	return false
}

func (t *TUI) toggleListStyle() {
	next := int(t.State.ListStyle) + 1
	if next < len(ListStyles) {
		t.State.ListStyle = ListStyles[next]
	} else {
		t.State.ListStyle = ListStyles[0]
	}

	t.Events.Update()
}

func (t *TUI) tuiKeyCapture(event *tcell.EventKey) *tcell.EventKey {
	// tab to flip between list styles
	if event.Key() == tcell.KeyTab {
		t.toggleListStyle()
	}

	// exit out
	if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyEscape {
		t.ExitWithNoChange()
	}

	// print out selected row on enter press
	if event.Key() == tcell.KeyEnter {
		t.ExitWithSelected()
	}

	// move cursor around
	if event.Key() == tcell.KeyUp {
		t.moveCursorPosUp()
		t.Events.Update()
	}
	if event.Key() == tcell.KeyDown {
		t.moveCursorPosDown()
		t.Events.Update()
	}

	return event
}

func (t *TUI) resetCursorPos() {
	t.State.CursorPos = 0
}

func (t *TUI) moveCursorPosUp() {
	if t.State.CursorPos <= 0 {
		t.State.CursorPos = 0
	} else {
		t.State.CursorPos--
	}
}

func (t *TUI) moveCursorPosDown() {
	listLen := len(t.State.ListItems) - 1
	dirCount := len(t.App.Directories) - 1

	if t.State.CursorPos >= listLen {
		t.State.CursorPos = listLen
	} else if t.State.CursorPos >= dirCount {
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
		case evt := <-t.Events:
			logger.Log("event received", zap.Int("event", evt))

			switch evt {
			case lib.EventUpdate:
				_, _, _, height := view.GetInnerRect()
				t.State.ResultsListMaxH = height

				t.Screen.QueueUpdateDraw(func() {
					view.Clear()
					t.addResults(view)
				})

			case lib.EventDone:
				t.Stop()
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

		label := item.LabelForStyle(t.State.ListStyle)

		space := " "
		if i == t.State.CursorPos {
			space = ">"
			label = ColorBgGray.Sprintf(" %s ", label)
		} else {
			label = fmt.Sprintf(" %s", label)
		}

		space = ColorBgGray.Sprintf("%s", ColorFgRed.Sprint(space))
		colorize(line, ColorBgDefault.Sprintf("%s%s", space, label))

		view.AddItem(line, 1, 1, false)
	}
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

func (t *TUI) inputView() *tview.InputField {
	in := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorReset).
		SetLabelColor(ctoc(ColorFgBlue)).
		SetChangedFunc(func(text string) {
			t.resetCursorPos()
			go t.doSearch(text)
		})

	in.SetBackgroundColor(tcell.ColorReset)
	in.SetFieldTextColor(tcell.ColorReset)

	return in
}
