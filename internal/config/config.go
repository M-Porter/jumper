package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	JumperDirname    = ".jumper"
	Filename         = "config"
	Type             = "yml"
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
		"/node_modules",
		"/bin",
		"/temp",
		"/tmp",
		"/vendor",
		"/venv",
		"/ios/Pods",
	}
)

type Config struct {
	HomeDir           string
	JumperDir         string
	CacheFileFullPath string
	CacheFile         string   `mapstructure:"cache_file"`
	SearchIncludes    []string `mapstructure:"search_includes"`
	SearchExcludes    []string `mapstructure:"search_excludes"`
}

var C *Config = nil

func Init() {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)

	configDirFull := filepath.Join(hd, JumperDirname)
	if _, err := os.Stat(configDirFull); os.IsNotExist(err) {
		err := os.MkdirAll(configDirFull, os.ModePerm)
		cobra.CheckErr(err)
	}

	viper.SetConfigName(Filename)
	viper.SetConfigType(Type)
	viper.AddConfigPath(configDirFull)

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

	err = viper.ReadInConfig()
	cobra.CheckErr(err)

	C = &Config{}
	err = viper.Unmarshal(C)
	cobra.CheckErr(err)

	C.HomeDir = hd
	C.CacheFileFullPath = filepath.Join(C.HomeDir, JumperDirname, C.CacheFile)
	C.JumperDir = filepath.Join(C.HomeDir, JumperDirname)
}
