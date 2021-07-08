package core

import "time"

type listStyle uint

const (
	tickerTimeInterval = time.Millisecond * 10
)

const (
	listStyleShort    listStyle = iota
	listStyleLong     listStyle = iota
	listStyleDetailed listStyle = iota
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
