package core

import (
	"github.com/spf13/cobra"
	"os"
)

type Application struct {
	Directories []string
}

var app *Application

func (a *Application) Setup() {
	isStale, err := isCacheStale(Config.cacheFileFullPath)
	if os.IsNotExist(err) {
		Analyze()
	} else {
		cobra.CheckErr(err)
		if isStale {
			Analyze()
		}
	}

	c, err := readFromCache(Config.cacheFileFullPath)
	if c != nil {
		a.Directories = c.Directories
	}
}

func Run(args []string) error {
	app = &Application{Directories: []string{}}
	go app.Setup()
	return tui()
}

//func setup() {
//	isStale, err := isCacheStale(Config.cacheFileFullPath)
//	if os.IsNotExist(err) {
//		Analyze()
//	} else {
//		cobra.CheckErr(err)
//		if isStale {
//			Analyze()
//		}
//	}
//
//	c, err := readFromCache(Config.cacheFileFullPath)
//	if c != nil {
//		app.Directories = c.Directories
//	}
//}

//func mapToShape(dirs []string) []Dir {
//	var r []Dir
//
//	for _, d := range dirs {
//		r = append(r, Dir{
//			Path:  d,
//			Label: filepath.Base(d),
//		})
//	}
//
//	return r
//}
