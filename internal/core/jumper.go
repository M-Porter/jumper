package core

import "time"

type listStyle int

var (
	searchVal          = ""
	tickerTimeInterval = time.Millisecond * 10
	resultsListMaxH    = 0
)

const (
	listStyleShort listStyle = 0
	listStyleLong  listStyle = 1
)

func Run(args []string) error {
	app = &Application{Directories: []string{}}
	go app.Setup()
	tui := NewTUI(app)
	return tui.Run()
}
