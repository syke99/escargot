package try

import (
	"errors"
	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

type TryFunc func(args argument.Arguments) *shell.Shell

type CatchFunc func(err *err.EscargotError, args argument.Arguments) *shell.Shell

type FinallyFunc func(args argument.Arguments) *shell.Shell

// Trier will handle trying the TryFunc provided and execute the provided CatchFunc
// on error
type Trier struct {
	try     func(args argument.Arguments) *shell.Shell
	catch   func(err *err.EscargotError, args argument.Arguments) *shell.Shell
	finally func(args argument.Arguments) *shell.Shell
}

// NewTrier will return a new Trier with the provided TryFunc and CatchFunc
func NewTrier(try TryFunc, catch CatchFunc, finally FinallyFunc) (Trier, error) {
	if try == nil ||
		catch == nil {
		return Trier{}, errors.New("invalid Trier configuration")
	}

	return Trier{
		try:     try,
		catch:   catch,
		finally: finally,
	}, nil
}

// Try tries the Trier's TryFunc with the provided tryArgs, and on error,
// will execute the Trier's CatchFunc with the provided catchArgs. It will
// return a *shell.Shell to access any values and/or errors
func (t Trier) Try(tryArgs argument.Arguments, catchArgs argument.Arguments) *shell.Shell {
	result := t.try(tryArgs)

	if result.GetErrStatus() {
		result = t.catch(result.GetErr(), catchArgs)
	}

	return result
}

// TryFinally works just like Try, but executes a FinallyFunc after the TryFunc
// and/or CatchFunc, regardless of the outcome of either function
func (t Trier) TryFinally(tryArgs argument.Arguments, catchArgs argument.Arguments, finallyArgs argument.Arguments) *shell.Shell {
	result := t.try(tryArgs)

	if result.GetErrStatus() {
		result = t.catch(result.GetErr(), catchArgs)
	}

	for k, v := range result.GetValues() {
		finallyArgs.SetArg(k, v, nil)
	}

	result = t.finally(finallyArgs)

	return result
}
