package server

import (
	pb "github.com/clstb/phi/go/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

type LedgerServer struct {
	pb.BeanAccountServiceServer
	Logger *zap.SugaredLogger
}

var _, sugar = func() (*zap.Logger, *zap.SugaredLogger) {
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

func NewServer() *LedgerServer {
	s := &LedgerServer{Logger: sugar}
	return s
}
