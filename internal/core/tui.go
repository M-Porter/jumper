package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/rivo/tview"
	"io"
	"path/filepath"
	"time"
)

var (
	ColorBgDefault = color.BgDefault
	ColorBgGray    = color.HEX("#424242", true)
	ColorFgRed     = color.HEX("#E53935")
	ColorFgBlue    = color.HEX("#60A5FA")
)

type TUIState struct {
	CursorPos         int
	ListStyle         listStyle
	ResultsListMaxH   int
	ListItems         []ListItem
	ListLastUpdatedAt int64
}

type TUI struct {
	App    *Application
	Screen *tview.Application
	Ticker *time.Ticker
	Done   chan struct{}
	State  *TUIState
}

func NewTUI(app *Application) *TUI {
	return &TUI{
		App:    app,
		Screen: tview.NewApplication(),

		Ticker: time.NewTicker(tickerTimeInterval),
		Done:   make(chan struct{}),

		State: &TUIState{
			CursorPos:         0,
			ListStyle:         listStyleShort,
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
	t.Ticker.Stop()
	t.Screen.Stop()
}

func (t *TUI) ExitWithNoChange() {
	fmt.Print(".")
	close(t.Done)
}

func (t *TUI) ExitWithSelected() {
	fmt.Print(t.App.Directories[t.State.CursorPos])
	close(t.Done)
}

func (*TUI) beforeDrawFunc(screen tcell.Screen) bool {
	screen.Clear()
	return false
}

func (t *TUI) toggleListStyle() {
	next := int(t.State.ListStyle) + 1
	if next < len(listStyles) {
		t.State.ListStyle = listStyles[next]
	} else {
		t.State.ListStyle = listStyles[0]
	}
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
		results = filterDirectories(t.App.Directories, text)
	}

	now := time.Now().UnixNano()
	if now > t.State.ListLastUpdatedAt {
		t.State.ListItems = pathsToListItems(results)
		t.State.ListLastUpdatedAt = now
	}
}

func (t *TUI) inputView() *tview.InputField {
	in := tview.NewInputField().
		SetLabel("> ").
		SetFieldBackgroundColor(tcell.ColorReset).
		SetLabelColor(ctoc(ColorFgBlue)).
		SetChangedFunc(func(text string) {
			go t.doSearch(text)
		})

	in.SetBackgroundColor(tcell.ColorReset)
	in.SetFieldTextColor(tcell.ColorReset)

	return in
}

func colorize(v io.Writer, text string) {
	_, _ = tview.ANSIWriter(v).Write([]byte(text))
}

func ctoc(c color.RGBColor) tcell.Color {
	v := c.Values()
	return tcell.NewRGBColor(int32(v[0]), int32(v[1]), int32(v[2]))
}

func pathsToListItems(paths []string) []ListItem {
	var r []ListItem
	for _, path := range paths {
		r = append(r, ListItem{
			Path: path,
			Base: filepath.Base(path),
			Dir:  filepath.Dir(path),
		})
	}
	return r
}

type ListItem struct {
	Path string
	Base string
	Dir  string
}

func (li *ListItem) LabelForStyle(style listStyle) string {
	switch style {
	case listStyleDetailed:
		return fmt.Sprintf("%s (%s)", li.Base, li.Dir)
	case listStyleLong:
		return li.Path
	case listStyleShort:
		fallthrough
	default:
		return li.Base
	}
}
