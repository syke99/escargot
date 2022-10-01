package try

import (
	"errors"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

type tryFunc[A any] func(args ...A) *shell.Shell[A]

type catchFunc[A any] func(err *err.EscargotError, args ...A)

// Trier will handle trying the TryFunc provided and execute the provided CatchFunc
// on error
type Trier[A any] struct {
	tryFunc   func(args ...A) *shell.Shell[A]
	catchFunc func(err *err.EscargotError, args ...A)
}

// NewTrier will return a new Trier with the provided TryFunc and CatchFunc
func NewTrier[A any](try tryFunc[A], catch catchFunc[A]) (Trier[A], error) {
	if try == nil ||
		catch == nil {
		return Trier[A]{}, errors.New("invalid Trier configuration")
	}

	return Trier[A]{
		tryFunc:   try,
		catchFunc: catch,
	}, nil
}

// Try tries the Trier's TryFunc with the provided tryArgs, and on error,
// will execute the Trier's CatchFunc with the provided catchArgs. It will
// return a *shell.Shell to access any values and/or errors
func (t Trier[A]) Try(tryArgs []A, catchArgs []A) *shell.Shell[A] {
	result := t.tryFunc(tryArgs...)

	if result.GetErrStatus() {
		t.catchFunc(result.GetErr(), catchArgs...)
	}

	return result
}
