package tui

import (
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/rivo/tview"
)

func colorize(v io.Writer, text string) {
	_, _ = tview.ANSIWriter(v).Write([]byte(text))
}

func ctoc(c color.RGBColor) tcell.Color {
	v := c.Values()
	return tcell.NewRGBColor(int32(v[0]), int32(v[1]), int32(v[2]))
}
