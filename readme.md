[![Go](https://github.com/rclark/errors/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/rclark/errors/actions/workflows/go.yml)

# errors

Another package for errors with stack traces included.

```go
import "github.com/rclark/errors"
```

## Index

- [func As\(err error, target interface\{\}\) bool](<#As>)
- [func Errorf\(format string, args ...any\) error](<#Errorf>)
- [func Is\(err, target error\) bool](<#Is>)
- [func Join\(errs ...error\) error](<#Join>)
- [func New\(message string\) error](<#New>)
- [func Unwrap\(err error\) error](<#Unwrap>)
- [func UnwrapAny\(err error\) \[\]error](<#UnwrapAny>)
- [func WithStack\(err error, opts ...StackOption\) error](<#WithStack>)
- [type Error](<#Error>)
  - [func \(e Error\) Error\(\) string](<#Error.Error>)
  - [func \(e Error\) Format\(s fmt.State, verb rune\)](<#Error.Format>)
  - [func \(e Error\) StackTrace\(\) Stack](<#Error.StackTrace>)
  - [func \(e Error\) Unwrap\(\) error](<#Error.Unwrap>)
- [type Frame](<#Frame>)
  - [func \(f Frame\) Format\(s fmt.State, verb rune\)](<#Frame.Format>)
  - [func \(f Frame\) String\(\) string](<#Frame.String>)
- [type Stack](<#Stack>)
  - [func StackTrace\(err error\) \(Stack, bool\)](<#StackTrace>)
  - [func \(st Stack\) Format\(s fmt.State, verb rune\)](<#Stack.Format>)
  - [func \(st Stack\) IsZero\(\) bool](<#Stack.IsZero>)
- [type StackOption](<#StackOption>)
  - [func NoOverwrite\(\) StackOption](<#NoOverwrite>)
- [type StackTracer](<#StackTracer>)


<a name="As"></a>
## func [As](<https://github.com/rclark/errors/blob/main/actions.go#L33>)

```go
func As(err error, target interface{}) bool
```

As finds the first error in err's tree that matches target, and if one is found, sets target to that error value and returns true. Otherwise, it returns false.

The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap\(\) error or Unwrap\(\) \[\]error method. When err wraps multiple errors, As examines err followed by a depth\-first traversal of its children.

An error matches target if the error's concrete value is assignable to the value pointed to by target, or if the error has a method As\(interface\{\}\) bool such that As\(target\) returns true. In the latter case, the As method is responsible for setting target.

An error type might provide an As method so it can be treated as if it were a different error type.

As panics if target is not a non\-nil pointer to either a type that implements error, or to any interface type.

<a name="Errorf"></a>
## func [Errorf](<https://github.com/rclark/errors/blob/main/actions.go#L157>)

```go
func Errorf(format string, args ...any) error
```

Errorf formats according to a format specifier and returns the string as a value that satisfies the error interface. By default, the error will have a stack trace to the point where the function was called.

If the format specifier includes a %w verb with an error operand, the returned error will implement an Unwrap method returning the operand. If there is more than one %w verb, the returned error will implement an Unwrap method returning a \[\]error containing all the %w operands in the order they appear in the arguments. It is invalid to supply the %w verb with an operand that does not implement the error interface. The %w verb is otherwise a synonym for %v.

If used to wrap errors, the [NoOverwrite](<#NoOverwrite>) option can be provided as the final argument to prevent the first stack trace from from one of the wrapped errors being overridden.

<details><summary>Example</summary>
<p>



```go
package main

import (
	std "errors"
	"fmt"

	"github.com/rclark/errors"
)

func withStack() error {
	return errors.New("first error")
}

func noStack() error {
	return std.New("no stack trace")
}

func main() {
	// When wrapped with no options, the stack trace should point to where
	// errors.Errorf is called from.
	original := withStack()
	err := errors.Errorf("wrapper: %w", original)
	stack, _ := errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// When wrapped with the NoOverwrite option, the stack trace should point to
	// the withStack function.
	original = withStack()
	err = errors.Errorf("wrapper: %w", original, errors.NoOverwrite())
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// When wrapped with the NoOverwrite, but the underlying error has no stack
	// trace, the resulting stack trace should point to where errors.Errorf is
	// called from.
	original = noStack()
	err = errors.Errorf("wrapper: %w", original, errors.NoOverwrite())
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

}
```

#### Output

```
github.com/rclark/errors_test.ExampleErrorf
github.com/rclark/errors_test.withStack
github.com/rclark/errors_test.ExampleErrorf
```

</p>
</details>

<a name="Is"></a>
## func [Is](<https://github.com/rclark/errors/blob/main/actions.go#L55>)

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
## func [Join](<https://github.com/rclark/errors/blob/main/actions.go#L65>)

```go
func Join(errs ...error) error
```

Join returns an error that wraps the given errors. Any nil error values are discarded. Join returns nil if every value in errs is nil. The error formats as the concatenation of the strings obtained by calling the Error method of each element of errs, with a newline between each string.

A non\-nil error returned by Join implements the Unwrap\(\) \[\]error method.

<a name="New"></a>
## func [New](<https://github.com/rclark/errors/blob/main/actions.go#L10>)

```go
func New(message string) error
```

New returns an error with the supplied message and a stack trace to the point where the function was called.

<a name="Unwrap"></a>
## func [Unwrap](<https://github.com/rclark/errors/blob/main/actions.go#L74>)

```go
func Unwrap(err error) error
```

Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.

Unwrap only calls a method of the form "Unwrap\(\) error". In particular Unwrap does not unwrap errors returned by [Join](<#Join>).

<a name="UnwrapAny"></a>
## func [UnwrapAny](<https://github.com/rclark/errors/blob/main/actions.go#L84>)

```go
func UnwrapAny(err error) []error
```

UnwrapAny returns the result of calling the Unwrap method on err, whether it implements \`Unwrap\(\) \[\]error\` or \`Unwrap\(\) error\`.

<a name="WithStack"></a>
## func [WithStack](<https://github.com/rclark/errors/blob/main/actions.go#L126>)

```go
func WithStack(err error, opts ...StackOption) error
```

WithStack adds a [Stack](<#Stack>) to the provided error at the point where the function was called. If the error already has a [Stack](<#Stack>), it will be overridden unless the [NoOverwrite](<#NoOverwrite>) option is provided.

<details><summary>Example</summary>
<p>



```go
package main

import (
	std "errors"
	"fmt"

	"github.com/rclark/errors"
)

func main() {
	// When the original error has no stack trace, the resulting stack trace
	// should point to where errors.WithStack is called from.
	original := noStack()
	err := errors.WithStack(original)
	stack, _ := errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// When the original error does have a stack trace, it will be overwritten.
	original = withStack()
	err = errors.WithStack(original)
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// With the NoOverwrite option, the original stack trace will be preserved.
	original = withStack()
	err = errors.WithStack(original, errors.NoOverwrite())
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

}

func withStack() error {
	return errors.New("first error")
}

func noStack() error {
	return std.New("no stack trace")
}
```

#### Output

```
github.com/rclark/errors_test.ExampleWithStack
github.com/rclark/errors_test.ExampleWithStack
github.com/rclark/errors_test.withStack
```

</p>
</details>

<a name="Error"></a>
## type [Error](<https://github.com/rclark/errors/blob/main/error.go#L9-L13>)

Error implements the error interface and provides a stack trace.

```go
type Error struct {
    // contains filtered or unexported fields
}
```

<a name="Error.Error"></a>
### func \(Error\) [Error](<https://github.com/rclark/errors/blob/main/error.go#L36>)

```go
func (e Error) Error() string
```

Error returns the error message.

<a name="Error.Format"></a>
### func \(Error\) [Format](<https://github.com/rclark/errors/blob/main/error.go#L56>)

```go
func (e Error) Format(s fmt.State, verb rune)
```

Format formats the error according to the fmt.Formatter interface.

- %s \<message\>
- %\+s \<message\>: \[\<filename:line\> ...\]
- %v \<message\>
- %\+v \<message\>\\n\<package\>.\<function\>\\n\\t\<filepath\>:\<line\>\\n\\t...

<a name="Error.StackTrace"></a>
### func \(Error\) [StackTrace](<https://github.com/rclark/errors/blob/main/error.go#L41>)

```go
func (e Error) StackTrace() Stack
```

StackTrace returns the [Stack](<#Stack>).

<a name="Error.Unwrap"></a>
### func \(Error\) [Unwrap](<https://github.com/rclark/errors/blob/main/error.go#L46>)

```go
func (e Error) Unwrap() error
```

Unwrap returns the wrapped error, if any.

<a name="Frame"></a>
## type [Frame](<https://github.com/rclark/errors/blob/main/frame.go#L12-L20>)

Frame represents a program counter inside a stack trace.

```go
type Frame struct {
    File string
    Line int

    Function string
    // contains filtered or unexported fields
}
```

<a name="Frame.Format"></a>
### func \(Frame\) [Format](<https://github.com/rclark/errors/blob/main/frame.go#L48>)

```go
func (f Frame) Format(s fmt.State, verb rune)
```

Format formats the frame according to the fmt.Formatter interface.

- %s \<filename\>:\<line\>
- %v \<package\>.\<function\>\\n\\t\<filepath\>:\<line\>

<a name="Frame.String"></a>
### func \(Frame\) [String](<https://github.com/rclark/errors/blob/main/frame.go#L59>)

```go
func (f Frame) String() string
```



<a name="Stack"></a>
## type [Stack](<https://github.com/rclark/errors/blob/main/stack-trace.go#L9>)

Stack represents a stack trace.

```go
type Stack []Frame
```

<a name="StackTrace"></a>
### func [StackTrace](<https://github.com/rclark/errors/blob/main/actions.go#L99>)

```go
func StackTrace(err error) (Stack, bool)
```

StackTrace returns a [Stack](<#Stack>), if err has one. If none was found, the returned bool will be false.

<a name="Stack.Format"></a>
### func \(Stack\) [Format](<https://github.com/rclark/errors/blob/main/stack-trace.go#L15>)

```go
func (st Stack) Format(s fmt.State, verb rune)
```

Format formats the stack of Frames according to the fmt.Formatter interface.

- %s \[\<filename\>:\<line\> ...\]
- %v \<package\>.\<function\>\\n\\t\<filepath\>:\<line\>\\n\\t...

<a name="Stack.IsZero"></a>
### func \(Stack\) [IsZero](<https://github.com/rclark/errors/blob/main/stack-trace.go#L39>)

```go
func (st Stack) IsZero() bool
```

IsZero reports whether the stack trace is empty.

<a name="StackOption"></a>
## type [StackOption](<https://github.com/rclark/errors/blob/main/actions.go#L113>)

StackOption is an option for the WithStack function.

```go
type StackOption func(*options)
```

<a name="NoOverwrite"></a>
### func [NoOverwrite](<https://github.com/rclark/errors/blob/main/actions.go#L117>)

```go
func NoOverwrite() StackOption
```

NoOverwrite is an option that prevents the stack trace from being overridden if the error already has one.

<a name="StackTracer"></a>
## type [StackTracer](<https://github.com/rclark/errors/blob/main/stack-trace.go#L45-L47>)

StackTracer is implemented by [Error](<#Error>). It can be used in external contexts to check whether an error has a stack trace that this package can expose.

```go
type StackTracer interface {
    StackTrace() Stack
}
```

<details><summary>Example</summary>
<p>



```go
package main

import (
	std "errors"
	"fmt"
	"log"

	"github.com/rclark/errors"
)

func main() {
	err := std.New("no stack trace")

	var hasStackTrace errors.StackTracer
	if errors.As(err, &hasStackTrace) {
		log.Fatal("error should not have a stack trace")
	}

	err = errors.New("with stack trace")
	if errors.As(err, &hasStackTrace) {
		fmt.Println(hasStackTrace.StackTrace()[0].Function)
	}
}
```

#### Output

```
github.com/rclark/errors_test.ExampleStackTracer
```

</p>
</details>

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
