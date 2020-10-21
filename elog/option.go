package elog

type option struct {
	filename    string
	filenameSet bool
	tag         string
	tagSet      bool
}

type optionFunc func(*option)
func (opf optionFunc) apply(op *option) {
	opf(op)
}

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