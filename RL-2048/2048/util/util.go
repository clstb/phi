package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateLogger() *zap.Logger {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.DisableStacktrace = true
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	_logger, err := loggerConfig.Build()
	Must(err)
	return _logger
}
