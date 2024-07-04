package errors

import (
	"runtime"
)

// Error implements the error interface and provides a stack trace.
type Error struct {
	err     error
	message string
	stack   []Frame
}

func newError(message string, skip int) Error {
	var full [32]uintptr
	n := runtime.Callers(skip, full[:])
	frames := make([]Frame, n)
	for i, pc := range full[:n] {
		frames[i] = newFrame(pc)
	}

	return Error{
		message: message,
		stack:   frames,
	}
}

func wrapError(err error, skip int) Error {
	e := newError(err.Error(), skip+1)
	e.err = err
	return e
}

// Error returns the error message.
func (e Error) Error() string {
	return e.message
}

// StackTrace returns the [Stack].
func (e Error) StackTrace() Stack {
	return e.stack
}

// Unwrap returns the wrapped error, if any.
func (e Error) Unwrap() error {
	return e.err
}
