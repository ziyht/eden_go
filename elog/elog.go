package elog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*  template
logger:
  dir          : var/log                 # default var/log
  group        : <HOSTNAME>              # default <HOSTNAME>, if set, real dir will be $Dir/$Group,  <HOSTNAME> represent hostname of current machine
  max_size     : 100                     # default 100, unit MB
  max_backups  : 7                       # default 7
  max_age      : 7                       # default 7
  compress     : false                   # default false
  level_console: info                    # default info      [debug, info, warn, error, fatal]
  level_file   : debug                   # default info      [debug, info, warn, error, fatal]
*/

// ElogCfg ...
type Cfg struct{
	Dir          	string `yaml:"dir"`
	Group           string `yaml:"group"`
	MaxSize     	int    `yaml:"max_size"`
	MaxBackups  	int    `yaml:"max_backups"`
	MaxAge      	int    `yaml:"max_age"`
	Compress    	bool   `yaml:"compress"`
	LevelConsole    string `yaml:"level_console"`
	LevelFile       string `yaml:"level_file"`
}

func DfCfg() Cfg {
	return Cfg{
		Dir: "var/log",
		Group: "<HOSTNAME>",
		MaxSize: 7,
		MaxBackups: 7,
		MaxAge: 7,
		Compress: false,
		LevelConsole: "info",
		LevelFile: "info",
	}
}

// Elogger ...
type Logger struct {
	cfg          *Cfg
	levelConsole zapcore.Level
	levelFile    zapcore.Level
}

type Log = zap.SugaredLogger

var dfLogger *Logger

// NewLogger ...
func NewLogger(cfg* Cfg) *Logger {

	if cfg == nil {
		return dfLogger
	}

	out := new(Logger)

	out.cfg = cfg
	out.validateCfg()

	return out
}

// DfLogger ...
func DfLogger() *Logger {
	return dfLogger
}

func InitDfLogger(cfg *Cfg) {

	if cfg == nil {
		if dfLogger == nil {
			dfCfg := DfCfg()
			dfLogger = NewLogger(&dfCfg)
		}
	}

	dfLogger = NewLogger(cfg)
}

// GetLog ...
func (l *Logger)GetLog(filename string, tag string) *Log {

	var logger *zap.Logger
	var path string

	if filepath.IsAbs(filename){
		path = filepath.Join(filename)
	} else if strings.HasSuffix(filename, ".log"){
		path = filepath.Join(l.cfg.Dir, filename)
	} else {
		path = filepath.Join(l.cfg.Dir, filename + ".log")
	}

	coreConsole := zapcore.NewCore(getEncoder(), zapcore.AddSync(os.Stdout), l.levelConsole)
    coreFile    := zapcore.NewCore(getEncoder(), l.getLogWriter(path), l.levelFile)

	logger = zap.New(zapcore.NewTee(coreConsole, coreFile))

	return logger.Sugar()
}

func (l *Logger) validateCfg() {

	cfg := Cfg{}
	cfg = *l.cfg

	if cfg.Dir == ""{
		cfg.Dir = filepath.Join("var/log")
	}
	if cfg.MaxSize < 1 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups < 1 {
		cfg.MaxBackups = 7
	}
	if cfg.MaxAge < 1 {
		cfg.MaxAge = 7
	}
	if cfg.Group != "" {
		hostname, _ := os.Hostname()
		cfg.Group = strings.Replace(cfg.Group, "<HOSTNAME>", hostname, -1)
		cfg.Dir = filepath.Join(cfg.Dir, cfg.Group)
	}

	levelConsole := zapcore.InfoLevel
	levelFile    := zapcore.InfoLevel

	switch cfg.LevelConsole {
		case "debug" : levelConsole = zapcore.DebugLevel
		case "info"  : levelConsole = zapcore.InfoLevel
		case "warn"  : levelConsole = zapcore.WarnLevel
		case "error" : levelConsole = zapcore.ErrorLevel
	}

	switch cfg.LevelFile {
		case "debug" : levelFile = zapcore.DebugLevel
		case "info"  : levelFile = zapcore.InfoLevel
		case "warn"  : levelFile = zapcore.WarnLevel
		case "error" : levelFile = zapcore.ErrorLevel
	}

	l.cfg = &cfg
	l.levelConsole = levelConsole
	l.levelFile    = levelFile
}

func myTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02_15:04:05.000"))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime  = myTimeEncoder // zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (l *Logger)getLogWriter(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    l.cfg.MaxSize,
		MaxBackups: l.cfg.MaxBackups,
		MaxAge:     l.cfg.MaxAge,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func init() {
	InitDfLogger(nil)
}

