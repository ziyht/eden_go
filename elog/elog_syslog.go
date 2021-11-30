package elog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	syslogCfg = genSyslogCfg()
	syslog Elog
)

func genSyslogCfg() Cfg {
	return Cfg{
		Dir              : "",
		Group            : "",
		FileName         : "",
		MaxSize          : 0,
		MaxBackups       : 0,
		MaxAge           : 0,
		ConsoleLevel     : LEVELS_DEBUG,
		ConsoleColor     : true,
		ConsoleStackLevel: LEVELS_DEBUG,
		FileLevel        : "",
		FileColor        : true,
		FileStackLevel   : LEVELS_DEBUG,
		Compress         : false,
	}
}

func initSyslog() {

	opt := newOption(&syslogCfg)

	coreConsole := zapcore.NewCore(getEncoder(true, opt.consoleStackLevel), zapcore.AddSync(os.Stdout), zapcore.Level(opt.consoleLevel))
	logger := zap.New(zapcore.NewTee(coreConsole), zap.AddStacktrace(zapcore.Level(opt.consoleStackLevel))).Named("[ELOG]")
	syslog = logger.Sugar()
}