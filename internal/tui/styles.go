package tui

import "github.com/gookit/color"

var (
	ColorBgDefault = color.BgDefault
	ColorBgGray    = color.HEX("#424242", true)
	ColorFgRed     = color.HEX("#E53935")
	ColorFgBlue    = color.HEX("#60A5FA")
)

type ListStyle int

const (
	ListStyleShort ListStyle = iota
	ListStyleLong
	ListStyleDetailed
)

var (
	ListStyles = []ListStyle{ListStyleShort, ListStyleLong, ListStyleDetailed}
)
