package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

// Frame represents a program counter inside a stack trace.
type Frame struct {
	pc uintptr
	fn *runtime.Func

	File     string
	Line     int
	line     string
	Function string
}

func newFrame(c uintptr) Frame {
	if c == 0 {
		return Frame{}
	}

	f := Frame{pc: c}

	f.fn = runtime.FuncForPC(f.pc - 1)
	if f.fn == nil {
		return f
	}

	// file path & line number
	f.File, f.Line = f.fn.FileLine(f.pc - 1)
	f.line = strconv.Itoa(f.Line)

	// function name with package prefix
	f.Function = f.fn.Name()

	return f
}

// Format formats the frame according to the fmt.Formatter interface.
//
//   - %s <filename>:<line>
//   - %v <package>.<function>\n\t<filepath>:<line>
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = io.WriteString(s, path.Base(f.File)+":"+f.line)
	case 'v':
		_, _ = io.WriteString(s, f.Function)
		_, _ = io.WriteString(s, "\n\t")
		_, _ = io.WriteString(s, f.File+":"+f.line)
	}
}

func (f Frame) String() string {
	return fmt.Sprintf("%+v", f)
}
