### Escargot

[![Go Reference](https://pkg.go.dev/badge/github.com/syke99/escargot.svg)](https://pkg.go.dev/github.com/syke99/escargot)
[![go reportcard](https://goreportcard.com/badge/github.com/syke99/escargot)](https://goreportcard.com/report/github.com/syke99/escargot)
![Go version](https://img.shields.io/github/go-mod/go-version/syke99/escargot)</br>
# escargot
A simple package for cutting down on code bloat and boilerplate code whenever implementing try/catch/finally

How do I use Escargot?
====

### Installation

```
go get github.com/syke99/escargot
```

### Simple Hello World Example

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/syke99/escargot/argument"
	"github.com/syke99/escargot/error"
	"github.com/syke99/escargot/internal/resources"
	"github.com/syke99/escargot/shell"
	"github.com/syke99/escargot/snail"
)
```

If you want to have functions that can be re-used, you can define them at whatever scope you'd like to be
able to re-use them where needed/applilcable. You aren't limited to top level functions, either. As long as
the signature matches a snail.TryFunc, snail.CatchFunc, or snail.FinallyFunc signature, you can use methods
defined on structs for the respective step of a try/catch/finally block with your Snail

```go
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

	return shell.FreshShell()
}

// errFunc is the function to be ran in case of error. The function signature
// must match func(e *err.EscargotError, args args argument.Arguments) *shell.Shell
func errFunc(e *error.EscargotError, args argument.Arguments) *shell.Shell {
	log.Fatal(e.Unwrap())

	return shell.FreshShell()
}

// finalFunc is the function to be ran in the Finally portion of a 
// Try/Catch/Finally block of code
func finalFunc(args argument.Arguments) *shell.Shell {
	fmt.Println("Try/Catch/Finally block complete")
	
	return shell.FreshShell()
}
```

As an entry point, create a Snail to execute your try/catch/finally block(s) with

```go
func main() {
	snl := snail.NewSnail()
```

After creating your Snail, create any arguments needed for your functions to be executed.

```go
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
```

Finally, execute the chain of Try/Catch/(and if desired)Finally
by providing a TryFunc, CatchFunc, and FinallyFunc, respectively,
along with the appropriate arguments. 
```go
	results := snl.Try(printHelloWorld, *tArgs).
		        Catch(errFunc, *cArgs).
		        Finally(finalFunc, nil).
		        GetFinalResults()
	
	println(results.GetValues())
```

You can also pass the functions
in as anonymous functions without having to predefine them to match
a more similar style to executing try/catch/finally blocks in other languages.
This allows for the ability to reuse TryFuncs/CatchFuncs/FinallyFuncs accross
your program if desired and appropriately applicable to your use-case

```go
    results = snl.Try(func(args argument.Arguments) *shell.Shell {
                    res := shell.FreshShell()
            
                    helloWorld, err := args.GetArg("hello")
            
                    if err != nil {
                        res.Err(err, "")
            
                        return res
                    }
            
                    fmt.Println(helloWorld.(string))
            
                    return res
                }, *tArgs).
		Catch(func(e *error.EscargotError, args argument.Arguments) *shell.Shell {
			log.Print(fmt.Sprintf("log: Level %s Error: %v Message: %s", e.Level, e.Unwrap(), e.Msg))

			return shell.FreshShell()
		}, *cArgs).
		Finally(func(args argument.Arguments) *shell.Shell {
			fmt.Println("Try/Catch/Finally block complete")

			return shell.FreshShell()
			}, nil).
		GetFinalResults()

	println(results.GetValues())

}
```

### Complete Example

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/syke99/escargot/argument"
	"github.com/syke99/escargot/error"
	"github.com/syke99/escargot/internal/resources"
	"github.com/syke99/escargot/shell"
	"github.com/syke99/escargot/snail"
)

func printHelloWorld(args argument.Arguments) *shell.Shell {
	res := shell.FreshShell()

	helloWorld, err := args.GetArg("hello")

	if err != nil {
		res.Err(err, "")

		return res
	}

	fmt.Println(helloWorld.(string))

	return shell.FreshShell()
}

func errFunc(e *error.EscargotError, args argument.Arguments) *shell.Shell {
	log.Fatal(e.Unwrap())

	return shell.FreshShell()
}

func finalFunc(args argument.Arguments) *shell.Shell {
	fmt.Println("Try/Catch/Finally block complete")
	
	return shell.FreshShell()
}

func main() {
	snl := snail.NewSnail()
	
	tArgs := argument.NewArguments()

	tArgs.SetArg("hello", "hello world", nil)
	
	cArgs := argument.NewArguments()
	
	results := snl.Try(printHelloWorld, *tArgs).
		        Catch(errFunc, *cArgs).
		        Finally(finalFunc, nil).
		        GetFinalResults()
	
	println(results.GetValues())
	
	// results = snl.Try(func(args argument.Arguments) *shell.Shell {
        //             res := shell.FreshShell()
        //     
        //             helloWorld, err := args.GetArg("hello")
        //     
        //             if err != nil {
        //                 res.Err(err, "")
        //     
        //                 return res
        //             }
        //     
        //             fmt.Println(helloWorld.(string))
        //     
        //             return res
        //         }, *tArgs).
	// 	Catch(func(e *error.EscargotError, args argument.Arguments) *shell.Shell {
	// 		log.Print(fmt.Sprintf("log: Level %s Error: %v Message: %s", e.Level, e.Unwrap(), e.Msg))
	// 
	// 		return shell.FreshShell()
	// 	}, *cArgs).
	// 	Finally(func(args argument.Arguments) *shell.Shell {
	// 		fmt.Println("Try/Catch/Finally block complete")
	// 
	// 		return shell.FreshShell()
	// 		}, nil).
	// 	GetFinalResults()
	// 
	// println(results.GetValues())
}
```
