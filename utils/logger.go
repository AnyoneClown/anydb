/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package utils

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var logPath string

func InitLogger() {
	homeDir, _ := os.UserHomeDir()

	logPath = filepath.Join(homeDir, ".anydb", "go.log")
	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", logPath}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}
}
