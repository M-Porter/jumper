package logger

import (
	"fmt"
	"github.com/m-porter/jumper/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"time"
)

var (
	Debug                      = false
	loggerInstance *zap.Logger = nil
)

func Log(msg string, fields ...zap.Field) {
	if Debug {
		if loggerInstance == nil {
			setup()
		}
		loggerInstance.Debug(msg, fields...)
	}
}

func setup() {
	t := time.Now()
	logFileName := fmt.Sprintf("%0.4d-%0.2d-%0.2d.log", t.Year(), t.Month(), t.Day())
	outputPath := filepath.Join(config.HomeDir(), config.JumperDirname, logFileName)
	c := zap.NewDevelopmentConfig()
	c.OutputPaths = []string{outputPath}
	c.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, _ := c.Build()
	defer logger.Sync()
	loggerInstance = logger
}
