package elog

import (
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Elogger a instance to handle multi elogs, you can use it to make a new elog or get a exist one
type Elogger struct {
	name           string
	cfg            *LoggerCfg
	option         *option
}

var (
	dfLoggerName = "__default__"
  dfLogger *Elogger
	loggers  = map[string]*Elogger{}
	mu       sync.Mutex
)

func newElogger(name string, cfg* LoggerCfg) *Elogger {
	mu.Lock()
	defer mu.Unlock()

	if err := cfg.validateAndCheck(); err != nil {
		syslog.Fatalf("validateAndCheck cfg failed: %s", err)
		return nil
	}

	if _, exist := loggers[name]; exist {
		syslog.Warnf("old logger named '%s' found, will be replace by new one", name)
	}

	out := genElogger(name, cfg)

	loggers[name] = out

	return out
}

func genElogger(name string, cfg *LoggerCfg) *Elogger {

	out := new(Elogger)

	out.name           = name
	out.cfg            = cfg
	out.option         = newOpt().applyCfg(cfg)

	return out
}

func getLogger(name ...string) *Elogger {
	if len(name) == 0 {
		return dfLogger
	}

	mu.Lock()
	defer mu.Unlock()

	dfSet := false

	for _, _name := range name{
		if logger, exist := loggers[_name]; exist {
			return logger
		}

		if _name == "" {
			dfSet = true
		}
	}

	if dfSet {
		return dfLogger
	}

	return nil
}

// GetLog ...
// param is to set tag and filename, the first one is tag and second one is filename
// if not set, it will using cfg in logger
func (l *Elogger)getLog(opts ...*option) Elog {

	console_needed := false
	file_needed    := false

	opt := l.option.clone().applyOptions(opts...)
	
	var cores []zapcore.Core

	// setting console output core
	cores = append(cores, l.getConsoleCore(1, opt))
	console_needed = true

	// setting output file core if neededs
	if opt.filename != "" {
		path := opt.filename

		if !filepath.IsAbs(opt.filename) {
			path = filepath.Join(l.cfg.Dir, l.cfg.Group, opt.filename)
			if !strings.HasSuffix(opt.filename, ".log"){
				path += ".log"
			}
		}

		file_needed = true
		cores = append(cores, l.getFileCore(getRepresentPathValue(path, l.name), opt))
	}
	
	// get lowest stack level
	stackLevel := LEVEL_NONE
	if console_needed && stackLevel > opt.consoleStackLevel { stackLevel = opt.consoleStackLevel }
	if file_needed    && stackLevel > opt.fileStackLevel    { stackLevel = opt.fileStackLevel }

	logger := zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(zapcore.Level(stackLevel)))
	if len(opt.tags) > 0 {
		logger = logger.Named("[" + strings.Join(opt.tags, ".") + "]")
	}

	return logger.Sugar()
}

func (l *Elogger)getConsoleCore(fd int, opt *option) zapcore.Core{
	return zapcore.NewCore(getEncoder(opt.consoleColor, opt.consoleStackLevel), getConsoleWriter(fd), zapcore.Level(opt.consoleLevel))
}

func (l *Elogger)getFileCore(path string, opt *option) zapcore.Core{
	return zapcore.NewCore(getEncoder(l.cfg.FileColor, opt.fileStackLevel), getFileWriter(path, l.cfg), zapcore.Level(opt.fileLevel))
}

func initDfLogger(cfg *Cfg) {

	if cfg == nil {
		if dfLogger == nil {
			dfLogger = NewLogger(dfLoggerName, &dfCfg)
		}
		return
	}

	dfLogger = NewLogger(dfLoggerName, cfg)
	dfCfg = *cfg
}