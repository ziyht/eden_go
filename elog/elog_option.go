package elog

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type option struct {
	tags               []string
	dir                  string
	group                string
	filename     	       string
	consoleLevel 	       zapcore.LevelEnabler
	consoleStackLevel    zapcore.LevelEnabler
	fileLevel    	       zapcore.LevelEnabler
	fileStackLevel       zapcore.LevelEnabler
	fileColor            colorSwitch
	consoleColor         colorSwitch

	reTagSet             bool
	dirSet               bool
	groupSet             bool
	filenameSet  	       bool
	consoleLevelSet      bool
	consoleStackLevelSet bool
	fileLevelSet         bool
	fileStackLevelSet    bool
	fileColorSet         bool
	consoleColorSet      bool
}

func newOpt() *option {
	return &option{}
}

// add tags for new log, it will append to the tags set in logger, empty tag will be skipped
func (o *option)Tags(tags ...string) *option   { for _, tag := range tags { if tag != "" { o.tags = append(o.tags, tag)} } ; return o}
// add tags for new log, tags set in logger will take no effect, empty tag will be skipped
func (o *option)ReTags(tags ...string) *option { o.reTagSet = true; return o.Tags(tags...)}

func (o *option)NoFile   () *option {return o.Filename("")}
func (o *option)NoConsole() *option {return o.ConsoleLevel(LEVEL_NONE)}

//    <HOSTNAME> -> hostname of current machine;
//    <APP_NAME> -> binary file name of current application;
//    <LOG_NAME> -> the name of current logger, in __default, it will set to elog;
func (o *option)Dir              (v string) *option {o.dir               = v; o.dirSet               = true; return o }
//    <HOSTNAME> -> hostname of current machine;
//    <APP_NAME> -> binary file name of current application;
//    <LOG_NAME> -> the name of current logger, in __default, it will set to elog;
func (o *option)Group            (v string) *option {o.group             = v; o.groupSet             = true; return o }
//    <HOSTNAME> -> hostname of current machine;
//    <APP_NAME> -> binary file name of current application;
//    <LOG_NAME> -> the name of current logger, in __default, it will set to elog;
func (o *option)Filename         (v string) *option {o.filename          = v; o.filenameSet          = true; return o }
func (o *option)FileLevel        (l zapcore.LevelEnabler)  *option {o.fileLevel         = l; o.fileLevelSet         = true; return o }
func (o *option)ConsoleLevel     (l zapcore.LevelEnabler)  *option {o.consoleLevel      = l; o.consoleLevelSet      = true; return o}
func (o *option)FileStackLevel   (l zapcore.LevelEnabler)  *option {o.fileStackLevel    = l; o.fileStackLevelSet    = true; return o}
func (o *option)ConsoleStackLevel(l zapcore.LevelEnabler)  *option {o.consoleStackLevel = l; o.consoleStackLevelSet = true; return o}
func (o *option)ConsoleColor     (b colorSwitch )  *option {o.consoleColor      = b; o.consoleColorSet      = true; return o}
func (o *option)FileColor        (b colorSwitch )  *option {o.fileColor         = b; o.fileColorSet         = true; return o}

type optionFunc func(*option)
func (update optionFunc) apply(op *option) {
	update(op)
}

func (opt *option)clone()*option{
	op := *opt
	op.tags = opt.tags[:]
	return &op
}

func (o *option)applyOptions(ns... *option)*option{
	for _, n := range ns{
		o.applyOption(n)
	}
	return o
}

func (o *option)applyOption(n *option) *option {
	if !n.reTagSet            { o.tags = append(o.tags, n.tags...)
	} else                    { o.tags = n.tags[:]}

	if n.dirSet               { o.Dir     (n.dir)      }
	if n.groupSet             { o.Group   (n.group)    }
	if n.filenameSet          { o.Filename(n.filename) }

	if n.consoleLevelSet      { o.ConsoleLevel     (n.consoleLevel     ) }
	if n.consoleStackLevelSet { o.ConsoleStackLevel(n.consoleStackLevel) }
	if n.fileLevelSet         { o.FileLevel        (n.fileLevel        ) }
	if n.fileStackLevelSet    { o.FileStackLevel   (n.fileStackLevel   ) }
	
	if n.fileColorSet         { o.FileColor(n.fileColor) }
  if n.consoleColorSet      { o.ConsoleColor(n.consoleColor) }

	return o
}

func (opt *option)applyCfg(cfg *Cfg)*option{
	if !opt.reTagSet             { if cfg.Tag != "" {opt.tags = append(opt.tags, cfg.Tag)} }

	if !opt.dirSet               { opt.dir          = cfg.Dir      }
	if !opt.groupSet             { opt.group        = cfg.Group    }
	if !opt.filenameSet          { opt.filename     = cfg.FileName }

	if !opt.consoleLevelSet      { opt.consoleLevel      = cfg.ConsoleLevel }
	if !opt.consoleStackLevelSet { opt.consoleStackLevel = cfg.ConsoleStackLevel }
	if !opt.fileLevelSet         { opt.fileLevel         = cfg.FileLevel}
	if !opt.fileStackLevelSet    { opt.fileStackLevel    = cfg.FileStackLevel }
	
	if !opt.fileColorSet         { opt.fileColor    = cfg.FileColor    }
  if !opt.consoleColorSet      { opt.consoleColor = cfg.ConsoleColor }

	return opt
}

func (opt *option)applyToLogCfg(cfg *LogCfg)*LogCfg{

	cfg = cfg.Clone()

	if opt.reTagSet               { cfg.Tag  = strings.Join(opt.tags, ".")
	} else if len(opt.tags) > 0   { cfg.Tag += strings.Join(opt.tags, ".")}

	if cfg.file {
		if opt.dirSet               { cfg.Dir        = opt.dir      }
		if opt.groupSet             { cfg.Group      = opt.group    }
		if opt.filenameSet          { cfg.FileName   = opt.filename }
		if opt.fileLevelSet         { cfg.Level      = opt.fileLevel }
		if opt.fileStackLevelSet    { cfg.StackLevel = opt.fileStackLevel }
		if opt.fileColorSet         { cfg.Color      = opt.fileColor }
	} else {
		if opt.filenameSet          { cfg.FileName   = opt.filename }
		if opt.consoleLevelSet      { cfg.Level      = opt.consoleLevel }
		if opt.consoleStackLevelSet { cfg.StackLevel = opt.consoleStackLevel }
		if opt.consoleColorSet      { cfg.Color      = opt.consoleColor }
	}

	return cfg
}

