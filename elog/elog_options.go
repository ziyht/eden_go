package elog

type Option interface {
	apply(op *option)
}

//    <HOSTNAME> -> hostname of current machine
//    <APP_NAME> -> binary file name of current application
//    <LOG_NAME> -> the name of current logger, in __default, it will set to elog
func Filename(filename string) Option{
	return optionFunc(func(op *option) {
		op.filenameSet = true
		op.filename    = filename
	})
}

func NoFile() Option{
	return Filename("")
}

func NoConsole() Option {
	return ConsoleLevel(LEVEL_NONE)
}

func FileLevel(level Level) Option {
	return optionFunc(func(op *option) {
		op.fileLevelSet = true
		op.fileLevel    = level
	})
}

func FileStackLevel(level Level) Option{
		return optionFunc(func(op *option) {
		op.fileStackLevelSet = true
		op.fileStackLevel    = level
	})
}

func ConsoleLevel(level Level) Option{
	return optionFunc(func(op *option) {
		op.consoleLevelSet = true
		op.consoleLevel    = level
	})
}

func ConsoleStackLevel(level Level) Option{
	return optionFunc(func(op *option) {
		op.consoleStackLevelSet = true
		op.consoleStackLevel    = level
	})
}



