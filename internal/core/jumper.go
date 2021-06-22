package core

import "time"

type listStyle int

var (
	resultsListMaxH = 0
)

const (
	tickerTimeInterval = time.Millisecond * 10

	listStyleShort listStyle = 0
	listStyleLong  listStyle = 1
)

func Run(args []string) error {
	app = &Application{Directories: []string{}}
	go app.Setup()
	tui := NewTUI(app)
	return tui.Run()
}
