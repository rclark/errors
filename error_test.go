package errors_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rclark/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorFormat(t *testing.T) {
	line := nextLine()
	e := errors.New("the error message")
	buf := strings.Builder{}
	fmt.Fprintf(&buf, "%+v", e)

	lines := strings.Split(buf.String(), "\n")
	assert.Equal(t, "the error message", lines[0], "%v first line should be error message")
	assert.Equal(t, "github.com/rclark/errors_test.TestErrorFormat", lines[1], "%v second line should be the test function")

	expect := fmt.Sprintf("error_test.go:%d", line)
	assert.Contains(t, lines[2], expect, "%v third line should be the test file path & line number")
	assert.Equal(t, "testing.tRunner", lines[3], "%v fourth line should be the test runner")
	assert.Contains(t, lines[4], "/testing.go:", "%v fifth line should be the test runner file path & line number")

	buf = strings.Builder{}
	fmt.Fprintf(&buf, "%+s", e)
	expect = fmt.Sprintf("the error message: [error_test.go:%d testing.go:", line)
	assert.Contains(t, buf.String(), expect, "%s should contain the error message and file path")
}
