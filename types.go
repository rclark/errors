package errors

type tracedError interface {
	StackTracer
	error
}

// UserFacingError is an error that carries a message that has been designated
// to be shown to a user external to the system.
//
// This is particularly useful in situations where you want an error to carry
// detailed, technical information to a logging system (via .Error()), while
// also carrying a more user-friendly description of the failure
// (via .Message()) to another part of the  application where it will be
// returned to the user.
type UserFacingError struct {
	err tracedError
	msg string
}

// UserFacingOption configures the creation of a [UserFacingError].
type UserFacingOption func(*options)

// FromError sets the [UserFacingError] to wrap the provided error.
func FromError(err error) UserFacingOption {
	return func(o *options) {
		o.underlying = err
	}
}

// OverwriteStackTrace sets the stack trace of a [UserFacingError] to the place that
// [NewUserFacingError] was called, overwriting any stack trace that may have
// been included in an underlying error provided via [FromError].
func OverwriteStackTrace() UserFacingOption {
	return func(o *options) {
		o.noOverwrite = false
	}
}

// Skip sets the number of stack frames to skip when creating a
// [UserFacingError].
func Skip(i int) UserFacingOption {
	return func(o *options) {
		o.skip = i
	}
}

// NewUserFacingError creates a new [UserFacingError]. The provided message is
// meant to be shown to a user external to the system. If no error is provided
// via [FromError], the provided message will also be used as the underlying
// error message.
func NewUserFacingError(msg string, opts ...UserFacingOption) error {
	uf := UserFacingError{msg: msg}

	o := options{noOverwrite: true, skip: 3}
	for _, opt := range opts {
		opt(&o)
	}

	if o.underlying == nil {
		uf.err = newError(msg, o.skip)
		return uf
	}

	var te tracedError
	if !As(o.underlying, &te) {
		uf.err = wrapError(o.underlying, o.skip)
		return uf
	}

	if o.noOverwrite {
		uf.err = te
		return uf
	}

	uf.err = wrapError(o.underlying, o.skip)
	return uf
}

// StackTrace returns the [Stack].
func (uf UserFacingError) StackTrace() Stack {
	return uf.err.StackTrace()
}

// Unwrap returns the underlying error, if any.
func (uf UserFacingError) Unwrap() error {
	return uf.err
}

// Error returns the underlying error message.
func (uf UserFacingError) Error() string {
	return uf.err.Error()
}

// Message returns the error message intended for the user external to the
// system.
func (uf UserFacingError) Message() string {
	return uf.msg
}

type userFacing interface {
	Message() string
}

// UserFacingMessage returns a message intended for a user external to the
// system, if the error provides one.
func UserFacingMessage(err error) (string, bool) {
	var uf userFacing
	if As(err, &uf) {
		return uf.Message(), true
	}

	return "", false
}

// ErrorType are generalized categories of errors that can be used to represent
// different kinds of common application failures. Using categories like this
// can help to provide more context to callers about how they may wish to handle
// the error.
type ErrorType interface {
	BadInputError | NotAllowedError | MissingError | ConflictError | TimeoutError | UnexpectedError
}

// NewError creates a new error of the provided generic type with the given
// message intended for a user external to the system.
func NewError[T ErrorType](msg string, opts ...UserFacingOption) error {
	opts = append([]UserFacingOption{Skip(4)}, opts...)
	uf := NewUserFacingError(msg, opts...).(UserFacingError)
	return error(T{UserFacingError: uf})
}

func asType[T ErrorType](err error) (T, bool) {
	var e T
	return e, As(err, &e)
}

// BadInputError is an [ErrorType] that represents a situation where some input was
// invalid.
type BadInputError struct {
	UserFacingError
}

// IsBadInput reports whether the provided error is a [BadInputError] and
// returns it if so.
func IsBadInput(err error) (BadInputError, bool) {
	return asType[BadInputError](err)
}

// NotAllowedError is an [ErrorType] that represents a situation where some action was
// not allowed.
type NotAllowedError struct {
	UserFacingError
}

// IsNotAllowed reports whether the provided error is a [NotAllowedError] and
// returns it if so.
func IsNotAllowed(err error) (NotAllowedError, bool) {
	return asType[NotAllowedError](err)
}

// MissingError is an [ErrorType] that represents a situation where something was
// not found.
type MissingError struct {
	UserFacingError
}

// IsMissing reports whether the provided error is a [MissingError] and returns
// it if so.
func IsMissing(err error) (MissingError, bool) {
	return asType[MissingError](err)
}

// ConflictError is an [ErrorType] that represents a situation where some action
// could not be completed due to a conflict.
type ConflictError struct {
	UserFacingError
}

// IsConflict reports whether the provided error is a [ConflictError] and
// returns it if so.
func IsConflict(err error) (ConflictError, bool) {
	return asType[ConflictError](err)
}

// TimeoutError is an [ErrorType] that represents a situation where some action took
// too long to complete.
type TimeoutError struct {
	UserFacingError
}

// IsTimeout reports whether the provided error is a [TimeoutError] and returns
// it if so.
func IsTimeout(err error) (TimeoutError, bool) {
	return asType[TimeoutError](err)
}

// UnexpectedError is an [ErrorType] that represents a situation where an unexpected
// error occurred.
type UnexpectedError struct {
	UserFacingError
}

// IsUnexpected reports whether the provided error is an [UnexpectedError] and
// returns it if so.
func IsUnexpected(err error) (UnexpectedError, bool) {
	return asType[UnexpectedError](err)
}
