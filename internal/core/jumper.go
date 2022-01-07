package core

func RunAnalyzer(runInDebugMode bool) {
	app := NewApp(runInDebugMode)
	app.Setup()
	app.Analyze()
}
