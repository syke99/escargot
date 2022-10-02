package try

import (
	"errors"
	"github.com/syke99/escargot/argument"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
)

type TryFunc func(args ...any) *shell.Shell

type CatchFunc func(err *err.EscargotError, args ...any) *shell.Shell

type FinallyFunc func(args ...any) *shell.Shell

// Trier will handle trying the TryFunc provided and execute the provided CatchFunc
// on error
type Trier struct {
	try     func(args ...any) *shell.Shell
	catch   func(err *err.EscargotError, args ...any) *shell.Shell
	finally func(args ...any) *shell.Shell
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
	targs := make([]any, len(tryArgs.GetArgsSlice()))

	for _, arg := range tryArgs.GetArgsSlice() {
		targs = append(targs, arg)
	}

	result := t.try(targs...)

	if result.GetErrStatus() {
		cargs := make([]any, len(catchArgs.GetArgsSlice()))

		for _, arg := range catchArgs.GetArgsSlice() {
			cargs = append(cargs, arg)
		}

		result = t.catch(result.GetErr(), cargs...)
	}

	return result
}

// TryFinally works just like Try, but executes a FinallyFunc after the TryFunc
// and/or CatchFunc, regardless of the outcome of either function
func (t Trier) TryFinally(tryArgs argument.Arguments, catchArgs argument.Arguments, finallyArgs argument.Arguments) *shell.Shell {
	targs := make([]any, len(tryArgs.GetArgsSlice()))

	for _, arg := range tryArgs.GetArgsSlice() {
		targs = append(targs, arg)
	}

	result := t.try(targs...)

	if result.GetErrStatus() {
		cargs := make([]any, len(catchArgs.GetArgsSlice()))

		for _, arg := range catchArgs.GetArgsSlice() {
			cargs = append(cargs, arg)
		}

		result = t.catch(result.GetErr(), cargs...)
	}

	fargs := make([]any, len(finallyArgs.GetArgsSlice())+len(result.GetValues()))

	for _, arg := range finallyArgs.GetArgsSlice() {
		fargs = append(fargs, arg)
	}

	rargs := argument.NewArguments("", "")

	for k, v := range result.GetValues() {
		rargs.SetArg(k, v, nil)
	}

	for _, arg := range rargs.GetArgsSlice() {
		fargs = append(fargs, arg)
	}

	result = t.finally(fargs...)

	return result
}
