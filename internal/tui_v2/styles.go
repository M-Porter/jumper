package tui_v2

type listStyle int

const (
	listStyleShort listStyle = iota
	listStyleLong
	listStyleDetailed
)

var (
	listStyles = []listStyle{listStyleShort, listStyleLong, listStyleDetailed}
)
