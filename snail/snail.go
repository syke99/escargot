package snail

import (
	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

type TryFunc func(args argument.Arguments) *shell.Shell

type CatchFunc func(err *err.EscargotError, args argument.Arguments) *shell.Shell

type FinallyFunc func(args argument.Arguments) *shell.Shell

// Snail will handle trying the TryFunc provided and execute the provided CatchFunc
// on error
type Snail struct {
	shell *shell.Shell
}

// NewTrier will return a new Snail for chaining together a TryFunc, CatchFunc, and/or FinallyFunc
func NewSnail() *Snail {
	shl := shell.Shell{}

	snl := Snail{
		shell: &shl,
	}

	return &snl
}

// Try tries the provided TryFunc with the provided argument.Arguments
// and sets the value(s) for the Snail to use during Catch and/or Finally
func (t *Snail) Try(tryFunc TryFunc, tryArgs argument.Arguments) *Snail {
	t.shell = tryFunc(tryArgs)

	return t
}

// Catch checks if and error occurred, and if so, executes the provided CatchFunc
// with the provided argument.Arguments, and well as any currently set value(s)
// in the Snail
func (t *Snail) Catch(catchFunc CatchFunc, catchArgs argument.Arguments) *Snail {
	if t.shell.GetErrStatus() {
		for k, v := range t.shell.GetValues() {
			er := catchArgs.SetArg(k, v, nil)

			if er != nil {
				t.shell.Err(er, "")
			}
		}

		t.shell = catchFunc(t.shell.GetErr(), catchArgs)
	}

	return t
}

// Finally executes the provided FinallyFunc with the provided argument.Arguments
// along with any currently set value(s) in the Snail, regardless of error status
func (t *Snail) Finally(finallyFunc FinallyFunc, finallyArgs argument.Arguments) *Snail {
	for k, v := range t.shell.GetValues() {
		er := finallyArgs.SetArg(k, v, nil)

		if er != nil {
			t.shell.Err(er, "")
		}
	}

	t.shell = finallyFunc(finallyArgs)

	return t
}
