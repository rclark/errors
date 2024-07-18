[![Go](https://github.com/rclark/errors/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/rclark/errors/actions/workflows/go.yml)

# errors

Another package for errors with stack traces included.

In addition to methods like `.As`, `.Is`, and `.Join` that you're used to from the standard `errors` package, this package makes sure that:

- `errors.New` always creates an error with a stack trace from where it is called.
- `errors.Errorf` is a drop-in replacement for `fmt.Errorf`, resulting in an error with a stack trace from where it is called.
- `errors.WithStack` adds a stack trace to an error that may not already have one.
- `errors.StackTrace` returns the stack trace from an error, if it has one.
- If an error has a stack trace, printing is with the `%+s` or `%+v` formatting directive will write out the error message and the stack trace.

The package also provides a number of error types that can be used to represent common types of application failure situations. These errors can carry two messages, one with technical details (exposed by `.Error()`) and another with a description better suited for a user of the application (exposed by `.Message()`).

```go
package example

import (
  "fmt"

  "github.com/rclark/errors"
)

func main() {
  err := errors.New("something went wrong")
  fmt.Printf("%+v", err)

  // Output:
  // something went wrong
  // example.main
  //  /path/to/example.go:10
}
```

## Usage

See [usage.md](./usage.md).
