package callback

import (
	"context"
	"errors"

	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
)

// CallBack is used to perform callback functions without specific
// *shell.Shell values
type CallBack struct {
	args  []any
	argsx argument.Arguments
	cb    func(...any) (any, *err.EscargotError)
}

// CallBackX is used to perform callback functions on specific
// *shell.Shell values
type CallBackX struct {
	args  []any
	argsx argument.Arguments
	cb    func(...any) (any, *err.EscargotError)
}

// CallBack executes the callback function provided just like
// callback.CallBack.Callback(), the only difference is it executes without
// a specific *shell.Shell value added to the arguments in the callback
// function
func (c CallBack) CallBack() (any, *err.EscargotError) {
	args := make([]any, len(c.argsx.GetArgsSlice()))

	for _, v := range c.argsx.GetArgsSlice() {
		args = append(args, v)
	}

	return c.cb(c.args...)
}

// CallBackWithCancellation works just like CallBack, but takes a context
// to cancel execution
func (c CallBack) CallBackWithCancellation(ctx context.Context, cancel context.CancelFunc) (any, *err.EscargotError) {

	select {
	default:
		res, er := c.cb(c.args...)

		if er.Unwrap() != nil {
			cancel()
		}

		return res, er
	case <-ctx.Done():
		er := err.EscargotError{
			Level: "Cancel",
			Msg:   "context cancel signal received",
		}

		return nil, &er
	}
}

// CallBackX executes the callback function provided on the value of
// the current iteration of Ranging over the *shell.Shell values, to
// execute a callback function without a *shell.Shell value added to
// the arguments in the callback function, use CallBackX
func (c CallBackX) CallBackX(value any) (any, *err.EscargotError) {

	args := make([]any, len(c.argsx.GetArgsSlice())+1)

	args[0] = value

	for i, v := range c.argsx.GetArgsSlice() {
		args[i+1] = v
	}

	return c.cb(args...)
}

// CallBackXWithCancellation works just like CallBackX, but takes a context
// to cancel execution
func (c CallBackX) CallBackXWithCancellation(ctx context.Context, cancel context.CancelFunc, value any) (any, *err.EscargotError) {

	select {
	default:
		args := make([]any, len(c.args)+1)

		args[0] = value

		for i, v := range c.args {
			args[i+1] = v
		}

		res, er := c.cb(args...)

		if er.Unwrap() != nil {
			cancel()
		}

		return res, er
	case <-ctx.Done():
		er := err.EscargotError{
			Level: "Cancel",
			Msg:   "context cancel signal received",
		}

		return nil, &er
	}
}

// NewCallBackX returns a CallBackX used in ranging over *shell.Shell values
// it takes the callback function expected, plus any arguments expected,
// minus the value. The callback function signature must match
// func(...any) *shell.Shell and the first argument will be the value
// in the current iteration of the range. All following arguments to the
// callback function will be the arguments provided
func NewCallBackX(cb func(...any) (any, *err.EscargotError), args ...any) (CallBackX, error) {
	if cb == nil {
		return CallBackX{}, errors.New("invalid callback configuration")
	}

	return CallBackX{
		args: args,
		cb:   cb,
	}, nil
}

func NewCallBack(cb func(...any) (any, *err.EscargotError), args ...any) (CallBack, error) {
	if cb == nil {
		return CallBack{}, errors.New("invalid callback configuration")
	}

	return CallBack{
		args: args,
		cb:   cb,
	}, nil
}
