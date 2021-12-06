package elog

type option struct {
	dir                  string
	group                string
	filename     	       string
	tag          	       string
	consoleLevel 	       Level
	consoleStackLevel    Level
	fileLevel    	       Level
	fileStackLevel       Level
	fileColor            bool
	consoleColor         bool

	dirSet               bool
	groupSet             bool
	filenameSet  	       bool
	tagSet       	       bool
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

func (o *option)NoFile   () *option {return o.FileName("")}
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
func (o *option)FileName         (v string) *option {o.filename          = v; o.filenameSet          = true; return o }
func (o *option)FileLevel        (l Level)  *option {o.fileLevel         = l; o.fileLevelSet         = true; return o }
func (o *option)ConsoleLevel     (l Level)  *option {o.consoleLevel      = l; o.consoleLevelSet      = true; return o}
func (o *option)FileStackLevel   (l Level)  *option {o.fileStackLevel    = l; o.fileStackLevelSet    = true; return o}
func (o *option)ConsoleStackLevel(l Level)  *option {o.consoleStackLevel = l; o.consoleStackLevelSet = true; return o}
func (o *option)ConsoleColor     (b bool )  *option {o.consoleColor      = b; o.consoleColorSet      = true; return o}
func (o *option)FileColor        (b bool )  *option {o.fileColor         = b; o.fileColorSet         = true; return o}

type optionFunc func(*option)
func (update optionFunc) apply(op *option) {
	update(op)
}

func (opt *option)clone()*option{
	op := *opt
	return &op
}

func (o *option)applyOptions(ns... *option)*option{
	for _, n := range ns{
		o.applyOption(n)
	}
	return o
}

func (o *option)applyOption(n *option) *option {

	if n.dirSet               { o.Dir     (n.dir)      }
	if n.groupSet             { o.Group   (n.group)    }
	if n.filenameSet          { o.FileName(n.filename) }

	if n.consoleLevelSet      { o.ConsoleLevel     (n.consoleLevel     ) }
	if n.consoleStackLevelSet { o.ConsoleStackLevel(n.consoleStackLevel) }
	if n.fileLevelSet         { o.FileLevel        (n.fileLevel        ) }
	if n.fileStackLevelSet    { o.FileStackLevel   (n.fileStackLevel   ) }
	
	if n.fileColorSet         { o.FileColor(n.fileColor) }
  if n.consoleColorSet      { o.ConsoleColor(n.consoleColor) }

	return o

}

func (opt *option)applyCfg(cfg *Cfg)*option{
	if !opt.dirSet               { opt.dir          = cfg.Dir      }
	if !opt.groupSet             { opt.group        = cfg.Group    }
	if !opt.filenameSet          { opt.filename     = cfg.FileName }

	if !opt.consoleLevelSet      { opt.consoleLevel, _      = getLevelByStr(cfg.ConsoleLevel) }
	if !opt.consoleStackLevelSet { opt.consoleStackLevel, _ = getLevelByStr(cfg.ConsoleStackLevel) }
	if !opt.fileLevelSet         { opt.fileLevel, _         = getLevelByStr(cfg.FileLevel)}
	if !opt.fileStackLevelSet    { opt.fileStackLevel, _    = getLevelByStr(cfg.FileStackLevel) }
	
	if !opt.fileColorSet         { opt.fileColor    = cfg.FileColor    }
  if !opt.consoleColorSet      { opt.consoleColor = cfg.ConsoleColor }

	return opt
}



