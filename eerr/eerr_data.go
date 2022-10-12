package eerr

import "github.com/cespare/xxhash/v2"

type errorData struct {
	// err contains original error.
	err error
	// frames contains stack trace of an error.
	stack stack
}

// Error returns error message.
func (e *errorData) Error() string {
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *errorData) StackTrace() []Frame {
	return e.stack.Frames()
}

// Unwrap returns the original error.
func (e *errorData) StaskCause() Frame {
	return e.stack.Frame(0)
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

// a hash value caculated by file, line number and error string
// for this is a not frequently using function, now re caculated each time by call
func (e *errorData) Id() uint64 {
	h := xxhash.New()

	if len(e.stack) > 0 {
		f := e.stack.Frame(0)
		h.WriteString(f.file)
		line := int32(f.line)
		h.Write([]byte{0, byte(line), byte(line>>8), byte(line>>16), byte(line>>24), 0})  // 0 is a gap between file and number and left content
	}

	h.WriteString(e.err.Error())
	return h.Sum64()
}
