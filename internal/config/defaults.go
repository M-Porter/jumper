package config

var JumperDirname = ".config/jumper"
var Filename = "config"
var Type = "yaml"
var DefaultCacheFile = "cache"

var defaultSearchIncludes = []string{
	"development/",
	"dev/",
	"xcode-projects/",
	"repos/",
}

var defaultSearchExcludes = []string{
	"/node_modules",
	"/bin",
	"/temp",
	"/tmp",
	"/vendor",
	"/venv",
	"/ios/Pods",
}

var defaultSearchPathStops = []string{
	"/.git",
	"/Gemfile",
	"/package.json",
	"/go.mod",
	"/setup.py",
	"/pyproject.toml",
}

var defaultSearchMaxDepth = 1
