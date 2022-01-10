package config

var JumperDirname = ".jumper"
var Filename = "config"
var Type = "yml"
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

var defaultSearchMaxDepth = 6
