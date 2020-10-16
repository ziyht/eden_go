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

func Filename(filename string) Option{
	return optionFunc(func(op *option) {
		op.filenameSet = true
		op.filename    = filename
	})
}