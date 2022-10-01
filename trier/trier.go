package trier

import (
	"errors"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

// TryFunc is the function you wish to attempt and provided arguments. It will
// return a *shell.Shell that will allow access to both values and errors
type TryFunc func(args ...any) *shell.Shell

// CatchFunc upon an *error.EscargotError returned inside the *shell.Shell returned
// from the provided TryFunc, CatchFunc will execute with any provided arguments.
type CatchFunc func(err *err.EscargotError, args ...any)

// Trier will handle trying the TryFunc provided and execute the provided CatchFunc
// on error
type Trier struct {
	tryFunc   func(args ...any) *shell.Shell
	catchFunc func(err *err.EscargotError, args ...any)
}

// NewTrier will return a new Trier with the provided TryFunc and CatchFunc
func NewTrier(try TryFunc, catch CatchFunc) (Trier, error) {
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
func (t Trier) Try(tryArgs []any, catchArgs []any) *shell.Shell {
	result := t.tryFunc(tryArgs...)

	if result.GetErr().Unwrap() != nil {
		t.catchFunc(result.GetErr(), catchArgs...)
	}

	return result
}
