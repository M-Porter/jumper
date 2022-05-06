package core

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/m-porter/jumper/internal/config"
	"github.com/m-porter/jumper/internal/lib"
	"github.com/m-porter/jumper/internal/logger"
	"github.com/saracen/walker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Application struct {
	Directories []string
	Cache       *Cache
}

func (a *Application) Setup() {
	isStale, err := isCacheStale(config.C.CacheFileFullPath)
	if os.IsNotExist(err) {
		a.Analyze()
	} else {
		cobra.CheckErr(err)
		if isStale {
			a.Analyze()
		}
	}

	c, err := readFromCache(config.C.CacheFileFullPath)
	if c != nil {
		a.Directories = c.Directories
		a.Cache = c
	}
}

func (a *Application) Analyze() {
	excludeRegex := lib.RegexpJoinPartsOr(config.C.SearchExcludes)

	var projectDirs []string
	var wg sync.WaitGroup

	counter := 0

	wg.Add(len(config.C.SearchIncludes))

	for _, search := range config.C.SearchIncludes {
		fullSearch := filepath.Join(config.C.HomeDir, search)
		logger.Log("analyzing path", zap.String("path", fullSearch))

		go func(inclPath string) {
			defer wg.Done()

			// walker panics on directories that don't exist so lets make sure
			// it does first
			if _, err := os.Stat(inclPath); os.IsNotExist(err) {
				logger.Log("skipping directory: IsNotExist", zap.String("path", inclPath))
				return
			}

			var mDirs []string

			walkFn := func(p string, fi os.FileInfo) error {
				counter++

				if excludeRegex.MatchString(p) {
					logger.Log("directory matches excludes", zap.String("path", p))
					return filepath.SkipDir
				}

				for _, re := range config.C.SearchPathStops {
					if re.MatchString(p) {
						cleanPath := filepath.Dir(p)
						projectDirs = append(projectDirs, cleanPath)
						mDirs = append(mDirs, cleanPath)

						logger.Log("appending directory", zap.String("path", cleanPath))

						//SkipDir to tell the walker to not go any further
						return filepath.SkipDir
					}
				}

				if !canSearchDeeper(p) {
					//SkipDir to tell the walker to not go any further
					return filepath.SkipDir
				}

				return nil
			}

			errCallback := walker.WithErrorCallback(func(pathname string, err error) error {
				if os.IsNotExist(err) {
					return nil
				}
				if os.IsPermission(err) {
					return nil
				}
				return err
			})

			err := walker.Walk(inclPath, walkFn, errCallback)
			cobra.CheckErr(err)
		}(fullSearch)
	}

	wg.Wait()

	projectDirs = lib.RemoveDuplicates(projectDirs)

	logger.Log("number of directories walked", zap.Int("count", counter))
	logger.Log("projects found", zap.Int("count", len(projectDirs)))

	err := writeToCache(config.C.CacheFileFullPath, projectDirs)
	if err != nil {
		logger.Log("failed writing to cache")
		cobra.CheckErr(err)
	}
}

func canSearchDeeper(path string) bool {
	return len(strings.Split(filepath.Dir(path), string(filepath.Separator))) <= config.C.SearchMaxDepth
}

func NewApp() *Application {
	return &Application{
		Directories: []string{},
	}
}
