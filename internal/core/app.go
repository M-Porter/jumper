package core

import (
	"fmt"
	"github.com/saracen/walker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	// how we determine not to go any deeper when walking
	pathStops = []*regexp.Regexp{
		regexp.MustCompile("/\\.git$"),
		regexp.MustCompile("/Gemfile$"),
		regexp.MustCompile("/package\\.json$"),
		regexp.MustCompile("/go\\.mod$"),
		regexp.MustCompile("/setup\\.py$"),
		regexp.MustCompile("/pyproject\\.toml$"),
	}

	/*
		todo: implement this
		max depth - how far deep we attempt to go beyond the home directory
	*/
	maxDepth = 6
)

type Application struct {
	Directories []string
	Logger      *zap.Logger
	Cache       *Cache
}

func (a *Application) Setup() {
	isStale, err := isCacheStale(Config.CacheFileFullPath)
	if os.IsNotExist(err) {
		a.Analyze()
	} else {
		cobra.CheckErr(err)
		if isStale {
			a.Analyze()
		}
	}

	c, err := readFromCache(Config.CacheFileFullPath)
	if c != nil {
		a.Directories = c.Directories
		a.Cache = c
	}
}

func (a *Application) Analyze() {
	excludeRegex := regexpJoinPartsOr(Config.SearchExcludes)

	var projectDirs []string
	var wg sync.WaitGroup

	counter := 0

	wg.Add(len(Config.SearchIncludes))

	for _, search := range Config.SearchIncludes {
		fullSearch := filepath.Join(Config.HomeDir, search)
		a.Log("analyzing path", zap.String("path", fullSearch))

		go func(inclPath string) {
			defer wg.Done()

			// walker panics on directories that don't exist so lets make sure
			// it does first
			if _, err := os.Stat(inclPath); os.IsNotExist(err) {
				a.Log("skipping directory: IsNotExist", zap.String("path", inclPath))
				return
			}

			var mDirs []string

			walkFn := func(p string, fi os.FileInfo) error {
				counter++

				if excludeRegex.MatchString(p) {
					a.Log("directory matches excludes", zap.String("path", p))
					return filepath.SkipDir
				}

				for _, re := range pathStops {
					if re.MatchString(p) {
						cleanPath := filepath.Dir(p)
						projectDirs = append(projectDirs, cleanPath)
						mDirs = append(mDirs, cleanPath)

						a.Log("appending directory", zap.String("path", cleanPath))

						//SkipDir to tell the walker to not go any further
						return filepath.SkipDir
					}
				}

				if len(strings.Split(filepath.Dir(p), string(filepath.Separator))) > maxDepth {
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

	projectDirs = removeDuplicates(projectDirs)

	a.Log("number of directories walked", zap.Int("count", counter))
	a.Log("projects found", zap.Int("count", len(projectDirs)))

	err := writeToCache(Config.CacheFileFullPath, projectDirs)
	if err != nil {
		a.Log("failed writing to cache")
		cobra.CheckErr(err)
	}
}

func (a *Application) Log(msg string, fields ...zap.Field) {
	if a.Logger != nil {
		a.Logger.Debug(msg, fields...)
	}
}

func NewApp(debug bool) *Application {
	return &Application{
		Directories: []string{},
		Logger:      NewLogger(debug),
	}
}

func NewLogger(debug bool) *zap.Logger {
	if !debug {
		return nil
	}

	t := time.Now()
	logFileName := fmt.Sprintf("%0.4d-%0.2d-%0.2d.log", t.Year(), t.Month(), t.Day())
	outputPath := filepath.Join(Config.HomeDir, JumperDirname, logFileName)
	c := zap.NewDevelopmentConfig()
	c.OutputPaths = []string{outputPath}
	c.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, _ := c.Build()
	defer logger.Sync()
	return logger
}

func quoteParts(parts []string) []string {
	var escaped []string
	for _, part := range parts {
		escaped = append(escaped, regexp.QuoteMeta(part))
	}
	return escaped
}

func regexpJoinPartsOr(parts []string) *regexp.Regexp {
	return regexp.MustCompile(strings.Join(quoteParts(parts), "|"))
}

func removeDuplicates(dirs []string) []string {
	set := make(map[string]struct{})
	var r []string
	for _, dir := range dirs {
		if _, ok := set[dir]; !ok {
			r = append(r, dir)
			set[dir] = struct{}{}
		}
	}
	return r
}
