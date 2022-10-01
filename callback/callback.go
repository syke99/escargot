package callback

import (
	"context"
	"errors"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

// CallBack is used to perform callback functions without specific
// *shell.Shell values
type CallBack[A any] struct {
	args []A
	cb   func(...A) *shell.Shell[A]
}

// CallBackX is used to perform callback functions on specific
// *shell.Shell values
type CallBackX[A any] struct {
	args []A
	cb   func(...A) *shell.Shell[A]
}

// CallBack executes the callback function provided just like
// callback.CallBack.Callback(), the only difference is it executes without
// a specific *shell.Shell value added to the arguments in the callback
// function
func (c CallBack[A]) CallBack() *shell.Shell[A] {
	return c.cb(c.args...)
}

// CallBackWithCancellation works just like CallBack, but takes a context
// to cancel execution
func (c CallBack[A]) CallBackWithCancellation(ctx context.Context, cancel context.CancelFunc) *shell.Shell[A] {

	select {
	default:
		res := c.cb(c.args...)

		if res.GetErrStatus() {
			cancel()
		}

		return res
	case <-ctx.Done():
		er := err.EscargotError{
			Level: "Cancel",
			Msg:   "context cancel signal received",
		}

		res := shell.Shell[A]{}

		r := &res

		r.Err(&er)

		return r
	}
}

// CallBackX executes the callback function provided on the value of
// the current iteration of Ranging over the *shell.Shell values, to
// execute a callback function without a *shell.Shell value added to
// the arguments in the callback function, use CallBackX
func (c CallBackX[A]) CallBackX(value A) *shell.Shell[A] {

	args := make([]A, len(c.args)+1)

	args[0] = value

	for i, v := range c.args {
		args[i+1] = v
	}

	return c.cb(args...)
}

// CallBackXWithCancellation works just like CallBackX, but takes a context
// to cancel execution
func (c CallBackX[A]) CallBackXWithCancellation(ctx context.Context, cancel context.CancelFunc, value A) *shell.Shell[A] {

	select {
	default:
		args := make([]A, len(c.args)+1)

		args[0] = value

		for i, v := range c.args {
			args[i+1] = v
		}

		res := c.cb(args...)

		if res.GetErrStatus() {
			cancel()
		}

		return res
	case <-ctx.Done():
		er := err.EscargotError{
			Level: "Cancel",
			Msg:   "context cancel signal received",
		}

		res := shell.Shell[A]{}

		r := &res

		r.Err(&er)

		return r
	}
}

// NewCallBackX returns a CallBack used in ranging over *shell.Shell values
// it takes the callback function expected, plus any arguments expected,
// minus the value. The callback function signature must match
// func(...any) *shell.Shell and the first argument will be the value
// in the current iteration of the range. All following arguments to the
// callback function will be the arguments provided
func NewCallBackX[A any](cb func(...A) *shell.Shell[A], args ...A) (CallBackX[A], error) {
	if cb == nil {
		return CallBackX[A]{}, errors.New("invalid callback configuration")
	}

	return CallBackX[A]{
		args: args,
		cb:   cb,
	}, nil
}

func NewCallBack[A any](cb func(...A) *shell.Shell[A], args ...A) (CallBack[A], error) {
	if cb == nil {
		return CallBack[A]{}, errors.New("invalid callback configuration")
	}

	return CallBack[A]{
		args: args,
		cb:   cb,
	}, nil
}
