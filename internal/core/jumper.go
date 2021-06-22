package core

func Run(args []string) error {
	app = &Application{Directories: []string{}}
	go app.Setup()
	tui := NewTUI(app)
	return tui.Run()
}
