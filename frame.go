package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// Frame represents a program counter inside a stack trace.
type Frame struct {
	pc uintptr
	fn *runtime.Func

	File     string
	Line     int
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

	f.File, f.Line = f.fn.FileLine(f.pc - 1)

	name := f.fn.Name()
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	f.Function = name[i+1:]

	return f
}

/*
Format function is copied from code available at https://github.com/pkg/errors

Copyright (c) 2015, Dave Cheney <dave@cheney.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Format formats the frame according to the fmt.Formatter interface.
//
//	%s    source file
//	%d    source line
//	%n    function name
//	%v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+s   function name and path of source file relative to the compile time
//	      GOPATH separated by \n\t (<func>\n\t<path>)
//	%+v   equivalent to %+s:%d
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, f.Function)
			_, _ = io.WriteString(s, "\n\t")
			_, _ = io.WriteString(s, f.File)
		default:
			_, _ = io.WriteString(s, path.Base(f.File))
		}
	case 'd':
		_, _ = io.WriteString(s, strconv.Itoa(f.Line))
	case 'n':
		_, _ = io.WriteString(s, f.Function)
	case 'v':
		f.Format(s, 's')
		_, _ = io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

func (f Frame) String() string {
	return fmt.Sprintf("%+v", f)
}
