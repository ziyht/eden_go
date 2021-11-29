package elog

// NewLogger generator a new Elogger
func NewLogger(name string, cfg* Cfg) *Elogger {
	return newElogger(name, cfg)
}

// Logger - get the default logger
//  - no name passed, it returned default logger
//  - passed > 1, return the first match, and "" represent to default logger
// you can init the property of DfLogger with name '__default' from yaml file
//
func Logger(name ...string) *Elogger {

	if len(name) == 0 {
		return dfLogger
	}

	mu.Lock()
	defer mu.Unlock()

	for _, _name := range name{
		if logger, exist := loggers[_name]; exist {
			return logger
		}

		if _name == "" {
			return dfLogger
		}
	}

	syslog.DPanicf("can not found any logger in names: %s", name)

	return dfLogger
}

func (l *Elogger)Log(options ...Option) Elog {
	return l.getLog(options...)
}
