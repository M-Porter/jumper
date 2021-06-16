package core

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/logrusorgru/aurora/v3"
	"io"
	"sync"
	"time"
)

var (
	done               = make(chan struct{})
	tuiWg              sync.WaitGroup
	tickerTimeInterval = time.Millisecond * 100

	mainFrameID = "mainJumperFrame"

	bgColorGray uint8 = 238

	cursorPos = 0
)

func tui() error {
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
	if err := quitKeyBinds(g); err != nil {
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

func writeProjectLine(v io.Writer, project string) {
	fmt.Fprint(v, aurora.BgIndex(bgColorGray, " "))
	fmt.Fprintf(v, " %s\n", project)
}

func writeSelectedProjectLine(v io.Writer, project string) {
	line := fmt.Sprintf("  %s\n", project)
	fmt.Fprint(v, aurora.BgIndex(bgColorGray, line))
}

func quitKeyBinds(g *gocui.Gui) error {
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
