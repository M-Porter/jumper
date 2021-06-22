package core

import "time"

type listStyle uint8

const (
	tickerTimeInterval = time.Millisecond * 10

	listStyleShort    listStyle = 0
	listStyleLong     listStyle = 1
	listStyleDetailed listStyle = 2
)

var (
	listStyles = []listStyle{listStyleShort, listStyleLong, listStyleDetailed}
)

func Run(args []string) error {
	app := NewApp()
	go app.Setup()
	tui := NewTUI(app)
	return tui.Run()
}
