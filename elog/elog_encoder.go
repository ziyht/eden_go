package elog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


func getEncoder(color colorSwitch, stacklevel Level) zapcore.Encoder {
	ecfg := zap.NewProductionEncoderConfig()
  ecfg.EncodeTime       = getTimeEncoder()
	ecfg.EncodeLevel      = getLevelEncoder(color)
	ecfg.ConsoleSeparator = " "

	return NewConsoleEncoder(ecfg, stacklevel)
}
