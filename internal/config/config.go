package config

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

var config *Config

var conf *configure.Configure

type Config struct {
	// persisted to YAML
	CacheFile       string   `mapstructure:"cache_file"`
	SearchIncludes  []string `mapstructure:"search_includes"`
	SearchExcludes  []string `mapstructure:"search_excludes"`
	SearchPathStops []string `mapstructure:"search_path_stops"`
	SearchMaxDepth  int      `mapstructure:"search_max_depth"`
	NoNerdFont      bool     `mapstructure:"no_nerd_font"`

	// computed at load time, not persisted
	HomeDir              string
	JumperDir            string
	CacheFileFullPath    string
	SearchPathStopRegexp []*regexp.Regexp
}

func (c *Config) hydrate() {
	hd := HomeDir()
	c.HomeDir = hd
	c.CacheFileFullPath = filepath.Join(hd, JumperDirname, c.CacheFile)
	c.JumperDir = filepath.Join(hd, JumperDirname)

	for _, pathStop := range c.SearchPathStops {
		pathStopRegexp := regexp.MustCompile(fmt.Sprintf("%s$", regexp.QuoteMeta(pathStop)))
		c.SearchPathStopRegexp = append(c.SearchPathStopRegexp, pathStopRegexp)
	}
}

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
		Config{
			CacheFile:       DefaultCacheFile,
			SearchIncludes:  defaultSearchIncludes,
			SearchExcludes:  defaultSearchExcludes,
			SearchPathStops: defaultSearchPathStops,
			SearchMaxDepth:  defaultSearchMaxDepth,
		},
	))
}

func Get() *Config {
	setupConfigure()

	if config != nil {
		return config
	}

	config = &Config{}
	cobra.CheckErr(conf.Get(config))
	config.hydrate()

	return config
}

func Save(c *Config) {
	setupConfigure()

	if config == nil {
		cobra.CheckErr(errors.New("cannot save config: config is nil"))
	}

	config = c
	cobra.CheckErr(conf.Save(config))
}

func Filepath() string {
	return filepath.Join(HomeDir(), JumperDirname, fmt.Sprintf("%s.%s", Filename, Type))
}

func HomeDir() string {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)
	return hd
}
