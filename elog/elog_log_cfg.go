package elog

import (
	"fmt"
	"io/ioutil"
	"os"
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
	Tag               string        `yaml:"tag"`
	Dir          	    string        `yaml:"dir"`
	Group             string        `yaml:"group"`
	FileName          string        `yaml:"filename"`
	MaxSize     	    int           `yaml:"max_size"`
	MaxBackup   	    int           `yaml:"max_backup"`
	MaxAge      	    int           `yaml:"max_age"`
	Compress    	    bool          `yaml:"compress"`
	FileLevel         LevelEnabler  `yaml:"file_level"`
	FileColor         colorSwitch   `yaml:"file_color"`
	FileStackLevel    LevelEnabler  `yaml:"file_stack_level"`
	ConsoleLevel      LevelEnabler  `yaml:"console_level"`
	ConsoleColor      colorSwitch   `yaml:"console_color"`
	ConsoleStackLevel LevelEnabler  `yaml:"console_stack_level"`

	logDir  string
	path    string
}

type LogCfg struct {
	Name   	    string       `yaml:"name"` 	
	Tag    	    string       `yaml:"tag"`
	Level  	    LevelEnabler `yaml:"level"`
	StackLevel  LevelEnabler `yaml:"stack_level"`
	Color       colorSwitch 

	// for console settings
	Stream      string   `yaml:"stream"`
	
	// for file settings
	Dir         string   `yaml:"dir"`
	Group       string   `yaml:"group"`
	FileName    string   `yaml:"filename"`
	MaxSize     int      `yaml:"max_size"`
	MaxBackup   int      `yaml:"max_backup"`
	MaxAge      int      `yaml:"max_age"`
  Compress    bool     `yaml:"compress"`
}

type LoggerCfg struct {
	df   *Cfg
	cfgs []*LogCfg
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
		MaxBackup        : 7,
		MaxAge           : 7,
		Level            : LEVEL_INFO,
		Color            : ColorAuto,
		StackLevel       : LEVEL_ERROR,
		Compress         : false,
	}
}

func (cfg *Cfg)check() (err error) {
	if err = cfg.checkFileRotate(); err != nil {return}
	return cfg.checkLevelStr()
}

func (cfg *Cfg)genLogCfg()*LogCfg {
	return &LogCfg{
		Name             : "",
		Tag              : cfg.Tag,
		Dir              : cfg.Dir,
		Group            : cfg.Group,
		FileName         : cfg.FileName,
		MaxSize          : cfg.MaxSize,
		MaxBackup        : cfg.MaxBackup,
		MaxAge           : cfg.MaxAge,
		Level            : LEVEL_NONE,
		Color            : ColorAuto,
		StackLevel       : LEVEL_OFF,
		Compress         : false,	
	}
}

func (cfg *Cfg)checkLevelStr() (err error) {
	// if _, err = getLevelByStr(cfg.ConsoleLevel);      err != nil {return}
	// if _, err = getLevelByStr(cfg.FileLevel);         err != nil {return}
	// if _, err = getLevelByStr(cfg.ConsoleStackLevel); err != nil {return}
	// if _, err = getLevelByStr(cfg.FileStackLevel);    err != nil {return}

	return
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

func readCfgFromFile(file string) (cfgs map[string]*LoggerCfg) {

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
	tmpDfCfg, err := parsingDfCfg(&dfCfg, v, cfgRootKey, "")
	if err != nil{
		return nil, fmt.Errorf("parsing default cfg failed:\n %s", err)
	}

	// parsing cfgs 
	for key, val := range root {
		if _, ok := val.(map[string]interface{}); !ok{
			continue
		}
		tmpCfg, err := parsingLoggerCfg(tmpDfCfg, v, cfgRootKey, key)
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
	if v.IsSet(maxBackupsPath) {cfg.MaxBackup , err = getIntFromObj(v.Get(maxBackupsPath)); if err != nil { return nil, fmt.Errorf("parsing max_backups failed: %s", err) } }
	if v.IsSet(maxAgePath)		 {cfg.MaxAge    , err = getIntFromObj(v.Get(maxAgePath)    ); if err != nil { return nil, fmt.Errorf("parsing max_age failed: %s", err) } }
	if v.IsSet(compressPath)	 {cfg.Compress,_, err = getMultiBoolFromObj(v.Get(compressPath)); if err != nil { return nil, fmt.Errorf("parsing compress failed: %s", err) } }

	var errs []string
	
	if v.IsSet(fLevelPath) {
		if cfg.FileLevel, err = parsingLevel(v.Get(fLevelPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing f_level failed: %s", err))
		}
	}
	if v.IsSet(fStackLevelPath) {
		if cfg.FileStackLevel, err = parsingLevel(v.Get(fStackLevelPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing f_stack_level failed: %s", err))
		}
	}
	if v.IsSet(fColorPath) {
		if cfg.FileColor, err = parsingColor(v.Get(fColorPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing f_color failed: %s", err))
		}
	}
	if v.IsSet(cLevelPath) {
		if cfg.FileLevel, err = parsingLevel(v.Get(cLevelPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing c_level failed: %s", err))
		}
	}
	if v.IsSet(cStackLevelPath) {
		if cfg.FileStackLevel, err = parsingLevel(v.Get(cStackLevelPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing c_stack_level failed: %s", err))
		}
	}
	if v.IsSet(cColorPath) {
		if cfg.FileColor, err = parsingColor(v.Get(cColorPath)); err != nil {
			errs = append(errs, fmt.Sprintf("parsing c_color failed: %s", err))
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

func parsingLoggerCfg(dfCfg *Cfg, v *viper.Viper, rootKey string, curKey string) (cfg *LoggerCfg, err error) {

	cfg = &LoggerCfg{df: dfCfg}

	root := v.Get(rootKey + "." + curKey)

	switch obj := root.(type) {

	case []interface{}: 
		for idx, iter := range obj {
			curObj, ok := iter.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid val type in %s.%s[%d], need a map[string]interface{}", rootKey, curKey, idx)
			}
			logCfg, err := parsingLogCfg(dfCfg, curObj)
			if err != nil {
				return nil, err
			}
			cfg.cfgs = append(cfg.cfgs, logCfg)
		}

	case map[string]interface{}:
		logCfg, err := parsingLogCfg(dfCfg, obj)
		if err != nil {
			return nil, err
		}
		cfg.cfgs = append(cfg.cfgs, logCfg)
	}

	if len(cfg.cfgs) == 0 {
		return nil, fmt.Errorf("can not found any log configs")
	}

	return cfg, nil
}

func parsingLogCfg(dfCfg *Cfg, obj map[string]interface{})(cfg *LogCfg, err error){

	cfg = dfCfg.genLogCfg()

	tagPath 		     	:= "tag"
	dirPath 		     	:= "dir"
	groupPath 	     	:= "group"
	filenamePath	    := "filename"
	maxsizePath       := "max_size"
	maxBackupPath     := "max_backup"
	maxAgePath        := "max_age"
	compressPath      := "compress"
	levelPath         := "level"
	stackLevelPath    := "stack_level"
	colorPath         := "color"
	streamPath        := "stream"

	if tag,       ok := obj[tagPath];       ok {cfg.Tag,      _, err = getMultiStringFromObj(tag     ); if err != nil {return nil, fmt.Errorf("parsing tag failed: %s", err)}}
	if dir,       ok := obj[dirPath];       ok {cfg.Dir,      _, err = getMultiStringFromObj(dir     ); if err != nil {return nil, fmt.Errorf("parsing dir failed: %s", err)}}
	if group,     ok := obj[groupPath];     ok {cfg.Group,    _, err = getMultiStringFromObj(group   ); if err != nil {return nil, fmt.Errorf("parsing group failed: %s", err)}}
	if filename,  ok := obj[filenamePath];  ok {cfg.FileName, _, err = getMultiStringFromObj(filename); if err != nil {return nil, fmt.Errorf("parsing filename failed: %s", err)}}
 	if maxsize,   ok := obj[maxsizePath];   ok {cfg.MaxSize,     err = getIntFromObj(maxsize  );        if err != nil {return nil, fmt.Errorf("parsing max_size failed: %s", err)}} 
 	if maxBackup, ok := obj[maxBackupPath]; ok {cfg.MaxBackup,   err = getIntFromObj(maxBackup);        if err != nil {return nil, fmt.Errorf("parsing max_backup failed: %s", err)}} 
 	if maxAge,    ok := obj[maxAgePath];    ok {cfg.MaxAge,      err = getIntFromObj(maxAge   );        if err != nil {return nil, fmt.Errorf("parsing max_age failed: %s", err)}} 
	if compress,  ok := obj[compressPath];  ok {cfg.Compress, _, err = getMultiBoolFromObj(compress );  if err != nil {return nil, fmt.Errorf("parsing compress failed: %s", err)}} 
	if level,     ok := obj[levelPath];     ok {cfg.Level,       err = parsingLevel(level);             if err != nil {return nil, fmt.Errorf("parsing level failed: %s", err)}}
	if slevel,    ok := obj[stackLevelPath];ok {cfg.StackLevel,  err = parsingLevel(slevel);            if err != nil {return nil, fmt.Errorf("parsing stack_level failed: %s", err)}}
	if color,     ok := obj[colorPath];     ok {cfg.Color,       err = parsingColor(color);             if err != nil {return nil, fmt.Errorf("parsing color failed: %s", err)}}
	if stream,    ok := obj[streamPath];    ok {cfg.Stream,      err = parsingStream(stream);           if err != nil {return nil, fmt.Errorf("parsing stream failed: %s", err)}}

	return cfg, nil
}

type LevelEnablerFunc func(Level) bool

// Enabled calls the wrapped function.
func (f LevelEnablerFunc) Enabled(lvl Level) bool { return f(lvl) }

func parsingLevel(obj interface{}) (val LevelEnabler, err error) {

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

	return  LevelEnablerFunc(func(l Level) bool {
		if l >= l1l && l <= l2l {
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