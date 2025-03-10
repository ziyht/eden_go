package elog

var (
	syslogger    *Elogger
	
	syslog Elog
	syslogTag  = "ELOG"

	sysloggerCfg = Cfg{
		Tag              : syslogTag,
		Dir              : "logs",
		Group            : "<APP_NAME>",
		FileName         : "",
		MaxSize          : 100,
		MaxBackup        : 7,
		MaxAge           : 7,
		Compress         : false,
		ConsoleLevel     : LEVEL_DEBUG,
		FileLevel        : LEVEL_DEBUG,
		ConsoleStackLevel: LEVEL_WARN,
		FileStackLevel   : LEVEL_WARN,
		ConsoleColor     : ColorAuto,
		FileColor        : ColorAuto,
	}
)

func initSyslog() {
  syslogger = newElogger("syslog", sysloggerCfg.Clone())
	syslog = syslogger.getLog()
}