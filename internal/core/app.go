package core

import (
	"github.com/spf13/cobra"
	"os"
)

type Application struct {
	Directories []string
}

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

func NewApp() *Application {
	return &Application{Directories: []string{}}
}
