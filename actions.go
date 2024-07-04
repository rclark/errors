package errors

import "errors"

// New returns an error with the supplied message and a stack trace to the point
// where the function was called.
func New(message string) Error {
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

// StackTrace returns a [Stack], if err has one. If none was found, the returned
// [Stack] will be empty. This can be checked with the [Stack.IsZero] method.
func StackTrace(err error) Stack {
	var s StackTracer
	if As(err, &s) {
		return s.StackTrace()
	}

	return make(Stack, 0)
}

type options struct {
	noOverwrite bool
}

// WithStackOption is an option for the WithStack function.
type WithStackOption func(*options)

// NoOverwrite is an option that prevents the stack trace from being overridden
// if the error already has one.
func NoOverwrite() WithStackOption {
	return func(o *options) {
		o.noOverwrite = true
	}
}

// WithStack adds a [Stack] to the provided error at the point where the
// function was called. If the error already has a [Stack], it will be
// overridden unless the [NoOverwrite] option is provided.
func WithStack(err error, opts ...WithStackOption) error {
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