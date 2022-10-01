package callback

import (
	"context"
	"errors"
	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

// CallBack is used to perform callback functions without specific
// *shell.Shell values
type CallBack struct {
	args  []any
	argsx argument.Arguments
	cb    func(...any) *shell.Shell
}

// CallBackX is used to perform callback functions on specific
// *shell.Shell values
type CallBackX struct {
	args  []any
	argsx argument.Arguments
	cb    func(...any) *shell.Shell
}

// CallBack executes the callback function provided just like
// callback.CallBack.Callback(), the only difference is it executes without
// a specific *shell.Shell value added to the arguments in the callback
// function
func (c CallBack) CallBack() *shell.Shell {
	args := make([]any, c.argsx.GetArgsLength())

	for _, v := range c.argsx.GetArgsSlice() {
		args = append(args, v)
	}

	return c.cb(c.args...)
}

// CallBackWithCancellation works just like CallBack, but takes a context
// to cancel execution
func (c CallBack) CallBackWithCancellation(ctx context.Context, cancel context.CancelFunc) *shell.Shell {

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

		res := shell.Shell{}

		r := &res

		r.Err(&er)

		return r
	}
}

// CallBackX executes the callback function provided on the value of
// the current iteration of Ranging over the *shell.Shell values, to
// execute a callback function without a *shell.Shell value added to
// the arguments in the callback function, use CallBackX
func (c CallBackX) CallBackX(value any) *shell.Shell {

	args := make([]any, c.argsx.GetArgsLength()+1)

	args[0] = value

	for i, v := range c.argsx.GetArgsSlice() {
		args[i+1] = v
	}

	return c.cb(args...)
}

// CallBackXWithCancellation works just like CallBackX, but takes a context
// to cancel execution
func (c CallBackX) CallBackXWithCancellation(ctx context.Context, cancel context.CancelFunc, value any) *shell.Shell {

	select {
	default:
		args := make([]any, len(c.args)+1)

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

		res := shell.Shell{}

		r := &res

		r.Err(&er)

		return r
	}
}

// NewCallBackX returns a CallBackX used in ranging over *shell.Shell values
// it takes the callback function expected, plus any arguments expected,
// minus the value. The callback function signature must match
// func(...any) *shell.Shell and the first argument will be the value
// in the current iteration of the range. All following arguments to the
// callback function will be the arguments provided
func NewCallBackX(cb func(...any) *shell.Shell, args ...any) (CallBackX, error) {
	if cb == nil {
		return CallBackX{}, errors.New("invalid callback configuration")
	}

	return CallBackX{
		args: args,
		cb:   cb,
	}, nil
}

func NewCallBack(cb func(...any) *shell.Shell, args ...any) (CallBack, error) {
	if cb == nil {
		return CallBack{}, errors.New("invalid callback configuration")
	}

	return CallBack{
		args: args,
		cb:   cb,
	}, nil
}
