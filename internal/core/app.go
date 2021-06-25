package core

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

type Application struct {
	Directories []string
	Logger      *zap.Logger
}

func (a *Application) Setup() {
	isStale, err := isCacheStale(Config.CacheFileFullPath)
	if os.IsNotExist(err) {
		Analyze(false)
	} else {
		cobra.CheckErr(err)
		if isStale {
			Analyze(false)
		}
	}

	c, err := readFromCache(Config.CacheFileFullPath)
	if c != nil {
		a.Directories = c.Directories
	}
}

func NewApp(debug bool) *Application {
	return &Application{
		Directories: []string{},
		Logger:      NewLogger(debug),
	}
}

func NewLogger(debug bool) *zap.Logger {
	t := time.Now()
	logFileName := fmt.Sprintf("debug_%0.4d-%0.2d-%0.2d.log", t.Year(), t.Month(), t.Day())
	outputPath := filepath.Join(Config.HomeDir, JumperDirname, logFileName)

	level := zapcore.FatalLevel
	if debug {
		level = zapcore.DebugLevel
	}

	c := zap.NewProductionConfig()
	c.OutputPaths = []string{outputPath}
	c.Level = zap.NewAtomicLevelAt(level)
	logger, _ := c.Build()
	defer logger.Sync()
	return logger
}
