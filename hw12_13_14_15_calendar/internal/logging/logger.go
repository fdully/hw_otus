package logging

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey struct{}

var fallbackLogger *zap.SugaredLogger

func InitLog(severity int, logFile string) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "severity"
	config.Level = zap.NewAtomicLevelAt(zapcore.Level(severity))

	if logFile != "" {
		config.OutputPaths = append(config.OutputPaths, logFile)
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}
	fallbackLogger = logger.Named("calendar").Sugar()

	return nil
}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey{}).(*zap.SugaredLogger); ok {
		return logger
	}
	return fallbackLogger
}
