// Package tracerr makes error output more informative.
// It adds stack trace to error and can display error with source fragments.
//
// Check example of output here https://github.com/ztrue/tracerr
package eerr

import (
	"fmt"
	"runtime"
)

// DefaultCap is a default cap for frames array.
// It can be changed to number of expected frames
// for purpose of performance optimisation.
var DefaultCap = 20

// Error is an error with stack trace.
type Error interface {
	Error()      string
	StackTrace() []Frame
	Unwrap()     error
	Id()         uint64      // a hash value caculated by file, line number and error string 
}

// CustomError creates an error with provided frames.
func CustomError(err error, frames []Frame) Error {
	return &errorData{
		err:    err,
		frames: frames,
	}
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) Error {
	return trace(fmt.Errorf(message, args...), 2)
}

// Newf is the same as Errorf
func Newf(message string, args ...interface{}) Error {
	return trace(fmt.Errorf(message, args...), 2)
}

// New creates new error with stacktrace.
func New(message string) Error {
	return trace(fmt.Errorf(message), 2)
}

// Wrap adds stacktrace to existing error.
func Wrap(err error) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	return trace(err, 2)
}

// Unwrap returns the original error.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unwrap()
}

// StackTrace returns stack trace of an error.
// It will be empty if err is not of type Error.
func StackTrace(err error) []Frame {
	e, ok := err.(Error)
	if !ok {
		return nil
	}
	return e.StackTrace()
}

func HasStack(err error) bool {
	_, ok := err.(Error)
	return ok
}

func StackCause(err error)(f Frame, ok bool){
	e, ok := err.(*errorData)
	if !ok {
		return
	}

	if len(e.frames)== 0 {
		return
	}

	return e.frames[0], true
}

func trace(err error, skip int) Error {
	pc := make([]uintptr, DefaultCap)
	cnt := runtime.Callers(skip+1, pc)
	if cnt == 0 {
		return &errorData{
			err:    err,
		}
	}

	frames := make([]Frame, 0, cnt)
	fs := runtime.CallersFrames(pc)
	for ;; {
		f, ok := fs.Next()
		frame := Frame{
			Func: f.Func.Name(),
			Line: f.Line,
			Path: f.File,
		}
		frames = append(frames, frame)
		if !ok {
			break
		}
	}

	return &errorData{
		err:    err,
		frames: frames,
	}
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
