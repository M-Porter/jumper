package jumper

import (
	"github.com/spf13/cobra"
	"os"
)

type Runtime struct {
	Directories []string
}

var rt *Runtime

func Run(args []string) error {
	rt = &Runtime{
		Directories: []string{},
	}

	setup()
	return tui()
}

func setup() {
	go func() {
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
			rt.Directories = c.Directories
		}
	}()
}
