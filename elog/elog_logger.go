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
	cfg1           *Cfg
	cfg2           *LoggerCfg
	option         *option
}

var (
	dfLoggerName = "__default__"
  dfLogger *Elogger
	loggers  = map[string]*Elogger{}
	mu       sync.Mutex
)

func newElogger(name string, cfg* Cfg) *Elogger {
	mu.Lock()
	defer mu.Unlock()

	if err := cfg.validateAndCheck(); err != nil {
		syslog.Fatalf("validateAndCheck cfg failed: %s", err)
		return nil
	}

	if _, exist := loggers[name]; exist {
		syslog.Warnf("old logger named '%s' found, will be replace by new one", name)
	}

	out := genEloggerFromCfg(name, cfg)

	loggers[name] = out

	return out
}

func newEloggerV2(name string, cfg* LoggerCfg) *Elogger {
	mu.Lock()
	defer mu.Unlock()

	if err := cfg.validateAndCheck(); err != nil {
		syslog.Fatalf("validateAndCheck cfg failed: %s", err)
		return nil
	}

	if _, exist := loggers[name]; exist {
		syslog.Warnf("old logger named '%s' found, will be replace by new one", name)
	}

	out := genEloggerFromLoggerCfg(name, cfg)

	if name != "" {
		loggers[name] = out
	}

	return out
}

func genEloggerFromCfg(name string, cfg *Cfg) *Elogger {

	out := new(Elogger)

	out.name           = name
	out.cfg1           = cfg
	out.cfg2           = cfg.genLoggerCfg()
	out.option         = newOpt().applyCfg(cfg)

	return out
}

func genEloggerFromLoggerCfg(name string, cfg *LoggerCfg) *Elogger {

	out := new(Elogger)

	out.name           = name
	out.cfg2           = cfg

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

	var cores []zapcore.Core
	var sles []zapcore.LevelEnabler
	var tag string

	opt := newOpt().applyOptions(opts...)
	
	for _, cfg := range l.cfg2.logs {

		cfg = opt.applyToLogCfg(cfg)

		if tag == ""{
			tag = cfg.Tag
		}

		if cfg.file {
			if cfg.FileName == ""{
				continue
			}
			path := cfg.FileName

			if !filepath.IsAbs(cfg.FileName){
				path = filepath.Join(cfg.Dir, cfg.Group, cfg.FileName)
				if !strings.HasSuffix(cfg.FileName, ".log"){
					path += ".log"
				}
			}

			sles = append(sles, cfg.StackLevel)
			cores = append(cores, l.getFileCore(getRepresentPathValue(path, l.name), cfg))
		} else {
			sles = append(sles, cfg.StackLevel)
			cores = append(cores, l.getConsoleCore(cfg))
		}		
	}

	stackLevel := getLevelEnableFromLevelEnables(sles)
	if len(sles) == 1 { stackLevel = sles[0] }
	logger := zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(stackLevel))
	if tag != "" {
		logger = logger.Named("[" + tag + "]")
	}

	return logger.Sugar()
}

func (l *Elogger)getConsoleCore(cfg *LogCfg) zapcore.Core{
	return zapcore.NewCore(getEncoder(cfg.Color, cfg.StackLevel), getConsoleWriter(cfg.Console), cfg.Level)
}

func (l *Elogger)getFileCore(path string, cfg *LogCfg) zapcore.Core{
	return zapcore.NewCore(getEncoder(cfg.Color, cfg.StackLevel), getFileWriter(path, cfg), cfg.Level)
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