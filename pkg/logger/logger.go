package logger

import (
	"github.com/google/wire"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ProviderSet is a provider set for wire
var ProviderSet = wire.NewSet(NewLogger)

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewLogger(cfg *config.Config) (Logger, error) {
	config := zap.NewProductionConfig()

	// Set log level based on config
	switch cfg.LogLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		sugaredLogger: logger.Sugar(),
	}, nil
}

// Logger interface implementations
func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.sugaredLogger.Infof(template, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(template, args...)
}
