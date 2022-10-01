package shell

import (
	"errors"
	"fmt"
	"sync"

	"github.com/syke99/escargot/callback"
	err "github.com/syke99/escargot/error"
)

// Shell holds the value(s) and/or EscargotError produced whenever attempting to try the
// provided tryFunc
type Shell struct {
	values map[string]any
	err    *err.EscargotError
}

// Err sets an err
func (s *Shell) Err(err *err.EscargotError) {
	s.err = err
}

// GetErrStatus returns the status of whether err is set
func (s *Shell) GetErrStatus() bool {
	if s.err == nil {
		return false
	}
	return true
}

// GetErr returns the EscargotError created whenever attempting to try the provided tryFunc
func (s *Shell) GetErr() *err.EscargotError {
	return s.err
}

// GetValue returns the value stored in the Shell with the given key as
// an interface (any). You can then use the value as you would any interface
// value (casting, switching, reflection, etc.)
func (s *Shell) GetValue(key string) any {
	v, ok := s.values[key]
	if !ok {
		er := errors.New("attempt to access non-existent value")

		escErr := err.EscargotError{
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
func (s *Shell) SetValue(key string, value any) {
	s.values[key] = value
}

// RemoveValue removes the value from the shell with the given key if it exists
func (s *Shell) RemoveValue(key string) *err.EscargotError {
	_, ok := s.values[key]
	if !ok {
		er := errors.New("attempt to delete non-existent value")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   fmt.Sprintf("value with key %s does not exist", key),
		}

		escErr.Err(er)

		return &escErr
	}

	delete(s.values, key)

	return &err.EscargotError{}
}

// Range ranges over all the values in the shell and executes the given callback for each
// value
func (s *Shell) Range(cb callback.CallBack) []*Shell {
	results := []*Shell{}

	var wg sync.WaitGroup

	for _, v := range s.values {
		v := v

		wg.Add(1)

		go func(v any) {
			defer wg.Done()

			results = append(results, cb.CallBack(v))
		}(v)
	}

	wg.Wait()

	return results
}
