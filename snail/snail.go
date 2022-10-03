package snail

import (
	"fmt"
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
	tShell *shell.Shell
	cShell *shell.Shell
}

// NewTrier will return a new Snail for chaining together a TryFunc, CatchFunc, and/or FinallyFunc
func NewSnail() *Snail {
	tShl := shell.Shell{}
	cShl := shell.Shell{}

	snl := Snail{
		tShell: &tShl,
		cShell: &cShl,
	}

	return &snl
}

// Try tries the provided TryFunc with the provided argument.Arguments
// and sets the value(s) for the Snail to use during Catch and/or Finally
func (t *Snail) Try(tryFunc TryFunc, tryArgs argument.Arguments) *Snail {
	t.tShell = tryFunc(tryArgs)

	return t
}

// Catch checks if and error occurred, and if so, executes the provided CatchFunc
// with the provided argument.Arguments, and well as any currently set value(s)
// in the Snail
func (t *Snail) Catch(catchFunc CatchFunc, catchArgs argument.Arguments) *Snail {
	if t.tShell.GetErrStatus() {
		for k, v := range t.tShell.GetValues() {

			k = fmt.Sprintf("try-%s", k)

			er := catchArgs.SetArg(k, v, nil)

			if er != nil {
				t.tShell.Err(er, "")

				t.tShell = t.cShell
			}
		}

		t.cShell = catchFunc(t.tShell.GetErr(), catchArgs)
	}

	return t
}

// Finally executes the provided FinallyFunc with the provided argument.Arguments
// along with any currently set value(s) in the Snail, regardless of error status
func (t *Snail) Finally(finallyFunc FinallyFunc, finallyArgs argument.Arguments) *Snail {
	for k, v := range t.tShell.GetValues() {
		er := finallyArgs.SetArg(k, v, nil)

		if er != nil {
			t.tShell.Err(er, er.Error())
		}
	}

	for k, v := range t.cShell.GetValues() {

		k = fmt.Sprintf("catch-%s", k)

		er := finallyArgs.SetArg(k, v, nil)

		if er != nil {
			t.tShell.Err(er, er.Error())
		}
	}

	t.tShell = finallyFunc(finallyArgs)

	return t
}

func (t *Snail) GetFinalResults() *shell.Shell {
	return t.tShell
}
