package elog

import (
	"strings"
)

type sCfgMan struct {
	df         *Cfg
	loggerCfgs map[string]*LoggerCfg
}

var cfgMan *sCfgMan = &sCfgMan{ df: &dfLoggerCfg, loggerCfgs: make(map[string]*LoggerCfg) }

func (cm *sCfgMan)register(name string, cfg *LoggerCfg){
	cm.loggerCfgs[name] = cfg
}

func (cm *sCfgMan)parsingCfgsFromFile(file string) (cfgs map[string]*LoggerCfg) {

	cfgs = parsingCfgsFromFile(file)

  for name, cfg := range cfgs {
		cm.register(name, cfg)
	}

	return cfgs
}

func (cm *sCfgMan)findLogCfgs(path string, cache map[string]*LoggerCfg) []*LogCfg {

	keys := strings.SplitN(path, ".", 2)

	if len(keys) == 0 {
		return nil
	}

	var loggerCfg *LoggerCfg
	if cache != nil {
		loggerCfg = cache[keys[0]]
	}
	if loggerCfg == nil {
		loggerCfg = cm.loggerCfgs[keys[0]];
	}
	if loggerCfg == nil {
		return nil
	}

	if len(keys) == 1 {
		loggerCfg = loggerCfg.Clone()
		return loggerCfg.cfgs
	}

	logCfg := loggerCfg.FindLogCfg(keys[1])
	if logCfg != nil {
		return []*LogCfg{logCfg}
	}

	return nil

}

