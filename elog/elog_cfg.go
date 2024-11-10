package elog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

const (
	cfgRootKey       = "elog"          // root key in the config file for elog
)

var (
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
	consoleKey       = "console"
	maxSizeKey       = "max_size"
	maxBackupKey     = "max_backup"
	maxAgeKey        = "max_age"
	compressKey      = "compress"
	levelKey         = "level"
	stackLevelKey    = "slevel"
	colorKey         = "color"
	

	c_levelKey       = "c_level"
	c_stackLevelKey  = "c_slevel"
	c_colorKey       = "c_color"
	f_levelKey       = "f_level"
	f_stackLevelKey  = "f_slevel"
	f_colorKey       = "f_color"
)

var skipKeys = map[string]bool{
		inheritKey     : true,
		nameKey        : true,
		tagKey         : true,
		dirKey         : true,
		groupKey       : true,
		filenameKey    : true,
		consoleKey     : true,
		maxSizeKey     : true,
		maxBackupKey   : true,
		maxAgeKey      : true,
		compressKey    : true,
		colorKey       : true,
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
	cfgs []*LogCfg
}
func (c *LoggerCfg)Clone() (*LoggerCfg){
	out := *c
	out.cfgs = make([]*LogCfg, len(c.cfgs))

	for i, cfg := range c.cfgs {
		out.cfgs[i] = cfg.Clone()
	}

	return &out 
}
func (c *LoggerCfg)FindLogCfg(name string) *LogCfg {
	for _, cfg := range c.cfgs{
		if name == cfg.Name{	
			return cfg.Clone()
		}
	}

	return nil
}
func (c *LoggerCfg)validateAndCheck() error {

	return nil
}

func (cfg *Cfg)check() (err error) {
	return cfg.checkFileRotate();
}

func (cfg *Cfg)genLoggerCfg()*LoggerCfg {

	out := &LoggerCfg{}

	out.cfgs = append(out.cfgs, &LogCfg{
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
	out.cfgs = append(out.cfgs, &LogCfg{
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
	return cfg.genLoggerCfg().cfgs
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
		syslog.Warnf("readCfgFromFile failed from file '%s': %s", file, err)
	}

	ext := filepath.Ext(path)
	if len(ext) > 1 {
		ext = ext[1:]
	} else {
		syslog.Warnf("readCfgFromFile failed from file '%s': can not found ext in file like .yml .ini ...", file)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		syslog.Warnf("readCfgFromFile failed from file '%s': %s", file, err)
	}

	if cfgs, err = parsingCfgsFromStr(string(data), ext, cfgRootKey); err != nil {
		syslog.Warnf("readCfgFromFile failed from file '%s': %s", file, err)
	}

	return cfgs
}

// parsingCfgsFromStr
// ext - file extension or content type
func parsingCfgsFromStr(content string, ext string, rootKey string) (map[string]*LoggerCfg, error) {

	cfgs := map[string]*LoggerCfg{} 

	v := viper.New()
	v.SetConfigType(ext)
	if err := v.ReadConfig(strings.NewReader(string(content))); err != nil {
		return nil, fmt.Errorf("parsing failed: %s", err)
	}

	rootObj := v.Get(rootKey)
	if rootObj == nil {
		syslog.Warnf ("can not find key '%s', skipped parsingCfgs for elog", rootKey)
		return cfgs, nil
	}

	root, ok := rootObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type of '.%s' in config file", rootKey)
	}

	// parsing default cfg
	tmpDfCfg, err := parsingDfCfg(&dfLoggerCfg, root)
	if err != nil{
		return nil, fmt.Errorf("parsing default cfg failed:\n %s", err)
	}

	// parsing cfgs
	cache := make(map[string]*LoggerCfg)
	cache[""]        = tmpDfCfg.genLoggerCfg()
	cache["default"] = tmpDfCfg.genLoggerCfg()
	for key := range root {
		if skipKeys[key] {
			continue
		}
		tmpCfg, err := parsingLoggerCfg(tmpDfCfg, v, rootKey, key, cache)
		if err != nil {
			return nil, fmt.Errorf("parsing cfg for '%s' failed:\n %s", key, err)
		}
		
		cfgs[key] = tmpCfg
	}

	dfLoggerCfg = *tmpDfCfg

	return cfgs, nil
}

func getViperFromFile(file string)(*viper.Viper, error){
	path, err := filepath.Abs(file); 
	if err != nil {
		return nil, fmt.Errorf("readCfgFromFile failed from file '%s': %s", file, err)
	}

	ext := filepath.Ext(path)
	if len(ext) > 1 {
		ext = ext[1:]
	} else {
		return nil, fmt.Errorf("readCfgFromFile failed from file '%s': can not found ext in filename like %s", file, []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl"})
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readCfgFromFile failed from file '%s': %s", file, err)
	}

	return getViperFromString(string(data), ext)
}

func getViperFromString(content string, ext string)(*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(ext)
	if err := v.ReadConfig(strings.NewReader(string(content))); err != nil {
		return nil, fmt.Errorf("parsing failed: %s", err)
	}
	return v, nil
}

func parsingLoggerCfgFromFile(file string, keys string) (cfg *LoggerCfg, err error) {
	v, err := getViperFromFile(file)
	if err != nil {
		return nil, err
	}

	rootObj := v.Get(keys)
	if rootObj == nil {
		return nil, fmt.Errorf("keys '%s' can not be found in file '%s'", keys, file)
	}

	cfg, err = parsingLoggerCfgSmart(&dfLoggerCfg, rootObj, nil)
	if err != nil {
		return nil, fmt.Errorf("error occurred when parsing '%s' in file '%s': %s", keys, file, err)
	}

	return cfg, err
}

func parsingDfCfg(dfCfg *Cfg, root map[string]interface{})(cfg *Cfg, err error) {

	cfg = dfCfg.Clone()

	var level, stackLevel zapcore.LevelEnabler
	var color colorSwitch

	_, err = parsingInStrFromMap( &cfg.FileName , root, filenameKey , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", filenameKey, err)}
	_, err = parsingInStrFromMap( &cfg.Dir      , root, dirKey      , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", dirKey, err)}
	_, err = parsingInStrFromMap( &cfg.Group    , root, groupKey    , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", groupKey, err)}
	_, err = parsingInIntFromMap( &cfg.MaxSize  , root, maxSizeKey  , nil); if err != nil { return nil, fmt.Errorf("parsing %s failed: %s", maxSizeKey, err)}
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

func parsingLoggerCfg(df *Cfg, v *viper.Viper, rootKey string, curKey string, cache map[string]*LoggerCfg) (cfg *LoggerCfg, err error) {

	keyPath := rootKey
	if curKey != "" {
		keyPath = rootKey + "." + curKey
	}

	root := v.Get(keyPath)
	if root == nil {
		return nil, fmt.Errorf("can not found any settings in path '%s'", keyPath)
	}

	cfg, err = parsingLoggerCfgSmart(df, root, cache)
	if err != nil {
		return nil, fmt.Errorf("parsing failed in: %s.%s: %s", rootKey, curKey, err)
	}

	if cache != nil && curKey != ""{
		cache[curKey] = cfg
	}

	return cfg, nil
}

func parsingLoggerCfgSmart(df *Cfg, in interface{}, cache map[string]*LoggerCfg)(cfg *LoggerCfg, err error){

	cfg = &LoggerCfg{df: df}
	var cur map[string]interface{}
	switch obj := in.(type) {

	case []interface{}: 
		for idx, iter := range obj {
			switch tmp := iter.(type) {
				case map[interface{}]interface{}:
					cur, err = tryConvert(tmp)
					if err != nil {
						return nil, fmt.Errorf("invalid type in [%d]: %s", idx, err)
					}
				case map[string]interface{}:
					cur = tmp
				case map[string]string:
					cur = tryConvert2(tmp)
				default:
					return nil, fmt.Errorf("invalid type in [%d], need a type like map[string]interface{}", idx)
			}

			logCfg, err := parsingLogCfg(df, cur, cache)
			if err != nil {
				return nil, fmt.Errorf("parsing failed in [%d]: %s", idx, err)
			}
			cfg.cfgs = append(cfg.cfgs, logCfg...)
		}

		if len(cfg.cfgs) == 0 {
			return df.genLoggerCfg(), nil
			//return nil, fmt.Errorf("can not found any log configs")
		}

		return cfg, nil

	case map[string]interface{}:
		cur = obj
	case map[string]string:
		cur = tryConvert2(obj)
	default:
		return nil, fmt.Errorf("invalid type of input, need map[string]interface{} or []map[string]interface{}")
	}	
	
	dfCfg, err := parsingDfCfg(df, cur)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %s", err)
	}

	return dfCfg.genLoggerCfg(), nil
}



func tryConvert(in map[interface{}]interface{})(map[string]interface{}, error) {

	out := make(map[string]interface{})
	for key, v := range in {

		strKey, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("invalid key type of '%v': %t, need string", key, key)
		}

		out[strKey] = v
	}
	return out, nil
}
func tryConvert2(in map[string]string)(map[string]interface{}) {
	out := make(map[string]interface{})
	for key, v := range in {
		out[key] = v
	}
	return out
}

func parsingLogCfg(dfCfg *Cfg, root map[string]interface{}, cache map[string]*LoggerCfg)(cfgs []*LogCfg, err error){

	cfg := &LogCfg{}

	inherit, err := parsingInStrFromMap(&cfg.inherit, root, inheritKey, nil); if err != nil { return nil, fmt.Errorf("parsing inherit failed: %s", err)}
	if inherit {
		cfgs = cfgMan.findLogCfgs(cfg.inherit, cache)
		if cfgs == nil {
			return nil, fmt.Errorf("can not find needed inherit log '%s'", cfg.inherit)
		}
		for _, iter := range cfgs {
			if iter.file {
				cfg.FileName = iter.FileName
			} else {
				cfg.Console = iter.Console
			}
		}
	}

	_, err =  parsingInStrFromMap(&cfg.FileName, root, filenameKey, nil); if err != nil { return nil, fmt.Errorf("parsing filename failed: %s", err)}
	_, err =  parsingInStrFromMap(&cfg.Console , root, consoleKey , nil); if err != nil { return nil, fmt.Errorf("parsing console failed: %s", err)}
	if cfg.FileName != "" && cfg.Console != "" {
		return nil, fmt.Errorf("invalid config checked, filename, console are both set: (%s, %s), you should set only one of them", cfg.FileName, cfg.Console)
	}
	if cfg.FileName == "" && cfg.Console == "" {
		return nil, fmt.Errorf("invalid config checked, filename, console are both not set, you should set one of them")
	}

  cfgs = dfCfg.genLogCfg()
	var out []*LogCfg
	for _, df := range cfgs {

		if df.file {
			if cfg.FileName == ""{
				continue
			}
			cfg.file = true
			_, err = parsingInStrFromMap(&cfg.Dir      , root, dirKey      , &df.Dir      ); if err != nil { return nil, fmt.Errorf("parsing dir failed: %s", err)}
			_, err = parsingInStrFromMap(&cfg.Group    , root, groupKey    , &df.Group    ); if err != nil { return nil, fmt.Errorf("parsing group failed: %s", err)}
			_, err = parsingInIntFromMap(&cfg.MaxSize  , root, maxSizeKey  , &df.MaxSize  ); if err != nil { return nil, fmt.Errorf("parsing max_size failed: %s", err)}
			_, err = parsingInIntFromMap(&cfg.MaxBackup, root, maxBackupKey, &df.MaxBackup); if err != nil { return nil, fmt.Errorf("parsing max_backup failed: %s", err)}
			_, err = parsingInIntFromMap(&cfg.MaxAge   , root, maxAgeKey   , &df.MaxAge)   ; if err != nil { return nil, fmt.Errorf("parsing max_age failed: %s", err)}
			_, err = parsingInBoolFromMap(&cfg.Compress, root, compressKey , &df.Compress) ; if err != nil { return nil, fmt.Errorf("parsing compress failed: %s", err)}
		} else {
			if cfg.Console == ""{
				continue
			}
			err = checkStream(cfg.Console)
			if err != nil {
				return nil, err
			}
		}
		_, err = parsingInStrFromMap(  &cfg.Name      , root, nameKey      , &df.Name      ); if err != nil { return nil, fmt.Errorf("parsing name failed: %s", err)}
		_, err = parsingInStrFromMap(  &cfg.Tag       , root, tagKey       , &df.Tag       ); if err != nil { return nil, fmt.Errorf("parsing tag failed: %s", err)}
		_, err = parsingInLevelFromMap(&cfg.Level     , root, levelKey     , &df.Level     ); if err != nil { return nil, fmt.Errorf("parsing level failed: %s", err)}
		_, err = parsingInLevelFromMap(&cfg.StackLevel, root, stackLevelKey, &df.StackLevel); if err != nil { return nil, fmt.Errorf("parsing slevel failed: %s", err)}
		_, err = parsingInColorFromMap(&cfg.Color     , root, colorKey     , &df.Color     ); if err != nil { return nil, fmt.Errorf("parsing color failed: %s", err)}
		
		out = append(out, cfg)
	}

	return out, nil
}

func parsingInStrFromMap(dest *string, root map[string]interface{}, key string, df *string)(bool, error) {

	if df != nil {
		*dest = *df
	}

	if obj, exist := root[key]; exist{
		strs, err := getMultiStringFromObj(obj);

		if err != nil {
			return false, err
		}

		if len(strs) == 0 {
			return false, fmt.Errorf("no strings found")
		}

		*dest = strs[0]
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

	lvs, err := getMultiStringFromObj(obj)
	if  err != nil {
		return LEVEL_NONE, err
	}

	if len(lvs) == 0 {
		return LEVEL_NONE, fmt.Errorf("level not found")
	}

	var lvls []Level
	for idx, lv := range lvs {
		lvl, err := getLevelByStr(lv)
		if err != nil {
			return LEVEL_NONE, fmt.Errorf("invalid level set in [%d]: %s", idx, err)
		}
		lvls = append(lvls, lvl)
	}

	if len(lvls) == 1 {
		return lvls[0], nil
	} else if len(lvls) == 2 {
		if lvls[0] > lvls[1] {
			return LEVEL_NONE, fmt.Errorf("invalid span from %s to %s", lvs[0], lvs[1])
		}
		lvl1 := lvls[0]
		lvl2 := lvls[1]
		return  LevelEnablerFunc(func(l zapcore.Level) bool {
			if Level(l) >= lvl1 && Level(l) <= lvl2 {
				return true
			}
			return false
		}), nil
	}
  
	return LevelEnablerFunc(func(l zapcore.Level) bool {
		for _, lvl := range lvls {
			if Level(l) == lvl {
				return true
			}
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

var validStreams = []string{"stdout", "STDOUT", "stderr", "STDERR", "1", "2"}

func checkStream(stream string)(err error) {

  for _, valid := range validStreams {
		if stream == valid{
			return nil
		}
	}

	return fmt.Errorf("unsupport stream '%s', you can set: %s", stream, validStreams)
}

func getRepresentPathValue(path string, name string) string {

	if path == "" {
		return path
	}

	path = strings.ReplaceAll(path, "<HOSTNAME>", hostname)
	path = strings.ReplaceAll(path, "<HOST>", hostname)
	path = strings.ReplaceAll(path, "<APPNAME>", appName)
	path = strings.ReplaceAll(path, "<APP_NAME>", appName)
	path = strings.ReplaceAll(path, "<APP>", appName)
	path = strings.ReplaceAll(path, "<LOGGERNAME>", name)
	path = strings.ReplaceAll(path, "<LOGGER>", name)
	path = strings.ReplaceAll(path, "<LOGNAME>", name)
	path = strings.ReplaceAll(path, "<LOG>", name)
	
	return path
}

func AppName() string {
	appName, _ := os.Executable()
	appName = filepath.Base(appName)
	appName = strings.TrimSuffix(appName, ".exe")

	return appName
}