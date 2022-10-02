package argument

import (
	"errors"
	"sync"

	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/internal/override"
)

// Arguments holds the arguments you wish to use in a callback.CallBack. It is
// used as an added safety measure whenever type asserting inside your custom
// callback functions through the use of *argument.Arguments.GetArg(key string).
// In conjunction with the use of *argument.Arguments.SetArg(key string, value any),
// this allows for a more strict guarantee that the value will correctly be
// asserted to the desired type at runtime
type Arguments struct {
	*sync.RWMutex
	args map[string]*any
}

// NewArguments returns a pointer to a new Arguments struct so that
// arguments can be added and removed across your application
func NewArguments(key string, value any) *Arguments {
	args := Arguments{}

	argMap := make(map[string]*any)

	args.args = argMap

	return &args
}

// GetArgsSlice returns a slice of arguments currently set in the
// *argument.Arguments this method is called on
func (a Arguments) GetArgsSlice() []any {
	args := make([]any, len(a.args))

	for _, v := range a.args {
		args = append(args, v)
	}

	return args
}

// GetArg returns the argument set with the given key
func (a *Arguments) GetArg(key string) (any, err.EscargotError) {
	a.Lock()
	defer a.Unlock()
	arg, ok := a.args[key]

	if key == "" {
		er := errors.New("no key provided to retrieve argument with")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "no key provided",
		}

		escErr.Err(er)

		return nil, escErr
	}

	if !ok {
		er := errors.New("argument does not exist in this set of arguments")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "nox-existent argument in arguments",
		}

		escErr.Err(er)

		return nil, escErr
	}

	return &arg, err.EscargotError{}
}

// OverRide is used to signal to *argument.Arguments.SetArg() that an argument
// value should be allowed to be overriden
type OverRide *override.OverRider

// SetArg checks for the existence of a provided key, as well as the existence of
// an OverRider in case of a pre-existing key. If a key already exists but no
// OverRider is provided, this method will error. If the key does not exist, the
// value will be added to the arguments
func (a *Arguments) SetArg(key string, value any, override OverRide) err.EscargotError {
	a.Lock()
	defer a.Unlock()

	if key == "" {
		er := errors.New("no key provided to set argument with")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "no key provided",
		}

		escErr.Err(er)

		return escErr
	}

	_, ok := a.args[key]

	if ok && override == nil {
		er := errors.New("attempt to override argument value without a provided OverRider")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "override not explicitly allowed",
		}

		escErr.Err(er)

		return escErr
	}

	a.args[key] = &value

	return err.EscargotError{}
}

// RemoveArg removes the argument from the *argument.Arguments this method is
// called on with the given key
func (a *Arguments) RemoveArg(key string) {
	delete(a.args, key)
}
