package elog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
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
  dir           : logs                # default logs
  group         : <HOSTNAME>          # default <HOSTNAME>, if set, real dir will be $Dir/$Group
  filename      : <APP>_<LOG>         # default <LOG_NAME>, will not write to file if set empty, real file path will be $Dir/$Group/$File
  max_size      : 100                 # default 100, unit MB
  max_backups   : 7                   # default 7
  max_age       : 7                   # default 7
  compress      : false               # default false
  f_level       : debug               # default debug, for file, valid value is [debug, info, warn, error, fatal, panic]
  f_stack_level : warn                # default warn , for file, valid value is [debug, info, warn, error, fatal, panic]
  f_color       : false               # default false, for file
  c_level       : info                # default info , for console, valid value is [debug, info, warn, error, fatal, panic]
  c_stack_level : error               # default error, for console, valid value is [debug, info, warn, error, fatal, panic]
  c_color       : true                # default true , for console

  default:
  - name        : console
    stream      : stdout
    level       : info
    stack_level : error
    color       : auto                # auto[default], true, false
  - name        : file
    dir         : logs                # default logs
    group       : <HOSTNAME>          # default <HOSTNAME>, if set, real dir will be $dir/$group
    filename    : <APP>_<LOG>         # default <LOG_NAME>, will not write to file if set empty, real file path will be $dir/$group/$file_name
    max_size    : 100                 # default 100, unit MB
    max_backup  : 7                   # default 7
    max_age     : 7                   # default 7
    compress    : false               # default false
    level       : debug               # default debug, for file, valid value is [debug, info, warn, error, fatal, panic]
    stack_level : warn                # default warn , for file, valid value is [debug, info, warn, error, fatal, panic]
    color       : false               # default false, for file
`

const (
	cfgRootKey       = "elog"          // root key in the config file for elog
	dfDir            = "logs"
	dfGroup          = "<HOSTNAME>"  
	dfFileName       = "<APP>_<LOG>"
	dfTag            = ""
	dfFileCfgName    = "file"
	dfConsoleCfgName = "console"
)

var (
	dfCfg       = *genDfCfg()
	appName     = AppName()
	hostname, _ = os.Hostname()
)

const (
	inheritKey       = "inherit"
	nameKey          = "name"
	tagKey 		     	 = "tag"
	dirKey 		     	 = "dir"
	groupKey 	     	 = "group"
	filenameKey	     = "filename"
	maxsizeKey       = "max_size"
	maxBackupKey     = "max_backup"
	maxAgeKey        = "max_age"
	compressKey      = "compress"
	levelKey         = "level"
	stackLevelKey    = "stack_level"
	colorKey         = "color"
	consoleKey       = "console"

	c_levelKey       = "c_level"
	c_stackLevelKey  = "c_stack_level"
	c_colorKey       = "c_color"
	f_levelKey       = "f_levle"
	f_stackLevelKey  = "f_stack_level"
	f_colorKey       = "f_color"
)

var skipKeys = map[string]bool{
		inheritKey     : true,
		nameKey        : true,
		tagKey         : true,
		dirKey         : true,
		groupKey       : true,
		filenameKey    : true,
		maxsizeKey     : true,
		maxBackupKey   : true,
		maxAgeKey      : true,
		compressKey    : true,
		c_levelKey     : true,
		c_stackLevelKey: true,
		c_colorKey     : true,
		f_levelKey     : true,
		f_stackLevelKey: true,
		f_colorKey     : true,
	}

// Cfg ...
type Cfg struct{
	Tag               string
	Dir          	    string
	Group             string
	FileName          string
	MaxSize     	    int
	MaxBackup   	    int
	MaxAge      	    int
	Compress    	    bool
	FileLevel         zapcore.LevelEnabler
	FileColor         colorSwitch
	FileStackLevel    zapcore.LevelEnabler
	ConsoleLevel      zapcore.LevelEnabler
	ConsoleColor      colorSwitch
	ConsoleStackLevel zapcore.LevelEnabler

	// logDir  string
	// path    string
}


func (c *Cfg)Clone() (*Cfg){
	out := *c
	return &out 
}

type LogCfg struct {
  inherit     string
	file        bool     // is current cfg for file settings

	Name   	    string
	Tag    	    string 
	Level  	    zapcore.LevelEnabler
	StackLevel  zapcore.LevelEnabler
	Color       colorSwitch 

	// for console settings
	Console     string
	
	// for file settings
	Dir         string
	Group       string
	FileName    string
	MaxSize     int
	MaxBackup   int
	MaxAge      int
  Compress    bool
}
func (c *LogCfg)Clone() (*LogCfg){
	out := *c
	return &out 
}

type LoggerCfg struct {
	df   *Cfg
	logs []*LogCfg
}
func (c *LoggerCfg)Clone() (*LoggerCfg){
	out := *c

	for i, cfg := range c.logs {
		c.logs[i] = cfg.Clone()
	}

	return &out 
}
func (c *LoggerCfg)FindLogCfg(name string) *LogCfg {
	for _, cfg := range c.logs{
		if name == cfg.Name{	
			return cfg.Clone()
		}
	}

	return nil
}
func (c *LoggerCfg)validateAndCheck() error {

	return nil
}

func genDfCfg() *Cfg {
	return &Cfg{
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
}

func genDfFileLogCfg() *LogCfg {
	return &LogCfg{
		Name             : dfFileCfgName,
		Tag              : dfTag,
		Dir              : dfDir,
		Group            : dfGroup,
		FileName         : dfFileName,
		MaxSize          : 100,
		MaxBackup        : 7,
		MaxAge           : 7,
		Level            : LEVEL_DEBUG,
		Color            : ColorAuto,
		StackLevel       : LEVEL_WARN,
		Compress         : false,
		Console           : "",
	}
}

func genDfConsoleLogCfg() *LogCfg {
	return &LogCfg{
		Name             : dfConsoleCfgName,
		Tag              : dfTag,
		Dir              : dfDir,
		Level            : LEVEL_INFO,
		Color            : ColorAuto,
		StackLevel       : LEVEL_ERROR,
		Console          : "stdout",
	}
}

func (cfg *Cfg)check() (err error) {
	return cfg.checkFileRotate();
}

func (cfg *Cfg)genLoggerCfg()*LoggerCfg {

	out := &LoggerCfg{}

	out.logs = append(out.logs, &LogCfg{
		file             : true,
		Name             : dfFileCfgName,
		Tag              : cfg.Tag,
		Dir              : cfg.Dir,
		Group            : cfg.Group,
		FileName         : cfg.FileName,
		MaxSize          : cfg.MaxSize,
		MaxBackup        : cfg.MaxBackup,
		MaxAge           : cfg.MaxAge,
		Level            : cfg.FileLevel,
		Color            : cfg.FileColor,
		StackLevel       : cfg.FileStackLevel,
		Compress         : cfg.Compress,
	})
	out.logs = append(out.logs, &LogCfg{
		file             : false,
		Name             : dfConsoleCfgName,
		Tag              : cfg.Tag,
		Console           : "stdout", 
		Level            : cfg.ConsoleLevel,
		Color            : cfg.ConsoleColor,
		StackLevel       : cfg.ConsoleStackLevel,
	})

	return out
}

func (cfg *Cfg)genLogCfg()[]*LogCfg {
	return cfg.genLoggerCfg().logs
}

func (cfg *Cfg)checkFileRotate() (err error) {
	if cfg.MaxSize < 0 {
		return fmt.Errorf("invalid max_size(%d), should >= 0", cfg.MaxSize)
	}
	if cfg.MaxBackup < 0 {
		return fmt.Errorf("invalid max_backup(%d), should >= 0", cfg.MaxBackup)
	}
	if cfg.MaxAge < 0 {
		return fmt.Errorf("invalid max_age(%d), should >= 0", cfg.MaxAge)
	}

	return
}

func (cfg *Cfg)validate() (err error){

	if cfg.MaxSize    < 0 { cfg.MaxSize    = 0 }
	if cfg.MaxBackup  < 0 { cfg.MaxBackup  = 0 }
	if cfg.MaxAge     < 0 { cfg.MaxAge     = 0 }

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

func parsingCfgsFromFile(file string) (cfgs map[string]*LoggerCfg) {

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

// parsingCfgsFromStr
// ext - file extension or content type
func parsingCfgsFromStr(content string, ext string) (map[string]*LoggerCfg, error) {

	cfgs := map[string]*LoggerCfg{} 

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
	tmpDfCfg, err := parsingDfCfg(&dfCfg, root)
	if err != nil{
		return nil, fmt.Errorf("parsing default cfg failed:\n %s", err)
	}

	// parsing cfgs
	news := make(map[string]*LoggerCfg)
	for key, _ := range root {
		if skipKeys[key] {
			continue
		}
		tmpCfg, err := parsingLoggerCfg(tmpDfCfg, v, cfgRootKey, key, news)
		if err != nil {
			return nil, fmt.Errorf("parsing cfg for '%s' failed:\n %s", key, err)
		}
		
		cfgs[key] = tmpCfg
	}

	dfCfg = *tmpDfCfg

	return cfgs, nil
}

func parsingDfCfg(dfCfg *Cfg, root map[string]interface{})(cfg *Cfg, err error) {

	cfg = dfCfg.Clone()

	var level, stackLevel zapcore.LevelEnabler
	var color colorSwitch

	_, err = parsingInStrFromMap( &cfg.Dir      , root, dirKey      , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", dirKey, err)}
	_, err = parsingInStrFromMap( &cfg.Group    , root, groupKey    , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", groupKey, err)}
	_, err = parsingInIntFromMap( &cfg.MaxSize  , root, maxsizeKey  , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", maxsizeKey, err)}
	_, err = parsingInIntFromMap( &cfg.MaxBackup, root, maxBackupKey, nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", maxBackupKey, err)}
	_, err = parsingInIntFromMap( &cfg.MaxAge   , root, maxAgeKey   , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", maxAgeKey, err)}
	_, err = parsingInBoolFromMap(&cfg.Compress , root, compressKey , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", compressKey, err)}
	_, err = parsingInStrFromMap( &cfg.Tag      , root, tagKey      , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", tagKey, err)}

	lok, err := parsingInLevelFromMap(&level     , root, levelKey     , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", levelKey, err)}
	sok, err := parsingInLevelFromMap(&stackLevel, root, stackLevelKey, nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", stackLevelKey, err)}
	cok, err := parsingInColorFromMap(&color     , root, colorKey     , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", colorKey, err)}

	if lok { cfg.FileLevel      = level;      cfg.ConsoleLevel      = level     }
	if sok { cfg.FileStackLevel = stackLevel; cfg.ConsoleStackLevel = stackLevel}
	if cok { cfg.FileColor      = color;      cfg.ConsoleColor      = color     }

	_, err = parsingInLevelFromMap(&cfg.FileLevel     , root, f_levelKey     , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", f_levelKey, err)}
	_, err = parsingInLevelFromMap(&cfg.FileStackLevel, root, f_stackLevelKey, nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", f_stackLevelKey, err)}
	_, err = parsingInColorFromMap(&cfg.FileColor     , root, f_colorKey     , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", f_colorKey, err)}	

	_, err = parsingInLevelFromMap(&cfg.ConsoleLevel     , root, c_levelKey     , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", c_levelKey, err)}
	_, err = parsingInLevelFromMap(&cfg.ConsoleStackLevel, root, c_stackLevelKey, nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", c_stackLevelKey, err)}
	_, err = parsingInColorFromMap(&cfg.ConsoleColor     , root, c_colorKey     , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", c_colorKey, err)}

	if err = cfg.checkAndValidate(); err != nil {
		return nil, err
	}

	return cfg, nil

}

func parsingLoggerCfg(dfCfg *Cfg, v *viper.Viper, rootKey string, curKey string, news map[string]*LoggerCfg) (cfg *LoggerCfg, err error) {

	cfg = &LoggerCfg{df: dfCfg}

	root := v.Get(rootKey + "." + curKey)

	switch obj := root.(type) {

	case []interface{}: 
		for idx, iter := range obj {
			tmp, ok := iter.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid val type in %s.%s[%d], need a map[string]interface{}", rootKey, curKey, idx)
			}
			curObj := convert(tmp)
			logCfg, err := parsingLogCfg(dfCfg, curObj, news)
			if err != nil {
				return nil, err
			}
			cfg.logs = append(cfg.logs, logCfg...)
		}

	case map[string]interface{}:
		logCfg, err := parsingLogCfg(dfCfg, obj, news)
		if err != nil {
			return nil, err
		}
		cfg.logs = append(cfg.logs, logCfg...)
	}

	if len(cfg.logs) == 0 {
		return nil, fmt.Errorf("can not found any log configs")
	}

	if news != nil {
		news[curKey] = cfg
	}

	return cfg, nil
}

func convert(in map[interface{}]interface{})map[string]interface{} {

	out := make(map[string]interface{})
	for key, v := range in {
		out[key.(string)] = v
	}
	return out
}

func parsingLogCfg(dfCfg *Cfg, root map[string]interface{}, news map[string]*LoggerCfg)(cfgs []*LogCfg, err error){

	cfg := &LogCfg{}

	set, err := parsingInStrFromMap(&cfg.inherit, root, inheritKey, nil); if err != nil { return nil, fmt.Errorf("parsing inherit failed: %s", err)}
	if set {
		cfgs = cfgMan.findLogCfgs(cfg.inherit, news)
		if cfgs == nil {
			return nil, fmt.Errorf("can not find needed inherit log '%s'", cfg.inherit)
		}
	}

	_, err =  parsingInStrFromMap(&cfg.FileName, root, filenameKey, nil); if err != nil { return nil, fmt.Errorf("parsing filename failed: %s", err)}
	_, err =  parsingInStrFromMap(&cfg.Console , root, consoleKey , nil); if err != nil { return nil, fmt.Errorf("parsing console failed: %s", err)}
	if cfg.FileName != "" && cfg.Console != "" {
		return nil, fmt.Errorf("invalid config checked, filename, console are both set: (%s, %s), you should set only one of them", cfg.FileName, cfg.Console)
	}

  cfgs = dfCfg.genLogCfg()
	var out []*LogCfg
	for _, df := range cfgs {

		if cfg.file {
			if cfg.FileName == ""{
				continue
			}
			_, err = parsingInStrFromMap(&cfg.Dir      , root, dirKey      , &df.Dir      ); if err != nil { return nil, fmt.Errorf("parsing dir failed: %s", err)}
		  _, err = parsingInStrFromMap(&cfg.Group    , root, groupKey    , &df.Group    ); if err != nil { return nil, fmt.Errorf("parsing group failed: %s", err)}
		  _, err = parsingInIntFromMap(&cfg.MaxSize  , root, maxsizeKey  , &df.MaxSize  ); if err != nil { return nil, fmt.Errorf("parsing max_size failed: %s", err)}
		  _, err = parsingInIntFromMap(&cfg.MaxBackup, root, maxBackupKey, &df.MaxBackup); if err != nil { return nil, fmt.Errorf("parsing max_backup failed: %s", err)}
		  _, err = parsingInIntFromMap(&cfg.MaxAge   , root, maxAgeKey   , &df.MaxAge)   ; if err != nil { return nil, fmt.Errorf("parsing max_age failed: %s", err)}
		  _, err = parsingInBoolFromMap(&cfg.Compress, root, compressKey , &df.Compress) ; if err != nil { return nil, fmt.Errorf("parsing compress failed: %s", err)}
		} else {
			if cfg.Console == ""{
				continue
			}
		}
		_, err = parsingInStrFromMap(  &cfg.Name      , root, nameKey      , &df.Name      ); if err != nil { return nil, fmt.Errorf("parsing name failed: %s", err)}
		_, err = parsingInStrFromMap(  &cfg.Tag       , root, tagKey       , &df.Tag       ); if err != nil { return nil, fmt.Errorf("parsing tag failed: %s", err)}
		_, err = parsingInLevelFromMap(&cfg.Level     , root, levelKey     , &df.Level     ); if err != nil { return nil, fmt.Errorf("parsing level failed: %s", err)}
		_, err = parsingInLevelFromMap(&cfg.StackLevel, root, stackLevelKey, &df.StackLevel); if err != nil { return nil, fmt.Errorf("parsing stack_level failed: %s", err)}
		_, err = parsingInColorFromMap(&cfg.Color     , root, colorKey     , &df.Color     ); if err != nil { return nil, fmt.Errorf("parsing color failed: %s", err)}
		
		out = append(out, df)
	}

	return out, nil
}

func parsingInStrFromMap(dest *string, root map[string]interface{}, key string, df *string)(bool, error) {

	if df != nil {
		*dest = *df
	}

	if obj, exist := root[key]; exist{
		val, _, err := getMultiStringFromObj(obj);

		if err != nil {
			return false, err
		}

		*dest = val
		return true, nil
	} 

	return false, nil
}

func parsingInIntFromMap(dest *int, root map[string]interface{}, key string, df *int)(bool, error) {

	if df != nil {
		*dest = *df
	}

	if obj, exist := root[key]; exist{
		val, err := getIntFromObj(obj);

		if err != nil {
			return false, err
		}

		*dest = val
		return true, nil
	} 

	return false, nil
}

func parsingInBoolFromMap(dest *bool, root map[string]interface{}, key string, df *bool)(bool, error) {

	if df != nil {
		*dest = *df
	}

	if obj, exist := root[key]; exist{
		val, _, err := getMultiBoolFromObj(obj);

		if err != nil {
			return false, err
		}

		*dest = val
		return true, nil
	} 

	return false, nil
}

func parsingInLevelFromMap(dest *zapcore.LevelEnabler, root map[string]interface{}, key string, df *zapcore.LevelEnabler)(bool, error) {

	if df != nil {
		*dest = *df
	}

	if obj, exist := root[key]; exist{
		val, err := parsingLevel(obj);

		if err != nil {
			return false, err
		}

		*dest = val
		return true, nil
	} 

	return false, nil
}

func parsingInColorFromMap(dest *colorSwitch, root map[string]interface{}, key string, df *colorSwitch)(bool, error) {

	if df != nil {
		*dest = *df
	}

	if obj, exist := root[key]; exist{
		val, err := parsingColor(obj);

		if err != nil {
			return false, err
		}

		*dest = val
		return true, nil
	} 

	return false, nil
}

func parsingLevel(obj interface{}) (val zapcore.LevelEnabler, err error) {

	l1, l2, err := getMultiStringFromObj(obj)
	if  err != nil {
		return LEVEL_NONE, err
	}

	l1l, err := getLevelByStr(l1)
	if err != nil {
		return LEVEL_NONE, err
	}

	if l2 == "" {
		 return l1l, nil
	}

	l2l, err := getLevelByStr(l2)
	if err != nil {
		return LEVEL_NONE, err
	}

	if l1l > l2l {
		return LEVEL_NONE, fmt.Errorf("invalid span from %s to %s", l1, l2)
	}

	return  LevelEnablerFunc(func(l zapcore.Level) bool {
		if Level(l) >= l1l && Level(l) <= l2l {
			return true
		}
		return false
	}), nil

}

func parsingColor(obj interface{}) (val colorSwitch, err error) {

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

var _validStreams =  []string{"stdout", "STDOUT", "stderr", "STDERR", "1", "2"}

func parsingStream(obj interface{})(val string, err error) {

	val, ok := obj.(string)
	if !ok {
		return 
	}

  for _, valid := range _validStreams {
		if val == valid{
			return val, nil
		}
	}

	return "", fmt.Errorf("unsupport stream '%s', you can set: %s", val, _validStreams)
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