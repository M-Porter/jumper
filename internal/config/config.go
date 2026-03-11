package config

import "C"
import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/m-porter/configure/v3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type Config struct {
	HomeDir           string
	JumperDir         string
	CacheFileFullPath string
	CacheFile         string
	SearchIncludes    []string         // see configFromFile.SearchIncludes
	SearchExcludes    []string         // see configFromFile.SearchExcludes
	SearchPathStops   []*regexp.Regexp // see configFromFile.SearchPathStops
	SearchMaxDepth    int              // see configFromFile.SearchMaxDepth
	LineIndicator     string           // see configFromFile.LineIndicator

	searchPathStops []string
}

func (c *Config) toFileConfig() *configFromFile {
	return &configFromFile{
		CacheFile:       c.CacheFile,
		SearchIncludes:  c.SearchIncludes,
		SearchExcludes:  c.SearchExcludes,
		SearchPathStops: c.searchPathStops,
		SearchMaxDepth:  c.SearchMaxDepth,
		LineIndicator:   c.LineIndicator,
	}
}

// the config structure as written to the file
type configFromFile struct {
	CacheFile string `mapstructure:"cache_file"`
	// Which paths to include in the search. The starting points.
	SearchIncludes []string `mapstructure:"search_includes"`
	// Which paths to ignore from the search if come across within the search excludes.
	SearchExcludes []string `mapstructure:"search_excludes"`
	// how we determine not to go any deeper when walking
	SearchPathStops []string `mapstructure:"search_path_stops"`
	// how far deep we attempt to search beyond the home directory
	SearchMaxDepth int `mapstructure:"search_max_depth"`
	// The line indicator
	LineIndicator string `mapstructure:"line_indicator"`
}

func (i *configFromFile) toConfig() *Config {
	c := &Config{
		HomeDir:         HomeDir(),
		SearchIncludes:  i.SearchIncludes,
		SearchExcludes:  i.SearchExcludes,
		CacheFile:       i.CacheFile,
		SearchMaxDepth:  i.SearchMaxDepth,
		searchPathStops: i.SearchPathStops,
		LineIndicator:   i.LineIndicator,
	}

	c.CacheFileFullPath = filepath.Join(c.HomeDir, JumperDirname, c.CacheFile)
	c.JumperDir = filepath.Join(c.HomeDir, JumperDirname)

	for _, pathStop := range c.searchPathStops {
		pathStopRegexp := regexp.MustCompile(fmt.Sprintf("%s$", regexp.QuoteMeta(pathStop)))
		c.SearchPathStops = append(config.SearchPathStops, pathStopRegexp)
	}

	return c
}

var config *Config = nil

var conf *configure.Configure

func Init() {
	setupConfigure()
}

func setupConfigure() {
	if conf != nil {
		return
	}

	hd := HomeDir()

	configDirFull := filepath.Join(hd, JumperDirname)
	if _, err := os.Stat(configDirFull); os.IsNotExist(err) {
		err := os.MkdirAll(configDirFull, os.ModePerm)
		cobra.CheckErr(err)
	}

	conf = configure.New()
	cobra.CheckErr(conf.SetConfigName(Filename))
	cobra.CheckErr(conf.SetConfigType(Type))
	cobra.CheckErr(conf.SetConfigDir(configDirFull))
	cobra.CheckErr(conf.SetWriteIfNotExists(true))
	cobra.CheckErr(conf.SetDefaults(
		configFromFile{
			CacheFile:       DefaultCacheFile,
			SearchIncludes:  defaultSearchIncludes,
			SearchExcludes:  defaultSearchExcludes,
			SearchPathStops: defaultSearchPathStops,
			SearchMaxDepth:  defaultSearchMaxDepth,
			LineIndicator:   defaultLineIndicator,
		},
	))
}

func Get() *Config {
	setupConfigure()

	if config != nil {
		return config
	}

	internalConfig := &configFromFile{}
	cobra.CheckErr(conf.Get(internalConfig))

	config = internalConfig.toConfig()

	return config
}

func Save(c *Config) {
	setupConfigure()

	if config == nil {
		cobra.CheckErr(errors.New("cannot save config: config is nil"))
	}

	config = c
	cobra.CheckErr(conf.Save(config.toFileConfig()))
}

func Filepath() string {
	return filepath.Join(HomeDir(), JumperDirname, fmt.Sprintf("%s.%s", Filename, Type))
}

func HomeDir() string {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)
	return hd
}
