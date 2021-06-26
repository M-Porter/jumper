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

func Run(runInDebugMode bool) error {
	app := NewApp(runInDebugMode)
	tui := NewTUI(app)
	return tui.Run()
}

func RunAnalyzer(runInDebugMode bool) {
	app := NewApp(runInDebugMode)
	app.Setup()
	app.Analyze()
}
