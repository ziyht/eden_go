// Package tracerr makes error output more informative.
// It adds stack trace to error and can display error with source fragments.
//
// Check example of output here https://github.com/ztrue/tracerr
package eerr

import (
	"fmt"

	"github.com/cespare/xxhash/v2"
)

// DefaultCap is a default cap for frames array.
// It can be changed to number of expected frames
// for purpose of performance optimisation.
var DefaultCap = 20

// Error is an error with stack trace.
type Error interface {
	Error()      string
	StackTrace() []Frame
	StaskCause() Frame
	Unpack()     error
	Id()         uint64      // a hash value caculated by file, line number and error string 
}

// Errorf creates new error with stacktrace and formatted message.
// Formatting works the same way as in fmt.Errorf.
func Errorf(message string, args ...interface{}) Error {
	return &errorData{fmt.Errorf(message, args...), Trace(2)}
}

// Newf is the same as Errorf
func Newf(message string, args ...interface{}) Error {
	return &errorData{fmt.Errorf(message, args...), Trace(2)}
}

// NewSkipf is the same as Errorf
// The parameter `skip` specifies the stack callers skipped amount.
func NewSkipf(skip int, message string, args ...interface{}) Error {
	return &errorData{fmt.Errorf(message, args...), Trace(skip+2)}
}

// New creates new error with stacktrace.
func New(message string) Error {
	return &errorData{fmt.Errorf(message), Trace(2)}
}

// New creates new error with stacktrace.
// The parameter `skip` specifies the stack callers skipped amount.
func NewSkip(skip int, message string) Error {
	return &errorData{fmt.Errorf(message), Trace(skip+2)}
}

// Pack adds stacktrace to existing error. 
// Note: it takes no effect for eerr.Error
func Pack(err error) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	return &errorData{err, Trace(2)}
}

// PackSkip adds stacktrace to existing error.
// The parameter `skip` specifies the stack callers skipped amount.
// Note: it takes no effect for eerr.Error
func PackSkip(skip int, err error) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if ok {
		return e
	}
	return &errorData{err, Trace(skip+2)}
}

// Unpack returns the original error.
func Unpack(err error) error {
	if err == nil {
		return nil
	}
	e, ok := err.(Error)
	if !ok {
		return err
	}
	return e.Unpack()
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

// a hash value caculated by file, line number and error string
// if not a err, only using the error string
func Id(err error) uint64 {
	if err == nil {
		return 0
	}

	e, ok := err.(Error)
	if ok {
		return e.Id()
	}

	return xxhash.Sum64String(err.Error())
}

func StackCause(err error)(f Frame, ok bool){
	e, ok := err.(*errorData)
	if !ok {
		return
	}

	if len(e.stack)== 0 {
		return
	}

	return e.stack.Frame(0), true
}
