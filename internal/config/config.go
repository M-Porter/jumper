package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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

var config *Config = nil

func HomeDir() string {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)
	return hd
}

func Filepath() string {
	return filepath.Join(HomeDir(), JumperDirname, fmt.Sprintf("%s.%s", Filename, Type))
}

func Get() Config {
	if config == nil {
		initialize()
	}
	return *config
}

func initialize() {
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
	var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
	if errors.As(err, &configFileAlreadyExistsError) {
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

	// copy internalConf to config
	config = &Config{
		HomeDir:        hd,
		SearchIncludes: internalConf.SearchIncludes,
		SearchExcludes: internalConf.SearchExcludes,
		CacheFile:      internalConf.CacheFile,
		SearchMaxDepth: internalConf.SearchMaxDepth,
	}

	config.CacheFileFullPath = filepath.Join(config.HomeDir, JumperDirname, config.CacheFile)
	config.JumperDir = filepath.Join(config.HomeDir, JumperDirname)

	for _, pathStop := range internalConf.SearchPathStops {
		pathStopRegexp := regexp.MustCompile(fmt.Sprintf("%s$", regexp.QuoteMeta(pathStop)))
		config.SearchPathStops = append(config.SearchPathStops, pathStopRegexp)
	}
}

func (c *Config) Save() {
	hd := HomeDir()

	internalConf := c.toInternalConfig()

	configDirFull := filepath.Join(hd, JumperDirname)
	if _, err := os.Stat(configDirFull); os.IsNotExist(err) {
		err := os.MkdirAll(configDirFull, os.ModePerm)
		cobra.CheckErr(err)
	}

	out, err := yaml.Marshal(internalConf)
	cobra.CheckErr(err)

	fo, err := os.Open(Filepath())
	cobra.CheckErr(err)
	defer func() {
		cobra.CheckErr(fo.Close())
	}()

	_, err = fo.Write(out)
	cobra.CheckErr(err)

	config = c
}

func (ic *configFromFile) ToConfig() Config {
	return Config{
		HomeDir:        HomeDir(),
		SearchIncludes: ic.SearchIncludes,
		SearchExcludes: ic.SearchExcludes,
		CacheFile:      ic.CacheFile,
		SearchMaxDepth: ic.SearchMaxDepth,
	}
}

func (c *Config) toInternalConfig() configFromFile {
	return configFromFile{
		CacheFile:      c.CacheFile,
		SearchIncludes: c.SearchIncludes,
		SearchExcludes: c.SearchExcludes,
		SearchMaxDepth: c.SearchMaxDepth,
	}
}
