package jumper

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	ConfigDir      = ".jumper"
	ConfigFilename = "config"
	ConfigType     = "yml"

	defaultCacheFile = "cache"
	defaultIncludes  = []string{
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

type ConfigShape struct {
	homedir           string
	cacheFileFullPath string
	CacheFile         string   `mapstructure:"cache_file"`
	SearchIncludes    []string `mapstructure:"search_includes"`
	SearchExcludes    []string `mapstructure:"search_excludes"`
}

var Config *ConfigShape = nil

func Init() {
	hd, err := homedir.Dir()
	cobra.CheckErr(err)

	configDirFull := filepath.Join(hd, ConfigDir)
	if _, err := os.Stat(configDirFull); os.IsNotExist(err) {
		err := os.MkdirAll(configDirFull, os.ModePerm)
		cobra.CheckErr(err)
	}

	viper.AddConfigPath(configDirFull)
	viper.SetConfigName(ConfigFilename)
	viper.SetConfigType(ConfigType)

	viper.SetDefault("cache_file", defaultCacheFile)
	viper.SetDefault("search_includes", defaultIncludes)
	viper.SetDefault("search_excludes", defaultExcludes)

	err = viper.SafeWriteConfig()
	if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
		// ignore, this is ok. just means the jumper already exists so
		// we don't need to write a new one
	} else {
		cobra.CheckErr(err)
	}

	viper.AutomaticEnv()

	Config = &ConfigShape{}
	err = viper.Unmarshal(Config)
	cobra.CheckErr(err)

	Config.homedir = hd
	Config.cacheFileFullPath = filepath.Join(Config.homedir, ConfigDir, Config.CacheFile)
}
