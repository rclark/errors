package errors

import (
	"fmt"
	"io"
)

// Stack represents a stack trace.
type Stack []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//   - %s	[<filename>:<line> ...]
//   - %v	<package>.<function>\n\t<filepath>:<line>\n\t...
func (st Stack) Format(s fmt.State, verb rune) {
	if st.IsZero() {
		return
	}

	switch verb {
	case 'v':
		for _, f := range st {
			_, _ = io.WriteString(s, "\n")
			f.Format(s, verb)
		}
	case 's':
		_, _ = io.WriteString(s, "[")
		for i, f := range st {
			if i > 0 {
				_, _ = io.WriteString(s, " ")
			}
			f.Format(s, verb)
		}
		_, _ = io.WriteString(s, "]")
	}
}

// IsZero reports whether the stack trace is empty.
func (st Stack) IsZero() bool {
	return len(st) == 0
}

// StackTracer is implemented by [Error]. It can be used in external contexts
// to check whether an error has a stack trace that this package can expose.
type StackTracer interface {
	StackTrace() Stack
}
