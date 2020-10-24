package elog

type Option interface {
	apply(op *option)
}

func Tag(tag string) Option{
	return optionFunc(func(op *option) {
		op.tagSet = true
		op.tag    = tag
	})
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

func FileLevel(level Level) Option {
	return optionFunc(func(op *option) {
		op.fileLevelSet = true
		op.fileLevel    = level
	})
}

func ConsoleLevel(level Level) Option{
	return optionFunc(func(op *option) {
		op.consoleLevelSet = true
		op.consoleLevel    = level
	})
}

func NoConsole() Option {
	return optionFunc(func(op *option) {
		op.noConsoleSet = true
		op.noConsole    = true
	})
}