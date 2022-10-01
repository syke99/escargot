package shell

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/syke99/escargot/internal/override"

	"github.com/syke99/escargot/callback"
	err "github.com/syke99/escargot/error"
)

// Shell holds the value(s) and/or EscargotError produced whenever attempting to try the
// provided tryFunc
type Shell struct {
	sync.RWMutex
	values map[string]any
	err    *err.EscargotError
}

var errVal = Shell{
	values: make(map[string]any),
	err:    nil,
}

// Err sets an err
func (s *Shell) Err(err *err.EscargotError) {
	s.err = err
}

// GetErrStatus returns the status of whether err is set
func (s *Shell) GetErrStatus() bool {
	if s.err != nil {
		return true
	}
	return false
}

// GetErr returns the EscargotError created whenever attempting to try the provided tryFunc
func (s *Shell) GetErr() *err.EscargotError {
	return s.err
}

// GetValue returns the value stored in the Shell with the given key as
// an interface (any). You can then use the value as you would any interface
// value (casting, switching, reflection, etc.)
func (s *Shell) GetValue(key string) any {
	s.Lock()
	defer s.Unlock()

	if key == "" {
		er := errors.New("attempt to set value with non-existent key")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "no key provided",
		}

		escErr.Err(er)

		return &escErr
	}

	v, ok := s.values[key]

	if !ok {
		er := errors.New("attempt to access non-existent value")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   fmt.Sprintf("value with key %s does not exist", key),
		}

		escErr.Err(er)

		return &escErr
	}

	return v
}

// OverRide is used to signal to *shell.Shell.SetValue() that a
// value should be allowed to be overriden
type OverRide *override.OverRider

// SetValue sets the given value in the shell with the given key. To retrieve the value,
// use *Shell.GetValue(key string). To remove the value, use *Shell.RemoveValue(key string) *error.Escargot
func (s *Shell) SetValue(key string, value any) any {
	s.Lock()
	defer s.Unlock()

	if key == "" {
		er := errors.New("no key provided to set value with")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "no key provided",
		}

		escErr.Err(er)

		return escErr
	}

	s.values[key] = value

	return nil
}

// RemoveValue removes the value from the shell with the given key if it exists
func (s *Shell) RemoveValue(key string) *err.EscargotError {
	s.Lock()
	defer s.Unlock()

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

// CallBack executes a callback.CallBackX function on the value returned from
// the given key. This method will error if a key is not provided or a value is not
// set at the given key
func (s *Shell) CallBack(key string, cb callback.CallBackX) *Shell {

	if key == "" {
		er := errors.New("attempt to set value with non-existent key")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   "no key provided",
		}

		escErr.Err(er)

		errVal.Err(&escErr)

		return &errVal
	}

	v, ok := s.values[key]

	if !ok {
		er := errors.New("attempt to access non-existent value")

		escErr := err.EscargotError{
			Level: "Error",
			Msg:   fmt.Sprintf("value with key %s does not exist", key),
		}

		escErr.Err(er)

		errVal.Err(&escErr)

		return &errVal
	}

	sh := Shell{
		values: make(map[string]any),
		err:    nil,
	}

	val, er := cb.CallBackX(v)

	sh.SetValue("value", val)
	sh.Err(er)

	return &sh
}

// Range ranges over all the values in the shell and executes the given callback for each
// value
func (s *Shell) Range(cb callback.CallBackX) []*Shell {
	results := make([]*Shell, len(s.values))

	var wg sync.WaitGroup

	for _, v := range s.values {
		v := v

		wg.Add(1)

		go func(v any) {
			defer wg.Done()

			sh := Shell{
				values: make(map[string]any),
				err:    nil,
			}

			val, er := cb.CallBackX(v)

			sh.SetValue("value", val)
			sh.Err(er)

			results = append(results, &sh)
		}(v)
	}

	wg.Wait()

	return results
}

// RangeWithCancel works just like Range, but takes a context
// to cancel execution
func (s *Shell) RangeWithCancel(ctx context.Context, cb callback.CallBackX) ([]*Shell, context.CancelFunc) {
	results := make([]*Shell, len(s.values))

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	for _, v := range s.values {
		v := v

		wg.Add(1)

		go func(v any) {
			defer wg.Done()

			sh := Shell{
				values: make(map[string]any),
				err:    nil,
			}

			val, er := cb.CallBackXWithCancellation(ctx, cancel, v)

			sh.SetValue("value", val)
			sh.Err(er)

			results = append(results, &sh)
		}(v)
	}

	wg.Wait()

	return results, cancel
}
