package eerr

import "runtime"

type stack []uintptr

func Call(skip int) (f Frame) {
	pc, file, l, ok := runtime.Caller(skip + 1)
	if ok {
		f.file  = file
		f.line  = l
		f.func_ = runtime.FuncForPC(pc).Name() 
	}
	return
}

func Trace(skip int) stack {
	pc := make([]uintptr, DefaultCap)
	cnt := runtime.Callers(skip+1, pc)
	if cnt == 0 {
		return nil
	}
	
	return pc
}

func (s stack)Frames() []Frame {
	frames := make([]Frame, 0, len(s))
	fs := runtime.CallersFrames(s)
	for ;; {
		f, ok := fs.Next()
		frame := Frame{
			func_: f.Func.Name(),
			line : f.Line,
			file : f.File,
		}
		frames = append(frames, frame)
		if !ok {
			break
		}
	}
	return frames
}

func (s stack)Frame(deep int) (f Frame) {
	if deep < len(s) {
		fp := runtime.FuncForPC(s[deep])
		f.func_ = fp.Name()
		f.file, f.line = fp.FileLine(s[deep])
	}

	return
}

// func trace2(err error, skip int) Error {
// 	frames := make([]Frame, 0, DefaultCap)
// 	for {
// 		pc, path, line, ok := runtime.Caller(skip)
// 		if !ok {
// 			break
// 		}
// 		fn := runtime.FuncForPC(pc)
// 		frame := Frame{
// 			Func: fn.Name(),
// 			Line: line,
// 			Path: path,
// 		}
// 		frames = append(frames, frame)
// 		skip++
// 	}
// 	return &errorData{
// 		err:    err,
// 		frames: frames,
// 	}
// }