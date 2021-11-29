package elog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const SampleCfg =
`
#
# Tag representation for dir, group, filename
#    <HOSTNAME> -> hostname of current machine
#    <APP_NAME> -> binary file name of current application
#    <LOG_NAME> -> the name of current logger, in __default, it will set to ""
# Note:
#   
#

elog:
  __default:                          # default setting for all logs 
    dir          : 'var/log'               # default var/log
    group        : '<HOSTNAME>'            # default <HOSTNAME>, if set, real dir will be $Dir/$Group
    filename     : '<LOG_NAME>'            # default <LOG_NAME>, will not write to file if set empty, real file path will be $Dir/$Group/$File.log, if not set, logs will not be written
    max_size     : 100                     # default 100, unit MB
    max_backups  : 7                       # default 7
    max_age      : 7                       # default 7
    compress     : false                   # default true
		log_level    : ['info' , 'debug']      # [0] for console, [1] for file, you can set [debug, info, warn, error, fatal, panic]
		log_stack    : ['fatal', 'error']      # [0] for console, [1] for file, you can set [debug, info, warn, error, fatal, panic]
		log_color    : [ true  , true   ]      # [0] for console, [1] for file

  log1:
    group   : ""
    filename: log1                         # it will set from __default 

  log2:
    filename: log2
`

const (
	cfgRootKey    = "elog"          // root key in the config file for elog
	cfgDefaultKey = "__default"     // key of default setting for each log

	defaultTag = ""
	syslogTag  = "elog"
)

// dfCfg the default config for elog, note: it can be reset by config file by key set in cfgDefaultKey
var (
	dfCfg       = *genDfCfg()
	appName     = AppName()
	hostname, _ = os.Hostname()
)

// Cfg ...
type Cfg struct{
	Dir          	    string `yaml:"dir"`
	Group             string `yaml:"group"`
	FileName          string `yaml:"filename"`
	MaxSize     	    int    `yaml:"max_size"`
	MaxBackups  	    int    `yaml:"max_backups"`
	MaxAge      	    int    `yaml:"max_age"`
	ConsoleLevel      string `yaml:"console_level"`
	ConsoleColor      bool   `yaml:"console_color"`
	ConsoleStackLevel string `yaml:"console_stack_level"`
	FileLevel         string `yaml:"file_level"`
	FileColor         bool   `yaml:"file_color"`
	FileStackLevel    string `yaml:"file_stack_level"`
	Compress    	    bool   `yaml:"compress"`

	name    string
	logDir  string
	path    string
}

type rootCfg struct {
	Cfgs map[string]*Cfg  `yaml:"elog"`
}

func genDfCfg() *Cfg {
	return &Cfg{
		Dir              : "logs",
		Group            : "<HOSTNAME>",
		FileName         : "<APP_NAME>",
		MaxSize          : 100,
		MaxBackups       : 7,
		MaxAge           : 7,
		ConsoleLevel     : LEVELS_INFO,
		ConsoleColor     : true,
		ConsoleStackLevel: LEVELS_ERROR,
		FileLevel        : LEVELS_DEBUG,
		FileColor        : true,
		FileStackLevel   : LEVELS_WARN,
		Compress         : false,
		name             : cfgDefaultKey,
	}
}

func (cfg *Cfg)check() (err error) {
	if err = cfg.checkFileRotate(); err != nil {return}
	return cfg.checkLevelStr()
}

func (cfg *Cfg)checkLevelStr() (err error) {
	if _, err = getLevelByStr(cfg.ConsoleLevel);      err != nil {return}
	if _, err = getLevelByStr(cfg.FileLevel);         err != nil {return}
	if _, err = getLevelByStr(cfg.ConsoleStackLevel); err != nil {return}
	if _, err = getLevelByStr(cfg.FileStackLevel);    err != nil {return}

	return
}

func (cfg *Cfg)checkFileRotate() (err error) {
	if cfg.MaxSize < 0 {
		return fmt.Errorf("invalid max_size(%d), should >= 0", cfg.MaxSize)
	}
	if cfg.MaxBackups < 0 {
		return fmt.Errorf("invalid max_backups(%d), should >= 0", cfg.MaxBackups)
	}
	if cfg.MaxAge < 0 {
		return fmt.Errorf("invalid max_age(%d), should >= 0", cfg.MaxAge)
	}

	return
}

func (cfg *Cfg)validate(){
	if cfg.MaxSize    < 0 { cfg.MaxSize    = 0 }
	if cfg.MaxBackups < 0 { cfg.MaxBackups = 0 }
	if cfg.MaxAge     < 0 { cfg.MaxAge     = 0 }

	cfg.logDir = getRepresentPathValue(cfg.Dir + "/" + cfg.Group, cfg.name)
	cfg.path   = getRepresentPathValue(cfg.logDir + "/" + cfg.FileName, cfg.name)
}

func (cfg *Cfg)checkAndValidate()(err error) {
	if err = cfg.check(); err != nil { 
		return
	}

  cfg.validate()
	return 
}

func (cfg *Cfg)validateAndCheck()(err error) {

	cfg.validate()

	if err = cfg.check(); err != nil { 
		return
	}
	return 
}

func readCfgFromYaml(file string) *rootCfg {

	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()
	if err != nil {
		syslog.Fatalf("readCfgFromYaml failed from %s:\n %s", file, err.Error())
	}

	var cfgs rootCfg

	file_, err := filepath.Abs(file)
	if err != nil {
		syslog.Fatalf("readCfgFromYaml failed from %s:\n %s", file, err.Error())
	}

	yamlFile, err := ioutil.ReadFile(file_)
	if err != nil {
		syslog.Fatalf("readCfgFromYaml failed from %s:\n %s", file_, err.Error())
	}

	err = yaml.UnmarshalStrict(yamlFile, &cfgs)
	if err != nil {
		syslog.Fatalf("readCfgFromYaml failed from %s:\n Unmarshal failed:\n %s", file_, err.Error())
	}


	// parsing __default cfg
	{
		// generate default config if not set
		if _, ok := cfgs.Cfgs[cfgDefaultKey]; !ok {
			cfgs.Cfgs[cfgDefaultKey] = genDfCfg()
		}
		__dfCfg := cfgs.Cfgs[cfgDefaultKey]

		err := validateCfg(__dfCfg, &dfCfg, v, cfgRootKey, cfgDefaultKey)
		if err != nil {
			syslog.Fatalf("readCfgFromYaml failed from %s:\n validate '%s' cfg failed:\n %s", file_, cfgDefaultKey, err.Error())
		}
	}

	curDfCfg := cfgs.Cfgs[cfgDefaultKey]

	for name, curcfg := range cfgs.Cfgs {

		if name == cfgDefaultKey {
			continue
		}

		if curcfg == nil {
			curcfg = &Cfg{}
			cfgs.Cfgs[name] = curcfg
		}

		err := validateCfg(curcfg, curDfCfg, v, cfgRootKey, name)
		if err != nil {
			syslog.Fatalf("readCfgFromYaml failed from %s:\n validate '%s' cfg failed:\n %s", file_, name, err.Error())
		}
	}

	return &cfgs
}

func validateCfg(dest *Cfg, dfCfg *Cfg, v *viper.Viper, rootKey string, curKey string) error {

	dirCfgPattern 		     	 := rootKey + "." + curKey + ".dir"
	groupCfgPattern 	     	 := rootKey + "." + curKey + ".group"
	filenameCfgPattern 	     := rootKey + "." + curKey + ".filename"
	maxsizeCfgPattern        := rootKey + "." + curKey + ".max_size"
	maxBackupsCfgPattern     := rootKey + "." + curKey + ".max_backups"
	maxAgeCfgPattern         := rootKey + "." + curKey + ".max_age"
	consoleLevelCfgPattern   := rootKey + "." + curKey + ".console_level"
	consoleColorCfgPattern 	 := rootKey + "." + curKey + ".console_color"
	consoleStackLevelPattern := rootKey + "." + curKey + ".console_stack_level"
	fileLevelCfgPattern    	 := rootKey + "." + curKey + ".file_level"
	fileColorCfgPattern      := rootKey + "." + curKey + ".file_color"
	fileStackLevelCfgPattern := rootKey + "." + curKey + ".file_stack_level"
	compressCfgPattern       := rootKey + "." + curKey + ".compress"

	if !v.IsSet(dirCfgPattern)  					{ dest.Dir   						 = dfCfg.Dir }
	if !v.IsSet(groupCfgPattern)					{ dest.Group 						 = dfCfg.Group }
	if !v.IsSet(filenameCfgPattern)				{ dest.FileName 				 = dfCfg.FileName }
	if !v.IsSet(maxsizeCfgPattern)				{ dest.MaxSize 					 = dfCfg.MaxSize  }
	if !v.IsSet(maxBackupsCfgPattern)			{ dest.MaxBackups 			 = dfCfg.MaxBackups }
	if !v.IsSet(maxAgeCfgPattern)					{ dest.MaxAge 					 = dfCfg.MaxAge }
	if !v.IsSet(consoleLevelCfgPattern)		{ dest.ConsoleLevel 		 = dfCfg.ConsoleLevel }
	if !v.IsSet(consoleColorCfgPattern)		{ dest.ConsoleColor 		 = dfCfg.ConsoleColor }
	if !v.IsSet(consoleStackLevelPattern)	{ dest.ConsoleStackLevel = dfCfg.ConsoleStackLevel }
	if !v.IsSet(fileLevelCfgPattern)			{ dest.FileLevel 			   = dfCfg.FileLevel }
	if !v.IsSet(fileColorCfgPattern)			{ dest.FileColor 			   = dfCfg.FileColor  }
	if !v.IsSet(fileStackLevelCfgPattern)	{ dest.FileStackLevel    = dfCfg.FileStackLevel }
	if !v.IsSet(compressCfgPattern)				{ dest.Compress 			   = dfCfg.Compress }

	return dest.checkAndValidate()
}


func getRepresentPathValue(path string, name string) string {

	if path == "" {
		return path
	}

	path = strings.ReplaceAll(path, "<HOSTNAME>", hostname)
	path = strings.ReplaceAll(path, "<APP_NAME>", appName)

	if name == cfgDefaultKey {
		path = strings.ReplaceAll(path, "<LOG_NAME>", defaultTag)
	} else {
		path = strings.ReplaceAll(path, "<LOG_NAME>", name)
	}

	return path
}

func AppName() string {
	appName, _ := os.Executable()
	appName = filepath.Base(appName)
	appName = strings.TrimSuffix(appName, ".exe")

	return appName
}