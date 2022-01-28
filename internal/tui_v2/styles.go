package tui_v2

import "github.com/charmbracelet/lipgloss"

const (
	colorGray = lipgloss.Color("#3F3F46")
	colorRed  = lipgloss.Color("#EF4444")
	colorBlue = lipgloss.Color("#0EA5E9")
)

type listStyle int

const (
	listStyleShort listStyle = iota
	listStyleLong
	listStyleDetailed
)

var (
	listStyles = []listStyle{listStyleShort, listStyleLong, listStyleDetailed}
)
