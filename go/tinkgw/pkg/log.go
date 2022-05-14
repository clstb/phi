package pkg

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

var _, Sugar = func() (*zap.Logger, *zap.SugaredLogger) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	_logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	_sugar := _logger.Sugar()
	return _logger, _sugar
}()
