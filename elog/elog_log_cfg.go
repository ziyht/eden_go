package elog

import (
	"fmt"
	"os"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const sampleCfg =
`
# Tag representation for dir, group, filename
#    <HOSTNAME> -> hostname of current machine
#    <APP_NAME> -> binary file name of current application
#    <LOG_NAME> -> the name of current logger, in default cfg, it will set to elog
#
#  note: 
#    1. the key like 'dir', 'group', ... under elog directly is to set default value,
#       you do not need to set it because all of them have a default value inside
#    2. for 'level', 'stack_level', 'color' you set only one val for console and file
#       settings, or you can set a arr with two val to set config for console and file
#       respectively 
#

elog:
  dir           : logs                    # default logs
  group         : <HOSTNAME>              # default <HOSTNAME>, if set, real dir will be $Dir/$Group
  filename      : <LOG_NAME>              # default <LOG_NAME>, will not write to file if set empty, real file path will be $Dir/$Group/$File
  max_size      : 100                     # default 100, unit MB
  max_backups   : 7                       # default 7
  max_age       : 7                       # default 7
  compress      : false                   # default false
  level         : [info,  debug]          # default info , debug, [0] for console, [1] for file, valid value is [debug, info, warn, error, fatal, panic]
  stack_level   : [fatal, warn ]          # default fatal, warn , [0] for console, [1] for file, valid value is [debug, info, warn, error, fatal, panic]
  color         : [true,  false]          # default true , true , [0] for console, [1] for file

  log1:
    group   : ""
    filename: log1                         # it will set from __default if not set and will no write to file if set empty

  log2:
    filename: log2
`

const (
	cfgRootKey     = "elog"          // root key in the config file for elog
	cfgDefaultName = "__default"     // key of default setting for each log
	
	defaultFilename = "elog"

	defaultTag = ""
)

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

func (c *Cfg)Clone(newName string) (*Cfg){
	out := *c
	out.name = newName
	return &out 
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
		name             : cfgDefaultName,
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

func (cfg *Cfg)validate() (err error){

	if cfg.MaxSize    < 0 { cfg.MaxSize    = 0 }
	if cfg.MaxBackups < 0 { cfg.MaxBackups = 0 }
	if cfg.MaxAge     < 0 { cfg.MaxAge     = 0 }

	logDir := path.Join(cfg.Dir   , cfg.Group) 
	path   := path.Join(cfg.logDir, cfg.FileName)

	if cfg.logDir, err = filepath.Abs(logDir); err != nil {
		return fmt.Errorf("do Abs() for logDir '%s' failed: %s", logDir, err)
	}
	if cfg.path, err = filepath.Abs(path); err != nil {
		return fmt.Errorf("do Abs() for path '%s' failed: %s", logDir, err)
	}

	return
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

func readCfgFromFile(file string) (cfgs *rootCfg) {

	path, err := filepath.Abs(file); 
	if err != nil {
		syslog.Fatalf("readCfgFromFile failed from file '%s':\n %s", file, err)
	}

	ext := filepath.Ext(path)
	if len(ext) > 1 {
		ext = ext[1:]
	} else {
		syslog.Fatalf("readCfgFromFile failed from file '%s':\n can not found ext in file like .yml .ini ...", file)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
	 	syslog.Fatalf("readCfgFromFile failed from file '%s':\n %s", file, err)
	}

	if cfgs, err = parsingCfgsFromStr(string(data), ext); err != nil {
		syslog.Fatalf("readCfgFromFile failed from file '%s':\n %s", file, err)
	}

	return cfgs
}

func parsingCfgsFromStr(content string, ext string) (*rootCfg, error) {

	cfgs := rootCfg{ map[string]*Cfg{} }

	v := viper.New()
	v.SetConfigType(ext)
	if err := v.ReadConfig(strings.NewReader(string(content))); err != nil {
		return nil, fmt.Errorf("parsing failed: %s", err)
	}

	rootObj := v.Get(cfgRootKey)
	if rootObj == nil {
		return &cfgs, nil
	}

	root, ok := rootObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type of '.%s' in config file", cfgRootKey)
	}

	// parsing default cfg
	tmpDfCfg, err := parsingCfg(&dfCfg, v, cfgRootKey, "")
	if err != nil{
		return nil, fmt.Errorf("parsing default cfg failed:\n %s", err)
	}

	// parsing cfgs 
	for key, val := range root {
		if _, ok := val.(map[string]interface{}); !ok{
			continue
		}
		tmpCfg, err := parsingCfg(tmpDfCfg, v, cfgRootKey, key)
		if err != nil {
			return nil, fmt.Errorf("parsing cfg for '%s' failed:\n %s", key, err)
		}
		
		cfgs.Cfgs[key] = tmpCfg
	}

	dfCfg = *tmpDfCfg
	dfCfg.name = cfgDefaultName

	return &cfgs, nil
}

func parsingCfg(dfCfg *Cfg, v *viper.Viper, rootKey string, curKey string) (cfg *Cfg, err error) {

	if curKey != "" { curKey += "." }

	cfg = dfCfg.Clone(curKey)

	dirPath 		     	:= rootKey + "." + curKey + "dir"
	groupPath 	     	:= rootKey + "." + curKey + "group"
	filenamePath	    := rootKey + "." + curKey + "filename"
	maxsizePath       := rootKey + "." + curKey + "max_size"
	maxBackupsPath    := rootKey + "." + curKey + "max_backups"
	maxAgePath        := rootKey + "." + curKey + "max_age"
	compressPath      := rootKey + "." + curKey + "compress"
	levelPath         := rootKey + "." + curKey + "level"
	stackLevelPath    := rootKey + "." + curKey + "stack_level"
	colorPath         := rootKey + "." + curKey + "color"

	if v.IsSet(dirPath)  			 {cfg.Dir,     _, err = getMultiStringFromObj(v.Get(dirPath)     ); if err != nil {return nil, fmt.Errorf("parsing dir failed: %s", err)}}
	if v.IsSet(groupPath)			 {cfg.Group,   _, err = getMultiStringFromObj(v.Get(groupPath)   ); if err != nil {return nil, fmt.Errorf("parsing group failed: %s", err)}}
	if v.IsSet(filenamePath)	 {cfg.FileName,_, err = getMultiStringFromObj(v.Get(filenamePath)); if err != nil {return nil, fmt.Errorf("parsing filename failed: %s", err)}}
	if v.IsSet(maxsizePath)		 {cfg.MaxSize   , err = getIntFromObj(v.Get(maxsizePath)   ); if err != nil { return nil, fmt.Errorf("parsing max_size failed: %s", err) } }
	if v.IsSet(maxBackupsPath) {cfg.MaxBackups, err = getIntFromObj(v.Get(maxBackupsPath)); if err != nil { return nil, fmt.Errorf("parsing max_backups failed: %s", err) } }
	if v.IsSet(maxAgePath)		 {cfg.MaxAge    , err = getIntFromObj(v.Get(maxAgePath)    ); if err != nil { return nil, fmt.Errorf("parsing max_age failed: %s", err) } }
	if v.IsSet(compressPath)	 {cfg.Compress,_, err = getMultiBoolFromObj(v.Get(compressPath)); if err != nil { return nil, fmt.Errorf("parsing compress failed: %s", err) } }

	var errs []string
	
	if v.IsSet(levelPath) {
		if cfg.ConsoleLevel, cfg.FileLevel, err = getMultiStringFromObj(v.Get(levelPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing level failed: %s", err))
		}
	}
	if v.IsSet(stackLevelPath) {
		if cfg.ConsoleStackLevel, cfg.FileStackLevel, err = getMultiStringFromObj(v.Get(stackLevelPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing stack_level failed: %s", err))
		}
	}
	if v.IsSet(colorPath) {
		if cfg.ConsoleColor, cfg.FileColor, err = getMultiBoolFromObj(v.Get(colorPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing color failed: %s", err))
		}
	}

	if len(errs) > 0 {
		err = fmt.Errorf("%s", strings.Join(errs, " | "))
		return nil, err
	}
	if err = cfg.checkAndValidate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getIntFromObj(obj interface{}) (val int, err error) {

	if v, ok := obj.(int); ok {
		return v, nil
	}
	return 0, fmt.Errorf("invalid type: %s", reflect.TypeOf(obj))
}

func getMultiStringFromObj(obj interface{}) (s1 string, s2 string, err error) {
	switch obj := obj.(type) {
  case string: return obj, obj, nil
	case []interface{}: 
		if len(obj) == 0 { return "", "", fmt.Errorf("val not set")
	  } else {
			obj = append(obj, obj[0])
			var v1, v2 string; var ok bool
			if v1, ok = obj[0].(string); !ok {
				return "", "", fmt.Errorf("invalid type of [0]: %s", reflect.TypeOf(obj[0]))
			}
			if v2, ok = obj[1].(string); !ok {
				return "", "", fmt.Errorf("invalid type of [1]: %s", reflect.TypeOf(obj[0]))
			}
			return v1, v2, nil
		}
  }  

	return "", "", fmt.Errorf("invalid type: %s", reflect.TypeOf(obj))
}

func getMultiBoolFromObj(obj interface{}) (b1 bool, b2 bool, err error) {
	switch obj := obj.(type) {
  case bool: return obj, obj, nil
	case []interface{}: 
		if len(obj) == 0 { return false, false, fmt.Errorf("val not set")
	  } else {
			obj = append(obj, obj[0])
			var v1, v2 bool; var ok bool
			if v1, ok = obj[0].(bool); !ok {
				return false, false, fmt.Errorf("invalid type of [0]: %s", reflect.TypeOf(obj[0]))
			}
			if v2, ok = obj[1].(bool); !ok {
				return false, false, fmt.Errorf("invalid type of [1]: %s", reflect.TypeOf(obj[0]))
			}
			return v1, v2, nil
		}
  }  

	return false, false, fmt.Errorf("invalid type: %s", reflect.TypeOf(obj))
}

func getRepresentPathValue(path string, name string) string {

	if path == "" {
		return path
	}

	path = strings.ReplaceAll(path, "<HOSTNAME>", hostname)
	path = strings.ReplaceAll(path, "<APP_NAME>", appName)

	if name == cfgDefaultName {
		path = strings.ReplaceAll(path, "<LOG_NAME>", defaultFilename)
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