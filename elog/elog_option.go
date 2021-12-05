package elog

type option struct {
	filename     	       string
	tag          	       string
	consoleLevel 	       Level
	consoleStackLevel    Level
	fileLevel    	       Level
	fileStackLevel       Level
	fileColor            bool
	consoleColor         bool

	filenameSet  	       bool
	tagSet       	       bool
	consoleLevelSet      bool
	consoleStackLevelSet bool
	fileLevelSet         bool
	fileStackLevelSet    bool
	fileColorSet         bool
	consoleColorSet      bool
}

type optionFunc func(*option)
func (update optionFunc) apply(op *option) {
	update(op)
}

func newOption(cfg *Cfg, options... Option) *option {

	var op option

	op.applyOptions(options...)
	op.applyCfg(cfg)

	return &op
}

func (opt *option)clone()*option{
	op := *opt
	return &op
}

func (opt *option)applyOptions(options... Option)*option{
	for _, option := range options{
		option.apply(opt)
	}
	return opt
}

func (opt *option)applyCfg(cfg *Cfg)*option{

	if !opt.filenameSet          { opt.filename     = cfg.FileName }

	if !opt.consoleLevelSet      { opt.consoleLevel, _      = getLevelByStr(cfg.ConsoleLevel) }
	if !opt.consoleStackLevelSet { opt.consoleStackLevel, _ = getLevelByStr(cfg.ConsoleStackLevel) }
	if !opt.fileLevelSet         { opt.fileLevel, _         = getLevelByStr(cfg.FileLevel)}
	if !opt.fileStackLevelSet    { opt.fileStackLevel, _    = getLevelByStr(cfg.FileStackLevel) }
	
	if !opt.fileColorSet         { opt.fileColor    = cfg.FileColor }
  if !opt.consoleColorSet      { opt.consoleColor = cfg.FileColor }

	return opt
}



