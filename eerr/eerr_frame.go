package eerr

import "fmt"

// Frame is a single step in stack trace.
type Frame struct {
	line  int
	func_ string   
	file  string  
}

func (f *Frame)Func() string { return f.func_ }
func (f *Frame)File() string { return f.file }
func (f *Frame)Line() int    { return f.line }

// String formats Frame to string.
func (f *Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.file, f.line, f.func_)
}

