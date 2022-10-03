# escargot
Go's try/catch/finally package

Simple Hello World Example

```go
package main

import (
	"fmt"
	"github.com/syke99/escargot/argument"
	"github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
	"github.com/syke99/escargot/try"
	"log"
)

// printHelloWorld is the function to be tried. The function signature
// must match func(args argument.Arguments) *shell.Shell
func printHelloWorld(args argument.Arguments) *shell.Shell {
	// create a "Shell" to hold your values and/or error
	res := shell.FreshShell()
	
	// get the argument with the value "hello"
	helloWorld, err := args.GetArg("hello")
	
	// make sure there is no error. If not, add it to the result (*shell.Shell);
	// nesting the error like this supports the *shell.Shell.Range() and
	// *shell.Shell.RangeWithCancel() methods
	if err != nil {
		res.Err(err, "")

		return res
    }

	// to use the value returned from args.getArg(), simply cast it to
	// the necessary type
	fmt.Println(helloWorld.(string))

	return nil
}

// errFunc is the function to be ran in case of error. The function signature
// must match func(e *err.EscargotError, args args argument.Arguments) *shell.Shell
func errFunc(e *error.EscargotError, args argument.Arguments) *shell.Shell {
	log.Fatal(e.Unwrap())
	
	return nil
}

func main() {
	// create your trier
	tr, err := trier.NewTrier(printHelloWorld, errFunc)
	if err != nil {
		log.Fatal(err.Error())
	}
	
	// tArgs are arguments to be used in the tryFunc (printHelloWorld in this case)
	tArgs := argument.NewArguments()
	
	// set an argument with the value "hello world" with the key "hello"; this
	// allows for more of a guarantee that casting the value to the necessary type
	// in the tryFunc will be successful; if you want to update/override a value at
	// the given key, you should pass in an argument.OverRide as the last
	// argument instead of nil
	tArgs.SetArg("hello", "hello world", nil)

	// cArgs are arguments to be used in the catchFunc (errFunc in this case)
	cArgs := argument.NewArguments()

	tr.Try(*tArgs, *cArgs)
}

```
