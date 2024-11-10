package elog

var (
	version = "1.0.2"
)

const sampleCfg =
`
#
# Tag representation for dir, group, filename
#    <HOSTNAME> -> hostname of current machine
#    <APP>      -> binary file name of current application
#    <LOG>      -> the name of current logger, in default cfg, the name is 'default'
#
#  note: 
#    1. the key like 'dir', 'group', ... under elog directly is to set default value,
#       you do not need to set it because all of them have a default value inside
#

elog:
  
  # default settings
  dir        : logs                # default logs
  group      : <HOSTNAME>          # default <HOSTNAME>, if set, real dir will be $Dir/$Group
  filename   : <APP>_<LOG>         # default <LOG>, will not write to file if set empty, real file path will be $Dir/$Group/$File
  console    : stdout              # default stdout, you can set stderr instead
  max_size   : 100                 # default 100, unit MB
  max_backup : 7                   # default 7
  max_age    : 7                   # default 7
  compress   : false               # default false
  f_level    : debug               # default debug,       level for file, valid value is [debug, info, warn, error, fatal, panic]
  f_slevel   : warn                # default warn , stack level for file, valid value is [debug, info, warn, error, fatal, panic]
  f_color    : false               # default auto,        color for file, valid value is [auto, true, false]
  c_level    : debug               # default info ,       level for console, valid value is [debug, info, warn, error, fatal, panic]
  c_slevel   : warn                # default error, stack level for console, valid value is [debug, info, warn, error, fatal, panic]
  c_color    : true                # default auto ,       color for console, valid value is [auto, true, false]

  # mode 1
  log1:
    # filename: <APP>_<LOG>        # if not set, will inherit from default value set in elog.filename
    tag    :  log1
    c_level:  info
    f_level:  debug       

  # mode 2
  log2:
  - tag         : log2                # first no-empty tag will take effect, nexts will be skipped
    name        : console             # not used now
    console     : stdout              # console setting
    level       : info                # log level
    slevel      : error               # stack level
    color       : auto                # color 
  - name        : file
    dir         : logs                # default logs
    group       : <HOSTNAME>          # default <HOSTNAME>, if set, real dir will be $dir/$group
    filename    : <APP>_<LOG>         # default <LOG_NAME>, will not write to file if set empty, real file path will be $dir/$group/$file_name
    max_size    : 100                 # default 100, unit MB
    max_backup  : 7                   # default 7
    max_age     : 7                   # default 7
    compress    : false               # default false
    level       : debug               # default debug, for file, valid value is [debug, info, warn, error, fatal, panic]
    slevel      : warn                # default warn , for file, valid value is [debug, info, warn, error, fatal, panic]
    color       : false               # default false, for file

  # mode 2
  multi_file:
  - tag     : multi_file
    filename: <APP>_<LOG>_debug
    level   : [ debug, debug ]
  - filename: <APP>_<LOG>_info
    level   : [ info, info ]
  - filename: <APP>_<LOG>_warn
    level   : [ warn, warn ]
  - filename: <APP>_<LOG>_err
    level   : [ error, error ]

  only_console:
  - console: stdout
    level  : info
    slevel : error
`

// InitFromConfigFile - this will init loggers from a config file, support multi file types like yaml, yml, json, toml...
//  note that the logger with same name will be replaced with new one
//  you can get the sample cfg content from SampleCfg()
//  the root key is elog
func InitFromFile(path string) {
	cfgs := cfgMan.parsingCfgsFromFile(path)

	for name, cfg := range cfgs {
		newEloggerV2(name, cfg)
	}
}

func SampleCfgStr() string {
	return sampleCfg
}

// GenDfCfg - returned a new dfCfg, it will always be the same
//   note: dfLoggerCfg can be modified by InitFromFile
func NewDfCfg() *Cfg {
	return dfLoggerCfg.Clone()
}

// SetDfCfg - set internal default cfg, this will be used in new logger, have no effect to existing loggers
func SetDfCfg(cfg *Cfg) {
	err := cfg.validateAndCheck()
  if err != nil {
    syslog.Warnf("SetDfCfg() failed: %s")
  }
	dfLoggerCfg = *cfg.Clone()
}

// NewLogger generator a new Elogger
func NewLogger(name string, cfg* Cfg) *Elogger {
	return newElogger(name, cfg)
}

// SysLogger - returns the internal syslogger, you can using this to make Elog to logout before init logger operations
//   - SysLogger do not write to file in default
func SysLogger() *Elogger {
	return syslogger
}

// Logger - get a logger by name
//  - no name passed, it returned default logger
//  - passed > 1, return the first match, and "" represent to default logger(lowest priority)
//  - if not found, return nil
func Logger(name ...string) *Elogger {
	return getLogger(name...)
}

// LoggerFromFile - gen a logger instance from conf file, will only parse the specified fields
//  - keys can be path like "test.system.app1.log"
func LoggerFromFile(file string, keys string) *Elogger {
	cfg, err := parsingLoggerCfgFromFile(file, keys)
	if err != nil {
		syslog.Warnf("returned default logger for LoggerFromFile failed: %s", err.Error())
    return Logger()
	}
	return newEloggerV2("", cfg)
}

func LoggerFromInterface(in interface{})*Elogger{
	cfg, err := parsingLoggerCfgSmart(&dfLoggerCfg, in, nil)
	if err != nil {
		syslog.Warnf("returned default logger for LoggerFromInterface failed: %s", err.Error())
    return Logger()
	}
	return newEloggerV2("", cfg)
}

// Opt - gen a new empty option to set properties in cfg which you want to change
func Opt() *option {
	return newOpt()
}

// Log - get a Elog instance by a logger
//   - Elog is a typedef of *zap.SugaredLogger, so you can use Named() to set tags
//   - every call will create a new instance, recommend cache it first and then use it
func (l *Elogger)Log(opts ...*option) Elog {
	return l.getLog(opts...)
}

// return the name of logger
func (l *Elogger)Name() string {
	return l.name
}

// Log - get a log generated by dfLogger
// the passed in option have high priority, if no options passed, it returned a default Elog
func Log(opts ...*option) Elog {
	return dfLogger.getLog(opts...)
}

func Version() string{
	return version
}