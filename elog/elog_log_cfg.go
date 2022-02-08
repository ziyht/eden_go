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
#    <APP> -> binary file name of current application
#    <LOG> -> the name of current logger, in default cfg, it will set to elog
#
#  note: 
#    1. the key like 'dir', 'group', ... under elog directly is to set default value,
#       you do not need to set it because all of them have a default value inside
#    2. for 'level', 'stack_level', 'color' you set only one val for console and file
#       settings, or you can set a arr with two val to set config for console and file
#       respectively 
#

elog:
  tag           : [ELOG]                  # default [ELOG]
  dir           : logs                    # default logs
  group         : <HOSTNAME>              # default <HOSTNAME>, if set, real dir will be $Dir/$Group
  filename      : <APP>_<LOG>             # default <APP>_<LOG>, will not write to file if set empty, real file path will be $Dir/$Group/$File
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
	cfgRootKey       = "elog"          // root key in the config file for elog
	dfDir            = "logs"
	dfGroup          = "<HOSTNAME>"  
	dfFileName       = "<APP>_<LOG>"
	dfTag            = "ELOG"
	dfFileCfgName    = "file"
	dfConsoleCfgName = "console"
)

var (
	dfCfg       = *genDfCfg()
	appName     = AppName()
	hostname, _ = os.Hostname()
)

// Cfg ...
type Cfg struct{
	Tag               string `yaml:"tag"`
	Dir          	    string `yaml:"dir"`
	Group             string `yaml:"group"`
	FileName          string `yaml:"filename"`
	MaxSize     	    int    `yaml:"max_size"`
	MaxBackups  	    int    `yaml:"max_backups"`
	MaxAge      	    int    `yaml:"max_age"`
	Compress    	    bool   `yaml:"compress"`
	FileLevel         string `yaml:"file_level"`
	FileColor         colorSwitch   `yaml:"file_color"`
	FileStackLevel    string `yaml:"file_stack_level"`
	ConsoleLevel      string `yaml:"console_level"`
	ConsoleColor      colorSwitch   `yaml:"console_color"`
	ConsoleStackLevel string `yaml:"console_stack_level"`

	logDir  string
	path    string
}

type LogCfg struct {
	Name   	    string      `yaml:"name"` 	
	Tag    	    string      `yaml:"tag"`
	Level  	    string      `yaml:"level"`
	StackLevel  string      `yaml:"stack_level"`
	Color       colorSwitch

	// for console settings
	Stream      string   `yaml:"stream"`
	
	// for file settings
	Dir         string   `yaml:"dir"`
	Group       string   `yaml:"group"`
	FileName    string   `yaml:"filename"`
	MaxSize     int      `yaml:"max_size"`
	MaxBackups  int      `yaml:"max_backups"`
	MaxAge      int      `yaml:"max_age"`
  Compress    bool     `yaml:"compress"`
}

func (c *Cfg)Clone() (*Cfg){
	out := *c
	return &out 
}

func genDfCfg() *Cfg {
	return &Cfg{
		Tag              : dfTag,
		Dir              : dfDir,
		Group            : dfGroup,
		FileName         : dfFileName,
		MaxSize          : 100,
		MaxBackups       : 7,
		MaxAge           : 7,
		ConsoleLevel     : LEVELS_INFO,
		ConsoleColor     : ColorAuto,
		ConsoleStackLevel: LEVELS_ERROR,
		FileLevel        : LEVELS_DEBUG,
		FileColor        : ColorAuto,
		FileStackLevel   : LEVELS_WARN,
		Compress         : false,
	}
}

func genDfFileLogCfg() *LogCfg {
	return &LogCfg{
		Name             : dfFileCfgName,
		Tag              : dfTag,
		Dir              : dfDir,
		Group            : dfGroup,
		FileName         : dfFileName,
		MaxSize          : 100,
		MaxBackups       : 7,
		MaxAge           : 7,
		Level            : LEVELS_DEBUG,
		Color            : ColorAuto,
		StackLevel       : LEVELS_WARN,
		Compress         : false,
	}
}

func genDfConsoleLogCfg() *LogCfg {
	return &LogCfg{
		Name             : dfConsoleCfgName,
		Tag              : dfTag,
		Dir              : dfDir,
		Group            : dfGroup,
		FileName         : dfFileName,
		MaxSize          : 100,
		MaxBackups       : 7,
		MaxAge           : 7,
		Level            : LEVELS_INFO,
		Color            : ColorAuto,
		StackLevel       : LEVELS_ERROR,
		Compress         : false,
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

func readCfgFromFile(file string) (cfgs map[string]*Cfg) {

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

func parsingCfgsFromStr(content string, ext string) (map[string]*Cfg, error) {

	cfgs := map[string]*Cfg{} 

	v := viper.New()
	v.SetConfigType(ext)
	if err := v.ReadConfig(strings.NewReader(string(content))); err != nil {
		return nil, fmt.Errorf("parsing failed: %s", err)
	}

	rootObj := v.Get(cfgRootKey)
	if rootObj == nil {
		return cfgs, nil
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
		
		cfgs[key] = tmpCfg
	}

	dfCfg = *tmpDfCfg

	return cfgs, nil
}

func parsingDfCfg(dfCfg *Cfg, v *viper.Viper, rootKey string, curKey string)(cfg *Cfg, err error) {

	if curKey != "" { curKey += "." }

	cfg = dfCfg.Clone()

	tagPath 		     	:= rootKey + "." + curKey + "tag"
	dirPath 		     	:= rootKey + "." + curKey + "dir"
	groupPath 	     	:= rootKey + "." + curKey + "group"
	filenamePath	    := rootKey + "." + curKey + "filename"
	maxsizePath       := rootKey + "." + curKey + "max_size"
	maxBackupsPath    := rootKey + "." + curKey + "max_backups"
	maxAgePath        := rootKey + "." + curKey + "max_age"
	compressPath      := rootKey + "." + curKey + "compress"
	fLevelPath        := rootKey + "." + curKey + "f_level"
	fStackLevelPath   := rootKey + "." + curKey + "f_stack_level"
	fColorPath        := rootKey + "." + curKey + "f_color"
	cLevelPath        := rootKey + "." + curKey + "c_level"
	cStackLevelPath   := rootKey + "." + curKey + "c_stack_level"
	cColorPath        := rootKey + "." + curKey + "c_color"

	if v.IsSet(tagPath)  			 {cfg.Tag,     _, err = getMultiStringFromObj(v.Get(tagPath)     ); if err != nil {return nil, fmt.Errorf("parsing tag failed: %s", err)}}
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
	if v.IsSet(fColorPath) {
		if cfg.FileColor, err = parsingColorStr(v.Get(fColorPath)); err != nil {
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

func parsingCfg(dfCfg *Cfg, v *viper.Viper, rootKey string, curKey string) (cfg *Cfg, err error) {

	if curKey != "" { curKey += "." }

	cfg = dfCfg.Clone()

	tagPath 		     	:= rootKey + "." + curKey + "tag"
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

	if v.IsSet(tagPath)  			 {cfg.Tag,     _, err = getMultiStringFromObj(v.Get(tagPath)     ); if err != nil {return nil, fmt.Errorf("parsing tag failed: %s", err)}}
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


func parsingColorStr(obj interface{}) (val colorSwitch, err error) {

	if v, ok := obj.(bool); ok {
		if v {
			return ColorOn, nil
		} else {
			return ColorOff, nil
		}
	}
  if v, ok := obj.(string); ok {
    v = strings.ToLower(v)
		switch v {
      case "true" : return ColorOn, nil
      case "false": return ColorOff, nil
      case "auto" : return ColorAuto, nil
		}

		return ColorOff, fmt.Errorf("invalid str: %s, you can set 'true', 'false', 'auto'", v)
  }
	return ColorOff, fmt.Errorf("invalid type: %s", reflect.TypeOf(obj))
}


func getRepresentPathValue(path string, name string) string {

	if path == "" {
		return path
	}

	if name == "" {
		name = dfFileName
	}

	path = strings.ReplaceAll(path, "<HOSTNAME>", hostname)
	path = strings.ReplaceAll(path, "<APP_NAME>", appName)
	path = strings.ReplaceAll(path, "<LOG_NAME>", name)

	return path
}

func AppName() string {
	appName, _ := os.Executable()
	appName = filepath.Base(appName)
	appName = strings.TrimSuffix(appName, ".exe")

	return appName
}