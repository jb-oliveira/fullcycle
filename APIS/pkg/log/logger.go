package log

import (
	"os"

	"go.uber.org/zap"
)

func init() {
	logger := zap.Must(zap.NewProduction())
	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	zap.ReplaceGlobals(logger)
}

func Info(msg string) {
	zap.L().Info(msg)
}

func Error(msg string) {
	zap.L().Error(msg)
}

func Warn(msg string) {
	zap.L().Warn(msg)
}

func Debug(msg string) {
	zap.L().Debug(msg)
}
