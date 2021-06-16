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

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, tuiQuit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, tuiQuit); err != nil {
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

	// write out the lines
	for i, dir := range rt.Directories {
		if i >= maxY {
			// don't bother render if we got way too many projects
			break
		}
		writeProjectLine(v, dir)
	}

	return nil
}

func writeProjectLine(v io.Writer, project string) {
	fmt.Fprint(v, aurora.BgIndex(bgColorGray, " "), " ")
	fmt.Fprintf(v, "%s\n", project)
}

func tuiQuit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	fmt.Print(".") // print . so the directory doesn't change
	return gocui.ErrQuit
}
