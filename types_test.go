package errors_test

import (
	std "errors"
	"fmt"
	"log"
	"testing"

	"github.com/rclark/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorTypes(t *testing.T) {
	t.Run("BadInputError", func(t *testing.T) {
		line := nextLine()
		err := errors.NewError[errors.BadInputError]("bad input")
		assert.Equal(t, "bad input", err.Error(), "error message should match")

		_, ok := errors.IsBadInput(err)
		assert.True(t, ok, "expected error to be of type BadInputError")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("NotAllowedError", func(t *testing.T) {
		line := nextLine()
		err := errors.NewError[errors.NotAllowedError]("not allowed")
		assert.Equal(t, "not allowed", err.Error(), "error message should match")

		_, ok := errors.IsNotAllowed(err)
		assert.True(t, ok, "expected error to be of type NotAllowedError")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("MissingError", func(t *testing.T) {
		line := nextLine()
		err := errors.NewError[errors.MissingError]("missing")
		assert.Equal(t, "missing", err.Error(), "error message should match")

		_, ok := errors.IsMissing(err)
		assert.True(t, ok, "expected error to be of type MissingError")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("ConflictError", func(t *testing.T) {
		line := nextLine()
		err := errors.NewError[errors.ConflictError]("conflict")
		assert.Equal(t, "conflict", err.Error(), "error message should match")

		_, ok := errors.IsConflict(err)
		assert.True(t, ok, "expected error to be of type ConflictError")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("TimeoutError", func(t *testing.T) {
		line := nextLine()
		err := errors.NewError[errors.TimeoutError]("timeout")
		assert.Equal(t, "timeout", err.Error(), "error message should match")

		_, ok := errors.IsTimeout(err)
		assert.True(t, ok, "expected error to be of type TimeoutError")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		line := nextLine()
		err := errors.NewError[errors.UnexpectedError]("unexpected")
		assert.Equal(t, "unexpected", err.Error(), "error message should match")

		_, ok := errors.IsUnexpected(err)
		assert.True(t, ok, "expected error to be of type UnexpectedError")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("wrapping an underlying error", func(t *testing.T) {
		line := nextLine()
		err := errors.New("underlying error")

		err = errors.NewError[errors.BadInputError]("bad input", errors.FromError(err))
		msg, ok := errors.UserFacingMessage(err)
		require.True(t, ok, "expected error to be a UserFacingError")
		assert.Equal(t, "bad input", msg, "external message should match")
		assert.Equal(t, "underlying error", err.Error(), "underlying error message should match")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("wrapping an underlying error with no stack trace", func(t *testing.T) {
		err := std.New("underlying error")

		line := nextLine()
		err = errors.NewError[errors.BadInputError]("bad input", errors.FromError(err))
		msg, ok := errors.UserFacingMessage(err)
		require.True(t, ok, "expected error to be a UserFacingError")
		assert.Equal(t, "bad input", msg, "external message should match")
		assert.Equal(t, "underlying error", err.Error(), "underlying error message should match")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})

	t.Run("wrapping an underlying error, overwrite stack trace", func(t *testing.T) {
		err := errors.New("underlying error")

		line := nextLine()
		err = errors.NewError[errors.BadInputError]("bad input", errors.FromError(err), errors.OverwriteStackTrace())
		msg, ok := errors.UserFacingMessage(err)
		require.True(t, ok, "expected error to be a UserFacingError")
		assert.Equal(t, "bad input", msg, "external message should match")
		assert.Equal(t, "underlying error", err.Error(), "underlying error message should match")

		trace, ok := errors.StackTrace(err)
		assert.True(t, ok, "expected error to have a stack trace")

		found := fmt.Sprintf("%s", trace)
		expect := fmt.Sprintf("[types_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have stack trace at correct location")
	})
}

func ExampleUserFacingError() {
	underlying := errors.New("failed to decode: string is not valid utf-8")
	err := errors.NewUserFacingError("string included invalid characters", errors.FromError(underlying))

	fmt.Println(err.Error())

	msg, ok := errors.UserFacingMessage(err)
	if !ok {
		log.Fatal("expected error to be a UserFacingError")
	}

	fmt.Println(msg)

	// Output:
	// failed to decode: string is not valid utf-8
	// string included invalid characters
}

func ExampleNewError() {
	underlying := errors.New("failed to decode: string is not valid utf-8")
	err := errors.NewError[errors.BadInputError]("invalid characters", errors.FromError(underlying))

	bad, ok := errors.IsBadInput(err)
	if !ok {
		log.Fatal("expected error to represent bad input")
	}

	fmt.Println(bad.Message())
	fmt.Println(bad.Error())
	fmt.Printf("%s", bad.StackTrace()[0])

	// Output:
	// invalid characters
	// failed to decode: string is not valid utf-8
	// types_test.go:167
}
