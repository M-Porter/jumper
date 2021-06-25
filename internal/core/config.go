package core

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	JumperDirname    = ".jumper"
	ConfigFilename   = "config"
	ConfigType       = "yml"
	DefaultCacheFile = "cache"
)

var (
	defaultIncludes = []string{
		"development/",
		"dev/",
		"xcode-projects/",
		"repos/",
	}
	defaultExcludes = []string{
		"node_modules/",
		"bin/",
		"temp/",
		"tmp/",
		".git/",
		"vendor/",
	}
)

type Configuration struct {
	HomeDir           string
	JumperDir         string
	CacheFileFullPath string
	CacheFile         string   `mapstructure:"cache_file"`
	SearchIncludes    []string `mapstructure:"search_includes"`
	SearchExcludes    []string `mapstructure:"search_excludes"`
}

var Config *Configuration = nil

func Init() {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)

	configDirFull := filepath.Join(hd, JumperDirname)
	if _, err := os.Stat(configDirFull); os.IsNotExist(err) {
		err := os.MkdirAll(configDirFull, os.ModePerm)
		cobra.CheckErr(err)
	}

	viper.AddConfigPath(configDirFull)
	viper.SetConfigName(ConfigFilename)
	viper.SetConfigType(ConfigType)

	viper.SetDefault("cache_file", DefaultCacheFile)
	viper.SetDefault("search_includes", defaultIncludes)
	viper.SetDefault("search_excludes", defaultExcludes)

	err = viper.SafeWriteConfig()
	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		// ignore, this is ok. just means the core already exists so
		// we don't need to write a new one
	} else {
		cobra.CheckErr(err)
	}

	viper.AutomaticEnv()

	Config = &Configuration{}
	err = viper.Unmarshal(Config)
	cobra.CheckErr(err)

	Config.HomeDir = hd
	Config.CacheFileFullPath = filepath.Join(Config.HomeDir, JumperDirname, Config.CacheFile)
	Config.JumperDir = filepath.Join(Config.HomeDir, JumperDirname)
}
