package elog

import "go.uber.org/zap/zapcore"

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
	PanicLevel = zapcore.PanicLevel
)

type option struct {
	filename     	string
	filenameSet  	bool
	tag          	string
	tagSet       	bool
	consoleLevel 	Level
	consoleLevelSet bool
	fileLevel    	Level
	fileLevelSet    bool
	noConsole       bool
	noConsoleSet    bool
}

type optionFunc func(*option)
func (update optionFunc) apply(op *option) {
	update(op)
}

func newOption(cfg *Cfg, options... Option) *option {

	var op option

	op._updateFromOptions(options...)
	op._updateFromCfg(cfg)

	return &op
}

func (opt *option)needConsole() bool {

	if opt.noConsoleSet{
		if opt.noConsole {
			return false
		}
	}

	return true
}

func (opt *option)_updateFromOptions(options... Option){
	for _, option := range options{
		option.apply(opt)
	}
}

func (opt *option)_updateFromCfg(cfg *Cfg){

	if !opt.consoleLevelSet {
		opt.consoleLevel = cfg.consoleLevel
	}

	if !opt.fileLevelSet {
		opt.fileLevel = cfg.fileLevel
	}

	if !opt.filenameSet {
		opt.filename = cfg.FileName
	}
}