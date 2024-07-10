package errors_test

import (
	std "errors"
	"fmt"
	"log"

	"github.com/rclark/errors"
)

func ExampleStackTracer() {
	err := std.New("no stack trace")

	var hasStackTrace errors.StackTracer
	if errors.As(err, &hasStackTrace) {
		log.Fatal("error should not have a stack trace")
	}

	err = errors.New("with stack trace")
	if errors.As(err, &hasStackTrace) {
		fmt.Println(hasStackTrace.StackTrace()[0].Function)
	}
	// Output: github.com/rclark/errors_test.ExampleStackTracer
}
