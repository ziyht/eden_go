package elog

var (
  dfLogger *Elogger

	dfDir            = "logs"
	dfGroup          = "<HOSTNAME>"  
	dfFileName       = "<APP>_<LOG>"
	dfTag            = ""
	dfFileCfgName    = "file"
	dfConsoleCfgName = "console"

	dfLoggerCfg      = Cfg{
		Tag              : dfTag,
		Dir              : dfDir,
		Group            : dfGroup,
		FileName         : dfFileName,
		MaxSize          : 100,
		MaxBackup        : 7,
		MaxAge           : 7,
		ConsoleLevel     : LEVEL_INFO,
		ConsoleColor     : ColorAuto,
		ConsoleStackLevel: LEVEL_ERROR,
		FileLevel        : LEVEL_DEBUG,
		FileColor        : ColorAuto,
		FileStackLevel   : LEVEL_WARN,
		Compress         : false,
	}
)

func init() {
	dfLogger = genEloggerFromCfg("", &dfLoggerCfg)
}