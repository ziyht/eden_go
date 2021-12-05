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
}

var (
  dfLogger *Elogger
	loggers  = map[string]*Elogger{}
	mu       sync.Mutex
)

func newElogger(name string, cfg* Cfg) *Elogger {
	mu.Lock()
	defer mu.Unlock()

	if err := cfg.validateAndCheck(); err != nil {
		return nil
	}

	if _, exist := loggers[name]; exist {
		syslog.Warnf("old logger named '%s' found, will be replace by new one", name)
	}

	out := genElogger(name, cfg)

	loggers[name] = out

	return out
}

func genElogger(name string, cfg *Cfg) *Elogger {

	out := new(Elogger)

	out.name           = name
	out.cfg            = cfg
	out.cfg.name       = name
	out.option         = newOption(cfg)   // cfg checked earlier 

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
	return zapcore.NewCore(getEncoder(opt.consoleColor, opt.consoleStackLevel), getConsoleWriter(fd), zapcore.Level(opt.consoleLevel))
}

func (l *Elogger)getFileCore(path string, opt *option) zapcore.Core{
	return zapcore.NewCore(getEncoder(l.cfg.FileColor, opt.fileStackLevel), getFileWriter(path, l.cfg), zapcore.Level(opt.fileLevel))
}

func initDfLogger(cfg *Cfg) {

	if cfg == nil {
		if dfLogger == nil {
			dfLogger = NewLogger(cfgDefaultName, &dfCfg)
		}
		return
	}

	dfLogger = NewLogger(cfgDefaultName, cfg)
	dfCfg = *cfg
}