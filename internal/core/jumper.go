package core

import "time"

const (
	tickerTimeInterval = time.Millisecond * 10
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
