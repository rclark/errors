[![Go](https://github.com/rclark/errors/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/rclark/errors/actions/workflows/go.yml)

# errors

Another package for errors with stack traces included.

```go
import "github.com/rclark/errors"
```

## Index

- [func As\(err error, target interface\{\}\) bool](<#As>)
- [func Is\(err, target error\) bool](<#Is>)
- [func Join\(errs ...error\) error](<#Join>)
- [func Unwrap\(err error\) error](<#Unwrap>)
- [func WithStack\(err error, opts ...WithStackOption\) error](<#WithStack>)
- [type Error](<#Error>)
  - [func New\(message string\) Error](<#New>)
  - [func \(e Error\) Error\(\) string](<#Error.Error>)
  - [func \(e Error\) StackTrace\(\) Stack](<#Error.StackTrace>)
  - [func \(e Error\) Unwrap\(\) error](<#Error.Unwrap>)
- [type Frame](<#Frame>)
  - [func \(f Frame\) Format\(s fmt.State, verb rune\)](<#Frame.Format>)
  - [func \(f Frame\) String\(\) string](<#Frame.String>)
- [type Stack](<#Stack>)
  - [func StackTrace\(err error\) Stack](<#StackTrace>)
  - [func \(st Stack\) Format\(s fmt.State, verb rune\)](<#Stack.Format>)
  - [func \(st Stack\) IsZero\(\) bool](<#Stack.IsZero>)
  - [func \(st Stack\) String\(\) string](<#Stack.String>)
- [type StackTracer](<#StackTracer>)
- [type WithStackOption](<#WithStackOption>)
  - [func NoOverwrite\(\) WithStackOption](<#NoOverwrite>)


<a name="As"></a>
## func As

```go
func As(err error, target interface{}) bool
```

As finds the first error in err's tree that matches target, and if one is found, sets target to that error value and returns true. Otherwise, it returns false.

The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap\(\) error or Unwrap\(\) \[\]error method. When err wraps multiple errors, As examines err followed by a depth\-first traversal of its children.

An error matches target if the error's concrete value is assignable to the value pointed to by target, or if the error has a method As\(interface\{\}\) bool such that As\(target\) returns true. In the latter case, the As method is responsible for setting target.

An error type might provide an As method so it can be treated as if it were a different error type.

As panics if target is not a non\-nil pointer to either a type that implements error, or to any interface type.

<a name="Is"></a>
## func Is

```go
func Is(err, target error) bool
```

Is reports whether any error in err's tree matches target.

The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap\(\) error or Unwrap\(\) \[\]error method. When err wraps multiple errors, Is examines err followed by a depth\-first traversal of its children.

An error is considered to match a target if it is equal to that target or if it implements a method Is\(error\) bool such that Is\(target\) returns true.

An error type might provide an Is method so it can be treated as equivalent to an existing error. For example, if MyError defines

```
func (m MyError) Is(target error) bool { return target == fs.ErrExist }
```

then Is\(MyError\{\}, fs.ErrExist\) returns true. See syscall.Errno.Is for an example in the standard library. An Is method should only shallowly compare err and the target and not call [Unwrap](<#Unwrap>) on either.

<a name="Join"></a>
## func Join

```go
func Join(errs ...error) error
```

Join returns an error that wraps the given errors. Any nil error values are discarded. Join returns nil if every value in errs is nil. The error formats as the concatenation of the strings obtained by calling the Error method of each element of errs, with a newline between each string.

A non\-nil error returned by Join implements the Unwrap\(\) \[\]error method.

<a name="Unwrap"></a>
## func Unwrap

```go
func Unwrap(err error) error
```

Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.

Unwrap only calls a method of the form "Unwrap\(\) error". In particular Unwrap does not unwrap errors returned by [Join](<#Join>).

<a name="WithStack"></a>
## func WithStack

```go
func WithStack(err error, opts ...WithStackOption) error
```

WithStack adds a [Stack](<#Stack>) to the provided error at the point where the function was called. If the error already has a [Stack](<#Stack>), it will be overridden unless the [NoOverwrite](<#NoOverwrite>) option is provided.

<a name="Error"></a>
## type Error

Error implements the error interface and provides a stack trace.

```go
type Error struct {
    // contains filtered or unexported fields
}
```

<a name="New"></a>
### func New

```go
func New(message string) Error
```

New returns an error with the supplied message and a stack trace to the point where the function was called.

<a name="Error.Error"></a>
### func \(Error\) Error

```go
func (e Error) Error() string
```

Error returns the error message.

<a name="Error.StackTrace"></a>
### func \(Error\) StackTrace

```go
func (e Error) StackTrace() Stack
```

StackTrace returns the [Stack](<#Stack>).

<a name="Error.Unwrap"></a>
### func \(Error\) Unwrap

```go
func (e Error) Unwrap() error
```

Unwrap returns the wrapped error, if any.

<a name="Frame"></a>
## type Frame

Frame represents a program counter inside a stack trace.

```go
type Frame struct {
    File     string
    Line     int
    Function string
    // contains filtered or unexported fields
}
```

<a name="Frame.Format"></a>
### func \(Frame\) Format

```go
func (f Frame) Format(s fmt.State, verb rune)
```

Format formats the frame according to the fmt.Formatter interface.

```
%s    source file
%d    source line
%n    function name
%v    equivalent to %s:%d
```

Format accepts flags that alter the printing of some verbs, as follows:

```
%+s   function name and path of source file relative to the compile time
      GOPATH separated by \n\t (<func>\n\t<path>)
%+v   equivalent to %+s:%d
```

<a name="Frame.String"></a>
### func \(Frame\) String

```go
func (f Frame) String() string
```



<a name="Stack"></a>
## type Stack

Stack represents a stack trace.

```go
type Stack []Frame
```

<a name="StackTrace"></a>
### func StackTrace

```go
func StackTrace(err error) Stack
```

StackTrace returns a [Stack](<#Stack>), if err has one. If none was found, the returned [Stack](<#Stack>) will be empty. This can be checked with the [Stack.IsZero](<#Stack.IsZero>) method.

<a name="Stack.Format"></a>
### func \(Stack\) Format

```go
func (st Stack) Format(s fmt.State, verb rune)
```

Format formats the stack of Frames according to the fmt.Formatter interface.

```
%s	lists source files for each Frame in the stack
%v	lists the source file and line number for each Frame in the stack
```

Format accepts flags that alter the printing of some verbs, as follows:

```
%+v   Prints filename, function, and line number for each Frame in the stack.
```

<a name="Stack.IsZero"></a>
### func \(Stack\) IsZero

```go
func (st Stack) IsZero() bool
```

IsZero reports whether the stack trace is empty.

<a name="Stack.String"></a>
### func \(Stack\) String

```go
func (st Stack) String() string
```



<a name="StackTracer"></a>
## type StackTracer

StackTracer is implemented by [Error](<#Error>). It can be used in external contexts to check whether an error has a stack trace that this package can expose.

```go
type StackTracer interface {
    StackTrace() Stack
}
```

<a name="WithStackOption"></a>
## type WithStackOption

WithStackOption is an option for the WithStack function.

```go
type WithStackOption func(*options)
```

<a name="NoOverwrite"></a>
### func NoOverwrite

```go
func NoOverwrite() WithStackOption
```

NoOverwrite is an option that prevents the stack trace from being overridden if the error already has one.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
