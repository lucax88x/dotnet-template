package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateLogger() *zap.Logger {
	encoder := zap.NewDevelopmentEncoderConfig()
	encoder.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	return logger
}
