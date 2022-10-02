package callback

import (
	"context"
	"errors"
	"github.com/syke99/escargot/internal/build"
	"github.com/syke99/escargot/internal/override"
	"github.com/syke99/escargot/internal/resources"

	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
)

// CallBack is used to perform callback functions without specific
// *shell.Shell values
type CallBack struct {
	args argument.Arguments
	cb   func(arguments argument.Arguments) ([]any, *err.EscargotError)
}

// CallBackX is used to perform callback functions on specific
// *shell.Shell values
type CallBackX struct {
	args argument.Arguments
	cb   func(arguments argument.Arguments) ([]any, *err.EscargotError)
}

// CallBack executes the callback function provided just like
// callback.CallBack.Callback(), the only difference is it executes without
// a specific *shell.Shell value added to the arguments in the callback
// function
func (c CallBack) CallBack() ([]any, *err.EscargotError) {
	return c.cb(c.args)
}

// CallBackWithCancellation works just like CallBack, but takes a context
// to cancel execution
func (c CallBack) CallBackWithCancellation(ctx context.Context, cancel context.CancelFunc) ([]any, *err.EscargotError) {

	select {
	default:
		res, er := c.cb(c.args)

		if er.Unwrap() != nil {
			cancel()
		}

		return res, er
	case <-ctx.Done():
		return nil, build.BuildErr(resources.ContextCancel, resources.ErrLevel, resources.CanceledContext.String())
	}
}

// CallBackX executes the callback function provided on the value of
// the current iteration of Ranging over the *shell.Shell values, to
// execute a callback function without a *shell.Shell value added to
// the arguments in the callback function, use CallBackX
func (c CallBackX) CallBackX(key string, value any, overRide *override.OverRider) ([]any, *err.EscargotError) {
	c.args.SetArg(key, value, overRide)

	return c.cb(c.args)
}

// CallBackXWithCancellation works just like CallBackX, but takes a context
// to cancel execution
func (c CallBackX) CallBackXWithCancellation(ctx context.Context, cancel context.CancelFunc, key string, value any, overRide *override.OverRider) ([]any, *err.EscargotError) {

	select {
	default:
		c.args.SetArg(key, value, overRide)

		res, er := c.cb(c.args)

		if er.Unwrap() != nil {
			cancel()
		}

		return res, er
	case <-ctx.Done():
		return nil, build.BuildErr(resources.ContextCancel, resources.ErrLevel, resources.CanceledContext.String())
	}
}

// NewCallBackX returns a CallBackX used in ranging over *shell.Shell values
// it takes the callback function expected, plus any arguments expected,
// minus the value. The callback function signature must match
// func(...any) *shell.Shell and the first argument will be the value
// in the current iteration of the range. All following arguments to the
// callback function will be the arguments provided
func NewCallBackX(cb func(arguments argument.Arguments) ([]any, *err.EscargotError), args argument.Arguments) (CallBackX, error) {
	if cb == nil {
		return CallBackX{}, errors.New("invalid callback configuration")
	}

	return CallBackX{
		args: args,
		cb:   cb,
	}, nil
}

func NewCallBack(cb func(arguments argument.Arguments) ([]any, *err.EscargotError), args argument.Arguments) (CallBack, error) {
	if cb == nil {
		return CallBack{}, errors.New("invalid callback configuration")
	}

	return CallBack{
		args: args,
		cb:   cb,
	}, nil
}
