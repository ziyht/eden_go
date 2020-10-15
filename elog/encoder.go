package elog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func myTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02_15:04:05.000"))
}

func getEncoder(colored bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime  = myTimeEncoder // zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if colored {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}