package elog

import (
	"os"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap/zapcore"
)

var (
	stdoutConsoleWriter = zapcore.AddSync(os.Stdout)
	stderrConsoleWriter = zapcore.AddSync(os.Stderr)

	fileWriters = map[string]zapcore.WriteSyncer{}
	fileWritersMu sync.Mutex
)

func getConsoleWriter(stream string) zapcore.WriteSyncer {

	if stream == "STDOUT" || stream == "stdout"  || stream == "1" {
		return stdoutConsoleWriter
	} else if stream == "STDERR" || stream == "stderr" || stream == "2" {
		return stderrConsoleWriter
	}

	syslog.Warnf("invalid stream(%s), now only support 1: stdout, 2: stderr", stream)
 
	return stdoutConsoleWriter
}

func getFileWriter(path string, cfg *LogCfg) zapcore.WriteSyncer {

	fileWritersMu.Lock()
	defer fileWritersMu.Unlock()

	fileWriter, exist := fileWriters[path]
	if !exist {
		fileWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   path,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackup,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		fileWriters[path] = fileWriter
	}

	return fileWriter
}