package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger struct {
	l *zap.Logger
}

var aLogger AppLogger

func NewLogger() error {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	z, err := config.Build()
	if err != nil {
		return err
	}

	aLogger = AppLogger{
		l: z,
	}
	return nil
}

func (al *AppLogger) Infof(format string, a ...interface{}) {
	al.l.Info(fmt.Sprintf(format, a...))
}

func (al *AppLogger) Errorf(format string, a ...interface{}) {
	al.l.Error(fmt.Sprintf(format, a...))
}

func (al *AppLogger) Fatalf(format string, a ...interface{}) {
	al.l.Fatal(fmt.Sprintf(format, a...))
}

func Info(message string, args ...zapcore.Field) {
	aLogger.l.Info(message, args...)
}

func Warning(message string, args ...zapcore.Field) {
	aLogger.l.Warn(message, args...)
}

func Error(message string, err error, args ...zapcore.Field) {
	if len(args) == 0 {
		aLogger.l.Error(message, zap.Error(err))
		return
	}
	aLogger.l.Error(message, prepareParams(err, args)...)
}

func Fatal(message string, err error, args ...zapcore.Field) {
	if len(args) == 0 {
		aLogger.l.Fatal(message, zap.Error(err))
		return
	}
	aLogger.l.Fatal(message, prepareParams(err, args)...)
}

func prepareParams(err error, args []zapcore.Field) []zapcore.Field {
	params := make([]zapcore.Field, 0, len(args)+1)
	params = append(params, zap.Error(err))
	params = append(params, args...)
	return params
}
