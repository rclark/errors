package errors_test

import (
	std "errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/rclark/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func nextLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line + 1
}

func TestStackTrace(t *testing.T) {
	t.Run("with stack trace", func(t *testing.T) {
		line := nextLine()
		err := errors.New("with stack trace")
		stack, ok := errors.StackTrace(err)
		require.True(t, ok, "should have stack trace")

		found := fmt.Sprintf("%s", stack)
		expect := fmt.Sprintf("[actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should format properly")
	})

	t.Run("without stack trace", func(t *testing.T) {
		err := std.New("without stack trace")
		_, ok := errors.StackTrace(err)
		assert.False(t, ok, "should not have stack trace")
	})

	t.Run("joined errors", func(t *testing.T) {
		a := std.New("a")
		line := nextLine()
		b := errors.New("b")
		c := errors.New("c")

		err := errors.Join(a, b, c)
		stack, ok := errors.StackTrace(err)
		require.True(t, ok, "should have stack trace")

		found := fmt.Sprintf("%s", stack)
		expect := fmt.Sprintf("[actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should expose stack trace from first error that has one")
	})
}

func TestWithStack(t *testing.T) {
	t.Run("no prior stack", func(t *testing.T) {
		err := std.New("the message")
		line := nextLine()
		err = errors.WithStack(err)

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should contain frame at correct line")
	})

	t.Run("do not overwrite existing stack", func(t *testing.T) {
		line := nextLine()
		err := errors.New("the message")
		err = errors.WithStack(err)

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should contain frame at correct line")
	})

	t.Run("overwrite existing stack", func(t *testing.T) {
		err := errors.New("the message")
		line := nextLine()
		err = errors.WithStack(err, errors.Overwrite())

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should contain frame at correct line")
	})
}

func ExampleWithStack() {
	// When the original error has no stack trace, the resulting stack trace
	// should point to where errors.WithStack is called from.
	original := noStack()
	err := errors.WithStack(original)
	stack, _ := errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// When the original error does have a stack trace, it will not be
	// overwritten.
	original = withStack()
	err = errors.WithStack(original)
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// With the Overwrite option, the original stack trace will be overwritten.
	original = withStack()
	err = errors.WithStack(original, errors.Overwrite())
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// Output:
	// github.com/rclark/errors_test.ExampleWithStack
	// github.com/rclark/errors_test.withStack
	// github.com/rclark/errors_test.ExampleWithStack
}

func TestUnwrapAny(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var err error
		assert.Equal(t, []error{}, errors.UnwrapAny(err), "should return empty slice")
	})

	t.Run("single error", func(t *testing.T) {
		a := std.New("a")
		assert.Equal(t, []error{}, errors.UnwrapAny(a), "should return empty slice")
	})

	t.Run("fmt.Errorf", func(t *testing.T) {
		a := std.New("a")
		b := fmt.Errorf("b: %w", a)
		assert.Equal(t, []error{a}, errors.UnwrapAny(b), "should return unwrapped error")
	})

	t.Run("joined errors", func(t *testing.T) {
		a := std.New("a")
		b := errors.New("b")
		c := errors.New("c")
		err := errors.Join(a, b, c)

		assert.Equal(t, []error{a, b, c}, errors.UnwrapAny(err), "should return all joined errors")
	})
}

func TestErrorf(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		line := nextLine()
		err := errors.Errorf("the message")

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have correct stack trace")
	})

	t.Run("wraps error", func(t *testing.T) {
		err := std.New("wrapped message")
		line := nextLine()
		err = errors.Errorf("the message: %w", err)

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: wrapped message: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have correct stack trace")
	})

	t.Run("wraps many errors", func(t *testing.T) {
		a := std.New("a")
		b := std.New("b")
		c := std.New("c")
		line := nextLine()
		err := errors.Errorf("the message: %w: %w: %w", a, b, c)

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: a: b: c: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have correct stack trace")
	})

	t.Run("overwrite", func(t *testing.T) {
		err := errors.New("wrapped message")
		line := nextLine()
		err = errors.Errorf("the message: %w", err, errors.Overwrite())

		found := fmt.Sprintf("%+s", err)
		expect := fmt.Sprintf("the message: wrapped message: [actions_test.go:%d testing.go:", line)
		assert.Contains(t, found, expect, "should have correct stack trace")
	})
}

func TestAsAny(t *testing.T) {
	err := errors.New("base")
	err = errors.NewError[errors.BadInputError]("bad input", errors.FromError(err))
	err = errors.NewError[errors.ConflictError]("conflict", errors.FromError(err))
	err = errors.Errorf("wrapped: %w", err)

	var (
		badInput errors.BadInputError
		conflict errors.ConflictError
		missing  errors.MissingError
	)

	ok := errors.AsAny(err, &badInput, &conflict, &missing)
	assert.True(t, ok, "should find at least one match")
	assert.NotEmpty(t, badInput.Message(), "should find BadInputError")
	assert.NotEmpty(t, conflict.Message(), "should find ConflictError")
	assert.Empty(t, missing.Message(), "should not find MissingError")
}

func withStack() error {
	return errors.New("first error")
}

func noStack() error {
	return std.New("no stack trace")
}

func ExampleErrorf() {
	// When wrapped with no options, the stack trace should not be overwritten,
	// and point to the withStack function.
	original := withStack()
	err := errors.Errorf("wrapper: %w", original)
	stack, _ := errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// When wrapped with the Overwrite option, the stack trace should point to
	// where errors.Errorf is called from.
	original = withStack()
	err = errors.Errorf("wrapper: %w", original, errors.Overwrite())
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// When the underlying error has no stack trace, the resulting stack trace
	// should point to where errors.Errorf is called from.
	original = noStack()
	err = errors.Errorf("wrapper: %w", original)
	stack, _ = errors.StackTrace(err)
	fmt.Println(stack[0].Function)

	// Output:
	// github.com/rclark/errors_test.withStack
	// github.com/rclark/errors_test.ExampleErrorf
	// github.com/rclark/errors_test.ExampleErrorf
}
