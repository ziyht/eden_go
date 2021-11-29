package elog

import (
	"os"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

var (
	stdoutConsoleWriter = zapcore.AddSync(os.Stdout)
	stderrConsoleWriter = zapcore.AddSync(os.Stderr)

	fileWriters = map[string]zapcore.WriteSyncer{}
	fileWritersMu sync.Mutex
)

func getConsoleWriter(id int) zapcore.WriteSyncer {

	if id == 2 { return stderrConsoleWriter }
	if id != 1 {
		syslog.Warnf("invalid id(%d), now only support 1: stdout, 2: stderr", id)
	} 
	return stdoutConsoleWriter
}

func getFileWriter(path string, cfg *Cfg) zapcore.WriteSyncer {

	fileWritersMu.Lock()
	defer fileWritersMu.Unlock()

	fileWriter, exist := fileWriters[path]
	if !exist {
		fileWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		fileWriters[path] = fileWriter
	}

	return fileWriter
}