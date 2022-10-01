package callback

import (
	"errors"
	"github.com/syke99/escargot/shell"
)

type CallBack struct {
	args []any
	cb   func(...any) *shell.Shell
}

// CallBack executes the callback function provided on the value of
// the current iteration of Ranging over the *shell.Shell values
func (c CallBack) CallBack(value any) *shell.Shell {

	args := make([]any, len(c.args)+1)

	args[0] = value

	for i, v := range c.args {
		args[i+1] = v
	}

	return c.cb(args...)
}

// NewCallBack returns a CallBack used in ranging over *shell.Shell values
// it takes the callback function expected, plus any arguments expected,
// minus the value. The callback function signature must match
// func(...any) *shell.Shell and the first argument will be the value
// in the current iteration of the range. All following arguments to the
// callback function will be the arguments provided
func NewCallBack(cb func(...any) *shell.Shell, args ...any) (CallBack, error) {
	if cb == nil {
		return CallBack{}, errors.New("invalid callback configuration")
	}

	return CallBack{
		args: args,
		cb:   cb,
	}, nil
}
