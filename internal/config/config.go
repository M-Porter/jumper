package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
}

var C *Config = nil

func Filepath() string {
	return filepath.Join(HomeDir(), JumperDirname, fmt.Sprintf("%s.%s", Filename, Type))
}

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
	viper.SetDefault("search_max_depth", defaultSearchMaxDepth)

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
	// had changed or a new config value was added.
	err = viper.WriteConfig()
	cobra.CheckErr(err)

	// copy internalConf to C
	C = &Config{
		HomeDir:        hd,
		SearchIncludes: internalConf.SearchIncludes,
		SearchExcludes: internalConf.SearchExcludes,
		CacheFile:      internalConf.CacheFile,
		SearchMaxDepth: internalConf.SearchMaxDepth,
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
