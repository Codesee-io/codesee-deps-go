package errutils

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"
)

const stackSize = 4 << 10 // 4KB

type unwrapper interface {
	Unwrap() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Fatal prints an error, prints its stack trace, and then exits with a 1 exit
// code.
func Fatal(err error) {
	stack := ErrStack(err)
	fmt.Println(err.Error())
	fmt.Println(string(stack))
	os.Exit(1)
}

// UnwrapErr takes in an error, goes through the error cause chain, and returns the
// deepest error that has a stacktrace. This is needed because by default, every
// subsequent call to errors.Wrap overwrites the previous stacktrace. So in the
// end, we don't know where the error originated from. By going to the deepest
// one, we can find exactly where it started.
func UnwrapErr(err error) error {
	for {
		var u unwrapper
		if !errors.As(err, &u) {
			break
		}
		unwrapped := u.Unwrap()
		var st stackTracer
		if !errors.As(unwrapped, &st) {
			break
		}
		err = unwrapped
	}
	return err
}

// ErrStack takes in an error and returns the stack trace in a byte slice. Go
// errors don't natural keep track of stack info, but the pkg/errors package
// does, so we pull the stack trace using that info.
func ErrStack(err error) []byte {
	var stack []byte

	// This is to support pkg/errors stackTracer interface.
	err = UnwrapErr(err)
	var st stackTracer
	if errors.As(err, &st) {
		stack = []byte(fmt.Sprintf("%+v", st.StackTrace()))
	} else {
		stack = make([]byte, stackSize)
		n := runtime.Stack(stack, true)
		stack = stack[:n]
	}

	return stack
}
