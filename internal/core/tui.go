package core

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/logrusorgru/aurora/v3"
	"io"
	"time"
)

type listStyle int

func _tui() error {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true

	g.SetManagerFunc(tuiLayout)

	if err := addArrowKeyBindings(g); err != nil {
		return err
	}
	if err := addEnterKeyBinding(g); err != nil {
		return err
	}
	if err := addTabKeyBinding(g); err != nil {
		return err
	}
	if err := quitKeyBinding(g); err != nil {
		return err
	}

	tuiWg.Add(1)
	go tuiTicker(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	tuiWg.Wait()

	return nil
}

func tuiTicker(g *gocui.Gui) {
	defer tuiWg.Done()

	ticker := time.NewTicker(tickerTimeInterval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			g.Update(tuiLayout)
		}
	}
}

func tuiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	_ = g.DeleteView(mainFrameID)

	v, err := g.SetView(mainFrameID, -1, -1, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	v.Frame = false
	//v.Editable = false

	fmt.Fprintf(v, "Cursor Pos: %d\n", cursorPos)

	// write out the lines
	for i, dir := range rt.Directories {
		if i >= maxY {
			// don't bother render if we got way too many projects
			break
		}

		if i == cursorPos {
			writeSelectedProjectLine(v, dir)
		} else {
			writeProjectLine(v, dir)
		}
	}

	return nil
}

func writeProjectLine(v io.Writer, project Dir) {
	fmt.Fprint(v, aurora.Inverse(" "))
	fmt.Fprintf(v, " %s\n", projectLabel(project))
}

func writeSelectedProjectLine(v io.Writer, project Dir) {
	fmt.Fprint(v, aurora.Inverse(fmt.Sprintf("❯ %s \n", projectLabel(project))))
}

func projectLabel(project Dir) string {
	if selectedListStyle == listStyleLong {
		return project.Path
	}
	return project.Label
}

func quitKeyBinding(g *gocui.Gui) error {
	quitHandler := func(g *gocui.Gui, v *gocui.View) error {
		close(done)
		fmt.Print(".") // print . so the directory doesn't change
		return gocui.ErrQuit
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitHandler); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quitHandler); err != nil {
		return err
	}

	return nil
}

func addArrowKeyBindings(g *gocui.Gui) error {
	_, maxY := g.Size()

	arrowKeyUpHandler := func(g *gocui.Gui, v *gocui.View) error {
		if cursorPos <= 0 {
			cursorPos = 0
		} else {
			cursorPos--
		}
		return nil
	}

	arrowKeyDownHandler := func(g *gocui.Gui, v *gocui.View) error {
		dirCount := len(rt.Directories) - 1
		if cursorPos >= dirCount {
			cursorPos = dirCount
		} else {
			if cursorPos >= maxY {
				cursorPos = maxY - 1
			} else {
				cursorPos++
			}
		}

		return nil
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, arrowKeyUpHandler); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, arrowKeyDownHandler); err != nil {
		return err
	}

	return nil
}

func addEnterKeyBinding(g *gocui.Gui) error {
	enterKeyBinding := func(g *gocui.Gui, v *gocui.View) error {
		close(done)
		fmt.Print(rt.Directories[cursorPos].Path)
		return gocui.ErrQuit
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, enterKeyBinding); err != nil {
		return err
	}

	return nil
}

func addTabKeyBinding(g *gocui.Gui) error {
	tabKeyBinding := func(g *gocui.Gui, v *gocui.View) error {
		if selectedListStyle == listStyleLong {
			selectedListStyle = listStyleShort
			return nil
		}
		selectedListStyle = listStyleLong
		return nil
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, tabKeyBinding); err != nil {
		return err
	}

	return nil
}
