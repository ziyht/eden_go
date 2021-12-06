package elog

var (
	syslogCfg = genSyslogCfg()
	syslogger *Elogger
	syslog Elog
	syslogTag  = "[ELOG]"
)

func genSyslogCfg() Cfg {
	return Cfg{
		Dir              : "logs",
		Group            : "<APP_NAME>",
		FileName         : "",
		MaxSize          : 100,
		MaxBackups       : 7,
		MaxAge           : 7,
		Compress         : false,
		ConsoleLevel     : LEVELS_DEBUG,
		FileLevel        : LEVELS_DEBUG,
		ConsoleStackLevel: LEVELS_DEBUG,
		FileStackLevel   : LEVELS_DEBUG,
		ConsoleColor     : true,
		FileColor        : true,
	}
}

func initSyslog() {
  syslogger = genElogger("syslog", &syslogCfg)
	syslog = syslogger.getLog().Named(syslogTag)
}