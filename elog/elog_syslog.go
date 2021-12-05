package elog

var (
	syslogCfg = genSyslogCfg()
	syslogger *Elogger
	syslog Elog
	syslogTag  = "[ELOG]"
)

func genSyslogCfg() Cfg {
	return Cfg{
		Dir              : "",
		Group            : "",
		FileName         : "",
		MaxSize          : 0,
		MaxBackups       : 0,
		MaxAge           : 0,
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