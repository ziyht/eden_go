package elog

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
RootKey       = "elog"
DefaultCfgKey = "__default"
)

// ElogCfg ...
type Cfg struct{
	Dir          	string `yaml:"dir"`
	Group           string `yaml:"group"`
	FileName        string `yaml:"filename"`
	MaxSize     	int    `yaml:"max_size"`
	MaxBackups  	int    `yaml:"max_backups"`
	MaxAge      	int    `yaml:"max_age"`
	Compress    	bool   `yaml:"compress"`
	ConsoleLevel    string `yaml:"console_level"`
	ConsoleColor    bool   `yaml:"console_color"`
	FileLevel       string `yaml:"file_level"`
	FileColor       bool   `yaml:"file_color"`
	StackLevel      string `yaml:"stack_level"`

	consoleLevel    zapcore.Level
	fileLevel       zapcore.Level
	stackLevel      zapcore.Level
}

func baseDfCfg() Cfg {
	return Cfg{
		Dir         : "var/log",
		Group       : "<HOSTNAME>",
		FileName    : "elog",
		MaxSize     : 100,
		MaxBackups  : 7,
		MaxAge      : 7,
		Compress    : true,
		ConsoleLevel: "debug",
		ConsoleColor: true,
		FileLevel   : "debug",
		FileColor   : false,
		StackLevel  : "warn",
	}
}

type Cfgs struct {
	Cfgs map[string]*Cfg `yaml:"elog"`
}

var dfCfg *Cfg

func initDfCfg(){
	if dfCfg == nil {
		cfg := baseDfCfg()
		dfCfg = &cfg
	}
}

func readCfgFromYaml(file string) *Cfgs {

	var cfgs Cfgs

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

	gfcfg, _ := gjson.Load(file_, false)
	baseDfCfg := baseDfCfg()
	// parsing __default cfg
	{
		if _, ok := cfgs.Cfgs[DefaultCfgKey]; !ok {
			cfgs.Cfgs[DefaultCfgKey] = &baseDfCfg
		}
		__defaultCfg := cfgs.Cfgs[DefaultCfgKey]
		err := validateCfgFromGfCfg(__defaultCfg, &baseDfCfg, gfcfg, RootKey, DefaultCfgKey)
		if err != nil {
			syslog.Fatalf("readCfgFromYaml failed from %s:\n validate '%s' cfg failed:\n %s", file_, DefaultCfgKey, err.Error())
		}
	}

	curDfCfg := cfgs.Cfgs[DefaultCfgKey]

	for name, curcfg := range cfgs.Cfgs {

		if name == DefaultCfgKey {
			continue
		}

		err := validateCfgFromGfCfg(curcfg, curDfCfg, gfcfg, RootKey, name)
		if err != nil {
			syslog.Fatalf("readCfgFromYaml failed from %s:\n validate '%s' cfg failed:\n %s", file_, name, err.Error())
		}
	}

	dfCfg = curDfCfg

	return &cfgs
}

func validateCfgFromGfCfg(cfg *Cfg, dfcfg *Cfg, gfcfg *gjson.Json, rootKey string, curKey string) error {

	var err error

	dirCfgPattern 		   := rootKey + "." + curKey + ".dir"
	groupCfgPattern 	   := rootKey + "." + curKey + ".group"
	filenameCfgPattern 	   := rootKey + "." + curKey + ".filename"
	maxsizeCfgPattern      := rootKey + "." + curKey + ".max_size"
	maxBackupsCfgPattern   := rootKey + "." + curKey + ".max_backups"
	maxAgeCfgPattern       := rootKey + "." + curKey + ".max_age"
	compressCfgPattern     := rootKey + "." + curKey + ".compress"
	consoleLevelCfgPattern := rootKey + "." + curKey + ".console_level"
	consoleColorCfgPattern := rootKey + "." + curKey + ".console_color"
	fileLevelCfgPattern    := rootKey + "." + curKey + ".file_level"
	fileColorCfgPattern    := rootKey + "." + curKey + ".file_color"
	stackLevelCfgPattern   := rootKey + "." + curKey + ".stack_level"


	if !gfcfg.Contains(dirCfgPattern) {
		cfg.Dir = dfcfg.Dir
	}
	if !gfcfg.Contains(groupCfgPattern) {
		cfg.Group = dfcfg.Group
	}
	if !gfcfg.Contains(filenameCfgPattern) {
		cfg.FileName = dfcfg.FileName
	}
	if !gfcfg.Contains(maxsizeCfgPattern) {
		cfg.MaxSize = dfcfg.MaxSize
	}
	if !gfcfg.Contains(maxBackupsCfgPattern) {
		cfg.MaxBackups = dfcfg.MaxBackups
	}
	if !gfcfg.Contains(maxAgeCfgPattern) {
		cfg.MaxAge = dfcfg.MaxAge
	}
	if !gfcfg.Contains(compressCfgPattern) {
		cfg.Compress = dfcfg.Compress
	}
	if !gfcfg.Contains(consoleLevelCfgPattern) {
		cfg.ConsoleLevel = dfcfg.ConsoleLevel
	}
	if !gfcfg.Contains(consoleColorCfgPattern) {
		cfg.ConsoleColor = dfcfg.ConsoleColor
	}
	if !gfcfg.Contains(fileLevelCfgPattern) {
		cfg.FileLevel = dfcfg.FileLevel
	}
	if !gfcfg.Contains(fileColorCfgPattern) {
		cfg.FileColor = dfcfg.FileColor
	}
	if !gfcfg.Contains(stackLevelCfgPattern) {
		cfg.StackLevel = dfcfg.StackLevel
	}

	if cfg.MaxSize < 1 {
		return fmt.Errorf("invalid max_size(%d), should > 0", cfg.MaxSize)
	}
	if cfg.MaxBackups < 0 {
		return fmt.Errorf("invalid max_backups(%d), should >= 0", cfg.MaxBackups)
	}
	if cfg.MaxAge < 0 {
		return fmt.Errorf("invalid max_age(%d), should >= 0", cfg.MaxAge)
	}
	if cfg.Group != "" {
		hostname, _ := os.Hostname()
		cfg.Group = strings.Replace(cfg.Group, "<HOSTNAME>", hostname, -1)
	}
	if cfg.FileName != "" {
		hostname, _ := os.Hostname()
		cfg.FileName = strings.Replace(cfg.FileName, "<HOSTNAME>", hostname, -1)
	} else {
		cfg.FileName = "elog"
	}

	cfg.consoleLevel, err = parsingLevelStr(cfg.ConsoleLevel)
	if err != nil {
		return fmt.Errorf("parsing console level: %s", err)
	}
	cfg.fileLevel, err    = parsingLevelStr(cfg.FileLevel)
	if err != nil {
		return fmt.Errorf("parsing file level: %s", err)
	}
	cfg.stackLevel, err    = parsingLevelStr(cfg.StackLevel)
	if err != nil {
		return fmt.Errorf("parsing stack level: %s", err)
	}

	return nil
}

func checkAndValidateCfg(cfg *Cfg) error {

	var err error

	if cfg.MaxSize < 1 {
		return fmt.Errorf("invalid max_size(%d), should > 0", cfg.MaxSize)
	}
	if cfg.MaxBackups < 0 {
		return fmt.Errorf("invalid max_backups(%d), should >= 0", cfg.MaxSize)
	}
	if cfg.MaxAge < 0 {
		return fmt.Errorf("invalid max_age(%d), should >= 0", cfg.MaxSize)
	}
	if cfg.Group != "" {
		hostname, _ := os.Hostname()
		cfg.Group = strings.Replace(cfg.Group, "<HOSTNAME>", hostname, -1)
	}
	if cfg.FileName != "" {
		hostname, _ := os.Hostname()
		cfg.FileName = strings.Replace(cfg.FileName, "<HOSTNAME>", hostname, -1)
	}

	cfg.consoleLevel, err = parsingLevelStr(cfg.ConsoleLevel)
	if err != nil {
		return fmt.Errorf("parsing console level: %s", err)
	}
	cfg.fileLevel, err    = parsingLevelStr(cfg.FileLevel)
	if err != nil {
		return fmt.Errorf("parsing file level: %s", err)
	}
	cfg.stackLevel, err    = parsingLevelStr(cfg.StackLevel)
	if err != nil {
		return fmt.Errorf("parsing stack level: %s", err)
	}

	return nil
}

func parsingLevelStr(str string) (zapcore.Level, error) {
	var level zapcore.Level

	switch str {
		case "debug" : level = zapcore.DebugLevel
		case "info"  : level = zapcore.InfoLevel
		case "warn"  : level = zapcore.WarnLevel
		case "error" : level = zapcore.ErrorLevel
		case "fatal" : level = zapcore.FatalLevel
		case "panic" : level = zapcore.PanicLevel
	default:
		return level, fmt.Errorf("unsupported level '%s'", str)
	}

	return level, nil
}

func init() {

}