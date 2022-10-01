# escargot
Go's try catch package

Simple Hello World Example

```go
package main

import (
	"errors"
	"fmt"
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/shell"
	"github.com/syke99/escargot/trier"
	"log"
)

func printHelloWorld(args ...any) *shell.Shell {
	e := err.EscargotError{
		Level: "",
		Msg:   "",
	}

	defer func(er *err.EscargotError) *shell.Shell {
		if r := recover(); r != nil {
			e.Level = "Panic"
			e.Msg = "recovered during tryFunc"
			e.Err(errors.New("recovered from panic"))
		}

		res := shell.Shell{}

		res.Err(er)

		return &res
	}(&e)

	if len(args) == 0 {
		e.Level = "Fatal"
		e.Msg = "No arguments provided to TryFunc"

		er := errors.New("invalid call to printHelloWorld")

		e.Err(er)

		res := &shell.Shell{}

		res.Err(&e)

		return res
	}

	helloWorld := args[0].(string)

	fmt.Println(helloWorld)

	return &shell.Shell{}
}

func errFunc(e *err.EscargotError, args ...any) {
	if e.Unwrap() != nil {
		log.Fatal(e.Unwrap())
	}
}

func main() {
	tr, e := trier.NewTrier(printHelloWorld, errFunc)
	if e != nil {
		log.Fatal(e.Error())
	}

	tr.Try([]any{}, []any{})
}

```