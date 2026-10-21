package errors

import (
	"errors"

	e "github.com/pkg/errors"
)

var (
	New         = e.New
	Errorf      = e.Errorf
	Wrap        = e.Wrap
	Wrapf       = e.Wrapf
	Cause       = e.Cause
	WithStack   = e.WithStack
	WithMessage = e.WithMessage

	Join   = errors.Join
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

func AsType[E error](err error) (E, bool) {
	return errors.AsType[E](err)
}

type stackTracer interface {
	StackTrace() e.StackTrace
}

// Stacks returns the stack trace of a error that wrapped via withStack.
func Stacks(err error) e.StackTrace {
	if er, ok := err.(stackTracer); ok {
		return er.StackTrace()
	}
	return e.StackTrace{}
}
