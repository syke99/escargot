package callback

import (
	"errors"
	"github.com/syke99/escargot/shell"
)

// CallBackX is used to perform callback functions on specific
// *shell.Shell values
type CallBackX struct {
	args []any
	cb   func(...any) *shell.Shell
}

// CallBack is used to perform callback functions without specific
// *shell.Shell values
type CallBack struct {
	args []any
	cb   func(...any) *shell.Shell
}

// CallBackX executes the callback function provided on the value of
// the current iteration of Ranging over the *shell.Shell values, to
// execute a callback function without a *shell.Shell value added to
// the arguments in the callback function, use CallBackX
func (c CallBackX) CallBackX(value any) *shell.Shell {

	args := make([]any, len(c.args)+1)

	args[0] = value

	for i, v := range c.args {
		args[i+1] = v
	}

	return c.cb(args...)
}

// CallBack executes the callback function provided just like
// callback.CallBack.Callback(), the only difference is it executes without
// a specific *shell.Shell value added to the arguments in the callback
// function
func (c CallBack) CallBack() *shell.Shell {
	return c.cb(c.args)
}

// NewCallBackX returns a CallBack used in ranging over *shell.Shell values
// it takes the callback function expected, plus any arguments expected,
// minus the value. The callback function signature must match
// func(...any) *shell.Shell and the first argument will be the value
// in the current iteration of the range. All following arguments to the
// callback function will be the arguments provided
func NewCallBackX(cb func(...any) *shell.Shell, args ...any) (CallBack, error) {
	if cb == nil {
		return CallBack{}, errors.New("invalid callback configuration")
	}

	return CallBack{
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
