package eerr

import "github.com/cespare/xxhash/v2"

type errorData struct {
	// err contains original error.
	err error
	// frames contains stack trace of an error.
	frames []Frame
}

// Error returns error message.
func (e *errorData) Error() string {
	return e.err.Error()
}

// StackTrace returns stack trace of an error.
func (e *errorData) StackTrace() []Frame {
	return e.frames
}

// Unwrap returns the original error.
func (e *errorData) StaskCause() (out Frame) {
	if len(e.frames) > 0 {
		return e.frames[0]
	}
	return
}

// Unwrap returns the original error.
func (e *errorData) Unwrap() error {
	return e.err
}

// a hash value caculated by file, line number and error string
// for this is a not frequently using function, now re caculated each time by call
func (e *errorData) Id() uint64 {
	h := xxhash.New()

	if len(e.frames) > 0 {
		h.WriteString(e.frames[0].Path)
		line := int32(e.frames[0].Line)
		h.Write([]byte{0, byte(line), byte(line>>8), byte(line>>16), byte(line>>24), 0})  // 0 is a gap between file and number and left content
	}

	h.WriteString(e.err.Error())
	return h.Sum64()
}
