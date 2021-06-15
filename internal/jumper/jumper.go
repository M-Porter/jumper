package jumper

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
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

	tick := 0

	for {
		fmt.Printf("loop: %d\n", tick)
		time.Sleep(time.Millisecond * 10)

		if tick > 1000 {
			break
		}
		tick++
	}

	return nil
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
