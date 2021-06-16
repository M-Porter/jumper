package core

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type Runtime struct {
	Directories []Dir
}

type Dir struct {
	Path  string
	Label string
}

var rt *Runtime

func Run(args []string) error {
	rt = &Runtime{Directories: []Dir{}}
	go setup()
	return tui()
}

func setup() {
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
		rt.Directories = mapToShape(c.Directories)
	}
}

func mapToShape(dirs []string) []Dir {
	var r []Dir

	for _, d := range dirs {
		r = append(r, Dir{
			Path:  d,
			Label: filepath.Base(d),
		})
	}

	return r
}
