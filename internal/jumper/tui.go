package jumper

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"sync"
	"time"
)

var (
	done  = make(chan struct{})
	tuiWg sync.WaitGroup
)

func tui() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
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

	for {
		select {
		case <-done:
			return
		case <-time.After(time.Millisecond * 100):
			g.Update(tuiLayout)
		}
	}
}

func tuiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	_ = g.DeleteView("hello")
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintf(v, "%d", time.Now().UnixNano())
	}
	return nil
}

func tuiQuit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	fmt.Print(".") // print . so the directory doesn't change
	return gocui.ErrQuit
}
