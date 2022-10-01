package shell

import (
	"errors"
	"fmt"

	"github.com/syke99/escargot/error"
)

// Shell holds the value(s) and/or EscargotError produced whenever attempting to try the
// provided tryFunc
type Shell struct {
	values map[string]any
	err    *error.EscargotError
}

// GetErr returns the EscargotError created whenever attempting to try the provided tryFunc
func (s Shell) GetErr() *error.EscargotError {
	return s.err
}

// GetValue returns the value stored in the Shell with the given key as
// an interface (any). You can then use the value as you would any interface
// value (casting, switching, reflection, etc.)
func (s Shell) GetValue(key string) any {
	v, ok := s.values[key]
	if !ok {
		er := errors.New("attempt to access non-existent value")

		escErr := error.EscargotError{
			Level: "Error",
			Msg:   fmt.Sprintf("value with key %s does not exist", key),
		}

		escErr.Err(er)

		return escErr
	}

	return v
}

// SetValue sets the given value in the shell with the given key. To retrieve the value,
// use *Shell.GetValue(key string). To remove the value, use *Shell.RemoveValue(key string) *error.Escargot
func (s Shell) SetValue(key string, value any) {
	s.values[key] = value
}

func (s Shell) RemoveValue(key string) *error.EscargotError {
	_, ok := s.values[key]
	if !ok {
		er := errors.New("attempt to delete non-existent value")

		escErr := error.EscargotError{
			Level: "Error",
			Msg:   fmt.Sprintf("value with key %s does not exist", key),
		}

		escErr.Err(er)

		return &escErr
	}

	delete(s.values, key)

	return &error.EscargotError{}
}
