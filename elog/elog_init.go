package elog

func init() {
	initSyslog()
	initDfLogger(nil)

	delete(loggers, dfLoggerName)
}