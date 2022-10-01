package try

import (
	"errors"
	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

type tryFunc func(args ...argument.Arguments) *shell.Shell

type catchFunc func(err *err.EscargotError, args ...argument.Arguments)

// Trier will handle trying the TryFunc provided and execute the provided CatchFunc
// on error
type Trier struct {
	tryFunc   func(args ...argument.Arguments) *shell.Shell
	catchFunc func(err *err.EscargotError, args ...argument.Arguments)
}

// NewTrier will return a new Trier with the provided TryFunc and CatchFunc
func NewTrier(try tryFunc, catch catchFunc) (Trier, error) {
	if try == nil ||
		catch == nil {
		return Trier{}, errors.New("invalid Trier configuration")
	}

	return Trier{
		tryFunc:   try,
		catchFunc: catch,
	}, nil
}

// Try tries the Trier's TryFunc with the provided tryArgs, and on error,
// will execute the Trier's CatchFunc with the provided catchArgs. It will
// return a *shell.Shell to access any values and/or errors
func (t Trier) Try(tryArgs []argument.Arguments, catchArgs []argument.Arguments) *shell.Shell {
	result := t.tryFunc(tryArgs...)

	if result.GetErrStatus() {
		t.catchFunc(result.GetErr(), catchArgs...)
	}

	return result
}
