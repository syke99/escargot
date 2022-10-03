package snail

import (
	"errors"
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
	shell   *shell.Shell
	try     func(args argument.Arguments) *shell.Shell
	catch   func(err *err.EscargotError, args argument.Arguments) *shell.Shell
	finally func(args argument.Arguments) *shell.Shell
}

// NewTrier will return a new Trier with the provided TryFunc and CatchFunc
func NewSnail(try TryFunc, catch CatchFunc, finally FinallyFunc) (Snail, error) {
	if try == nil ||
		catch == nil {
		return Snail{}, errors.New("invalid Trier configuration")
	}

	return Snail{
		try:     try,
		catch:   catch,
		finally: finally,
	}, nil
}

// Try tries the Trier's TryFunc with the provided tryArgs, and on error,
// will execute the Trier's CatchFunc with the provided catchArgs. It will
// return a *shell.Shell to access any values and/or errors
func (t *Snail) Try(tryArgs argument.Arguments, catchArgs argument.Arguments) {
	t.shell = t.try(tryArgs)
}

func (t *Snail) Catch(catchArgs argument.Arguments) {
	if t.shell.GetErrStatus() {
		t.shell = t.catch(t.shell.GetErr(), catchArgs)
	}
}

// Finally works just like Try, but executes a FinallyFunc after the TryFunc
// and/or CatchFunc, regardless of the outcome of either function
func (t *Snail) Finally(finallyArgs argument.Arguments) {
	for k, v := range t.shell.GetValues() {
		finallyArgs.SetArg(k, v, nil)
	}

	t.shell = t.finally(finallyArgs)
}
