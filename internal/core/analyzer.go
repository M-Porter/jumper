package core

import (
	"github.com/saracen/walker"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	pathStopRegexGit *regexp.Regexp
)

func init() {
	pathStopRegexGit = regexp.MustCompile("/\\.git$")
}

func Analyze() {
	excludeRegex := regexpJoinPartsOr(Config.SearchExcludes)

	var projectDirs []string
	var wg sync.WaitGroup
	wg.Add(len(Config.SearchIncludes))

	for _, search := range Config.SearchIncludes {
		fullSearch := filepath.Join(Config.homedir, search)

		go func(inclPath string) {
			defer wg.Done()

			// walker panics on directories that don't exist so lets make sure
			// it does first
			if _, err := os.Stat(inclPath); os.IsNotExist(err) {
				return
			}

			var mDirs []string

			walkFn := func(path string, fi os.FileInfo) error {
				if excludeRegex.MatchString(path) {
					return filepath.SkipDir
				}

				if pathStopRegexGit.MatchString(path) {
					projectDirs = append(projectDirs, path)
					mDirs = append(mDirs, path)
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

	err := writeToCache(Config.cacheFileFullPath, removeGitParts(projectDirs))
	cobra.CheckErr(err)
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

func removeGitParts(dirs []string) []string {
	var r []string
	for _, dir := range dirs {
		r = append(r, strings.TrimRight(dir, "/.git"))
	}
	return r
}
