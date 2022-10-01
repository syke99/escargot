package argument

import (
	"errors"
	"sync"
)

// Arguments holds the arguments you wish to use in a callback.CallBack. It is
// used as an added safety measure whenever type asserting inside your custom
// callback functions through the use of *argument.Arguments.GetArg(key string).
// In conjunction with the use of *argument.Arguments.SetArg(key string, value any),
// this allows for a more strict guarantee that the value will correctly be
// asserted to the desired type at runtime
type Arguments struct {
	sync.RWMutex
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

// GetArgsLength returns the current number of arguments set
func (a *Arguments) GetArgsLength() int {
	return len(a.args)
}

// GetArgsSlice returns a slice of arguments currently set in the
// *argument.Arguments this method is called on
func (a *Arguments) GetArgsSlice() []any {
	args := make([]any, len(a.args))

	for _, v := range a.args {
		args = append(args, v)
	}

	return args
}

// GetArg returns the argument set with the given key
func (a *Arguments) GetArg(key string) (any, error) {
	a.Lock()
	defer a.Unlock()
	arg, ok := a.args[key]

	if key == "" {
		return nil, errors.New("no key provided to retrieve argument with")
	}

	if !ok {
		return nil, errors.New("argument does not exist in this set of arguments")
	}

	return &arg, nil
}

type overRider struct{}

// OverRide is used to signal to *argument.Arguments.SetArg() that an argument
// value should be allowed to be overriden
type OverRide *overRider

// SetArg checks for the existence of a provided key, as well as the existence of
// an OverRider in case of a pre-existing key. If a key already exists but no
// OverRider is provided, this method will error. If the key does not exist, the
// value will be added to the arguments
func (a *Arguments) SetArg(key string, value any, override OverRide) error {
	a.Lock()
	defer a.Unlock()

	_, ok := a.args[key]

	if key == "" {
		return errors.New("no key provided to set argument with")
	}

	if ok && override == nil {
		return errors.New("attempt to override argument value without a provided OverRider")
	}

	a.args[key] = &value

	return nil
}

// RemoveArg removes the argument from the *argument.Arguments this method is
// called on with the given key
func (a *Arguments) RemoveArg(key string) {
	delete(a.args, key)
}
