package log

import (
	"log"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger ...
	Logger *zap.Logger
	// LogOpts ...
	LogOpts []grpc_zap.Option
)

func init() {
	// if dev use dev logger otherwise use prod
	logConfig := zap.NewDevelopmentConfig()
	//zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	//var customFunc grpc_zap.CodeToLevel
	LogOpts = []grpc_zap.Option{
	//grpc_zap.WithLevels(customFunc),
	}
	// make sure that log statements internal to gRPC
	// library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLogger(logger)
	logger.Info("DONE LOG")
}
