package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func New() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{logger: logger}
}

func NewDevelopment() *Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return &Logger{logger: logger}
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.logger.Sugar().Debugw(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.logger.Sugar().Infow(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.logger.Sugar().Warnw(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.logger.Sugar().Errorw(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.logger.Sugar().Fatalw(msg, fields...)
	os.Exit(1)
}

func (l *Logger) Sync() {
	_ = l.logger.Sync()
}
