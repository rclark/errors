package errors

import (
	"errors"
	"fmt"
)

// New returns an error with the supplied message and a stack trace to the point
// where the function was called.
func New(message string) error {
	return newError(message, 3)
}

// As finds the first error in err's tree that matches target, and if one is
// found, sets target to that error value and returns true. Otherwise, it
// returns false.
//
// The tree consists of err itself, followed by the errors obtained by
// repeatedly calling its Unwrap() error or Unwrap() []error method. When err
// wraps multiple errors, As examines err followed by a depth-first traversal of
// its children.
//
// An error matches target if the error's concrete value is assignable to the
// value pointed to by target, or if the error has a method As(interface{}) bool
// such that As(target) returns true. In the latter case, the As method is
// responsible for setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// AsAny runs [As] for each provided targets. It will return true if it finds a
// match for at least one of the targets. Otherwise, it will return false. The
// targets that match will be set to the first error in the tree that matches.
func AsAny(err error, targets ...interface{}) bool {
	result := false

	for i := range targets {
		if As(err, targets[i]) {
			result = true
		}
	}

	return result
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by
// repeatedly calling its Unwrap() error or Unwrap() []error method. When err
// wraps multiple errors, Is examines err followed by a depth-first traversal of
// its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for an
// example in the standard library. An Is method should only shallowly compare
// err and the target and not call [Unwrap] on either.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Join returns an error that wraps the given errors. Any nil error values are
// discarded. Join returns nil if every value in errs is nil. The error formats
// as the concatenation of the strings obtained by calling the Error method of
// each element of errs, with a newline between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type
// contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
//
// Unwrap only calls a method of the form "Unwrap() error". In particular Unwrap
// does not unwrap errors returned by [Join].
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

type joined interface {
	Unwrap() []error
}

// UnwrapAny returns the result of calling the Unwrap method on err, whether
// it implements `Unwrap() []error` or `Unwrap() error`.
func UnwrapAny(err error) []error {
	var u joined
	if As(err, &u) {
		return u.Unwrap()
	}

	if e := errors.Unwrap(err); e != nil {
		return []error{e}
	}

	return make([]error, 0)
}

// StackTrace returns a [Stack], if err has one. If none was found, the returned
// bool will be false.
func StackTrace(err error) (Stack, bool) {
	var s StackTracer
	if As(err, &s) {
		return s.StackTrace(), true
	}

	return make(Stack, 0), false
}

type options struct {
	noOverwrite bool
	underlying  error
	skip        int
}

// StackOption is an option for the WithStack function.
type StackOption func(*options)

// NoOverwrite is an option that prevents the stack trace from being overridden
// if the error already has one.
func NoOverwrite() StackOption {
	return func(o *options) {
		o.noOverwrite = true
	}
}

// WithStack adds a [Stack] to the provided error at the point where the
// function was called. If the error already has a [Stack], it will be
// overridden unless the [NoOverwrite] option is provided.
func WithStack(err error, opts ...StackOption) error {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.noOverwrite {
		var s StackTracer
		if As(err, &s) {
			return err
		}
	}

	return wrapError(err, 3)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies the error interface. By default, the error will have a
// stack trace to the point where the function was called.
//
// If the format specifier includes a %w verb with an error operand, the
// returned error will implement an Unwrap method returning the operand. If
// there is more than one %w verb, the returned error will implement an Unwrap
// method returning a []error containing all the %w operands in the order they
// appear in the arguments. It is invalid to supply the %w verb with an operand
// that does not implement the error interface. The %w verb is otherwise a
// synonym for %v.
//
// If used to wrap errors, the [NoOverwrite] option can be provided as the final
// argument to prevent the first stack trace from from one of the wrapped errors
// being overridden.
func Errorf(format string, args ...any) error {
	options := options{}

	operands := []any{}

	for i := len(args) - 1; i >= 0; i-- {
		if opt, ok := args[i].(StackOption); ok {
			opt(&options)
		} else {
			operands = args[:i+1]
			break
		}
	}

	e := fmt.Errorf(format, operands...)

	if options.noOverwrite {
		var s StackTracer
		if As(e, &s) {
			return Error{message: e.Error(), err: e, stack: s.StackTrace()}
		}
	}

	return wrapError(e, 3)
}
