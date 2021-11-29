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
	cfg            *Cfg
	option         *option
	fileWriters    map[string]zapcore.WriteSyncer		// <path, WriteSyncer>
	consoleWriters map[string]zapcore.WriteSyncer		// <name, WriteSyncer>

	fileCores      map[string]zapcore.Core
	consoleCores   map[string]zapcore.Core
}

var (
  dfLogger *Elogger
	loggers  map[string]*Elogger
	mu       sync.Mutex
)

func newElogger(name string, cfg* Cfg) *Elogger {
	mu.Lock()
	defer mu.Unlock()

	if err := cfg.validateAndCheck(); err != nil {
		return nil
	}

	if _, exist := loggers[name]; exist {
		syslog.Warnf("old logger named '%s' found, will be replace by new one")
	}

	out := new(Elogger)

	out.name           = name
	out.cfg            = cfg
	out.option         = newOption(cfg)   // cfg checked earlier 
	out.fileCores      = map[string]zapcore.Core{}
	out.fileWriters    = map[string]zapcore.WriteSyncer{}
	out.consoleCores   = map[string]zapcore.Core{}
	out.consoleWriters = map[string]zapcore.WriteSyncer{}

	if loggers == nil {
		loggers = map[string]*Elogger{}
	}

	return out
}


// GetLog ...
// param is to set tag and filename, the first one is tag and second one is filename
// if not set, it will using cfg in logger
func (l *Elogger)getLog(options ...Option) Elog {

	var logger *zap.Logger
	var path string

	opt := l.option.clone().applyOptions(options...)
	
	var cores []zapcore.Core
	cores = append(cores, l.getConsoleCore(1, opt))

	// setting output file core if neededs
	{
		if opt.filename != "" {
			
			if !filepath.IsAbs(opt.filename) {
				path = filepath.Join(l.cfg.Dir, l.cfg.Group, opt.filename)
				if !strings.HasSuffix(opt.filename, ".log"){
					path += ".log"
				}
			}
			path = getRepresentPathValue(path, l.name)
			
			cores = append(cores, l.getFileCore(path, opt))
		}
	}

	logger = zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(zapcore.Level(opt.fileStackLevel)))

	if opt.tagSet{
		return logger.Sugar().Named(opt.tag)
	}

	return logger.Sugar()
}

func (l *Elogger)getConsoleCore(fd int, opt *option) zapcore.Core{
	return zapcore.NewCore(getEncoder(opt.consoleColor), getConsoleWriter(fd), zapcore.Level(opt.consoleLevel))
}

func (l *Elogger)getFileCore(path string, opt *option) zapcore.Core{
	return zapcore.NewCore(getEncoder(l.cfg.FileColor), getFileWriter(path, l.cfg), zapcore.Level(opt.fileLevel))
}

func initDfLogger(cfg *Cfg) {

	if cfg == nil {
		if dfLogger == nil {
			dfLogger = NewLogger(cfgDefaultKey, &dfCfg)
		}
		return
	}

	dfLogger = NewLogger(cfgDefaultKey, cfg)
	dfCfg = *cfg
}