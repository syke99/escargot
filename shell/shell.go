package shell

import (
	"context"
	"fmt"
	"github.com/syke99/escargot/internal/resources"
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

func (s *Shell) GetValues() map[string]any {
	return s.values
}

func (s *Shell) buildErr(e resources.Err, m string) *err.EscargotError {
	escErr := err.EscargotError{
		Level: resources.ErrLevel,
		Msg:   m,
	}

	escErr.Err(e.Error())

	return &escErr
}

func (s *Shell) buildErrVal(e resources.Err, m string) *Shell {
	errVal.Err(s.buildErr(e, m))

	return &errVal
}

// GetValue returns the value stored in the Shell with the given key as
// an interface (any). You can then use the value as you would any interface
// value (casting, switching, reflection, etc.)
func (s *Shell) GetValue(key string) any {
	s.Lock()
	defer s.Unlock()

	if key == "" {
		return s.buildErr(resources.SetWithoutKey, resources.NoKeyProvided.String())
	}

	v, ok := s.values[key]

	if !ok {
		return s.buildErr(resources.AccessNonExistentValue, fmt.Sprintf(resources.NonExistentValue.String(), key))
	}

	return v
}

// OverRide is used to signal to *shell.Shell.SetValue() that a
// value should be allowed to be overriden
type OverRide *override.OverRider

// SetValue sets the given value in the shell with the given key. To retrieve the value,
// use *Shell.GetValue(key string). To remove the value, use *Shell.RemoveValue(key string) *error.Escargot
func (s *Shell) SetValue(key string, value any, override OverRide) any {
	s.Lock()
	defer s.Unlock()

	if key == "" {
		return s.buildErr(resources.NoKeyProvidedSet, resources.NoKeyProvided.String())
	}

	if override == nil {
		return s.buildErr(resources.OverRideWithoutOverRider, resources.OverRideNotAllowed.String())
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
		return s.buildErr(resources.DeleteNonExistentValue, resources.NonExistentValue.String())
	}

	delete(s.values, key)

	return nil
}

// CallBackX executes a callback.CallBackX function on the value returned from
// the given key. This method will error if a key is not provided or a value is not
// set at the given key
func (s *Shell) CallBackX(key string, cb callback.CallBackX, cbValOverRide OverRide) *Shell {

	if key == "" {
		return s.buildErrVal(resources.SetWithoutKey, resources.NoKeyProvided.String())
	}

	v, ok := s.values[key]

	if !ok {
		return s.buildErrVal(resources.AccessNonExistentValue, resources.NonExistentValue.String())
	}

	sh := Shell{
		values: make(map[string]any),
		err:    nil,
	}

	val, er := cb.CallBackX(key, v, cbValOverRide)

	for i, v := range val {
		sh.SetValue(fmt.Sprintf("cbValNum%d", i), v, cbValOverRide)
	}

	sh.Err(er)

	return &sh
}

// Range ranges over all the values in the shell and executes the given callback for each
// value
func (s *Shell) Range(cb callback.CallBackX, cbValOverRide OverRide) []*Shell {
	results := make([]*Shell, len(s.values))

	var wg sync.WaitGroup

	for k, v := range s.values {
		v := v

		key := k

		wg.Add(1)

		go func(key string, v any) {
			defer wg.Done()

			sh := Shell{
				values: make(map[string]any),
				err:    nil,
			}

			val, er := cb.CallBackX(k, v, cbValOverRide)

			for i, v := range val {
				sh.SetValue(fmt.Sprintf("cbValNum%d", i), v, cbValOverRide)
			}

			sh.Err(er)

			results = append(results, &sh)
		}(key, v)
	}

	wg.Wait()

	return results
}

// RangeWithCancel works just like Range, but takes a context
// to cancel execution
func (s *Shell) RangeWithCancel(ctx context.Context, cb callback.CallBackX, cbValOverRide OverRide) ([]*Shell, context.CancelFunc) {
	results := make([]*Shell, len(s.values))

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup

	for k, v := range s.values {
		v := v

		k := k

		wg.Add(1)

		go func(k string, v any) {
			defer wg.Done()

			sh := Shell{
				values: make(map[string]any),
				err:    nil,
			}

			val, er := cb.CallBackXWithCancellation(ctx, cancel, k, v, cbValOverRide)

			for i, v := range val {
				sh.SetValue(fmt.Sprintf("cbValNum%d", i), v, cbValOverRide)
			}

			sh.Err(er)

			results = append(results, &sh)
		}(k, v)
	}

	wg.Wait()

	return results, cancel
}
