package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"regexp"
)

const (
	JumperDirname    = ".jumper"
	Filename         = "config"
	Type             = "yml"
	DefaultCacheFile = "cache"
)

var (
	defaultSearchIncludes = []string{
		"development/",
		"dev/",
		"xcode-projects/",
		"repos/",
	}
	defaultSearchExcludes = []string{
		"/node_modules",
		"/bin",
		"/temp",
		"/tmp",
		"/vendor",
		"/venv",
		"/ios/Pods",
	}
	defaultSearchPathStops = []string{
		"/.git",
		"/Gemfile",
		"/package.json",
		"/go.mod",
		"/setup.py",
		"/pyproject.toml",
	}
)

type Config struct {
	HomeDir           string
	JumperDir         string
	CacheFileFullPath string
	CacheFile         string
	SearchIncludes    []string
	SearchExcludes    []string
	SearchPathStops   []*regexp.Regexp
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
}

var C *Config = nil

func Init() {
	hd := HomeDir()

	configDirFull := filepath.Join(hd, JumperDirname)
	if _, err := os.Stat(configDirFull); os.IsNotExist(err) {
		err := os.MkdirAll(configDirFull, os.ModePerm)
		cobra.CheckErr(err)
	}

	viper.SetConfigName(Filename)
	viper.SetConfigType(Type)
	viper.AddConfigPath(configDirFull)

	viper.SetDefault("cache_file", DefaultCacheFile)
	viper.SetDefault("search_includes", defaultSearchIncludes)
	viper.SetDefault("search_excludes", defaultSearchExcludes)
	viper.SetDefault("search_path_stops", defaultSearchPathStops)

	err := viper.SafeWriteConfig()
	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		// ignore, this is ok. just means the config already exists, so
		// we don't need to write a new one
	} else {
		cobra.CheckErr(err)
	}

	err = viper.ReadInConfig()
	cobra.CheckErr(err)

	internalConf := &configFromFile{}
	err = viper.Unmarshal(internalConf)
	cobra.CheckErr(err)

	// write the config after reading and setting defaults in case something
	// had changed or a new config value was added
	go func() {
		err = viper.WriteConfig()
		cobra.CheckErr(err)
	}()

	// copy internalConf to C
	C = &Config{
		HomeDir:        hd,
		SearchIncludes: internalConf.SearchIncludes,
		SearchExcludes: internalConf.SearchExcludes,
		CacheFile:      internalConf.CacheFile,
	}

	C.CacheFileFullPath = filepath.Join(C.HomeDir, JumperDirname, C.CacheFile)
	C.JumperDir = filepath.Join(C.HomeDir, JumperDirname)

	for _, pathStop := range internalConf.SearchPathStops {
		pathStopRegexp := regexp.MustCompile(fmt.Sprintf("%s$", regexp.QuoteMeta(pathStop)))
		C.SearchPathStops = append(C.SearchPathStops, pathStopRegexp)
	}
}

func HomeDir() string {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)
	return hd
}
