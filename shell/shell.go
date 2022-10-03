package shell

import (
	"fmt"
	"sync"

	"github.com/syke99/escargot/internal/build"
	"github.com/syke99/escargot/internal/resources"

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

func FreshShell() *Shell {
	shell := Shell{}

	return &shell
}

// Err sets an err
func (s *Shell) Err(er error, message string) {
	if message == "" {
		message = er.Error()
	}

	e := build.BuildCustomErr(er, resources.ErrLevel, message)

	s.err = e
}

// GetErrStatus returns the status of whether err is set
func (s *Shell) GetErrStatus() bool {
	return s.err != nil
}

// GetErr returns the EscargotError created whenever attempting to try the provided tryFunc
func (s *Shell) GetErr() *err.EscargotError {
	return s.err
}

func (s *Shell) GetValues() map[string]any {
	return s.values
}

func (s *Shell) buildErrVal(e resources.Err, l, m string) *Shell {
	errVal.err = build.BuildErr(e, l, m)

	return &errVal
}

// GetValue returns the value stored in the Shell with the given key as
// an interface (any). You can then use the value as you would any interface
// value (casting, switching, reflection, etc.)
func (s *Shell) GetValue(key string) (any, *err.EscargotError) {
	s.Lock()
	defer s.Unlock()

	if key == "" {
		return nil, build.BuildErr(resources.SetWithoutKey, resources.ErrLevel, resources.NoKeyProvided.String())
	}

	v, ok := s.values[key]

	if !ok {
		return nil, build.BuildErr(resources.AccessNonExistentValue, resources.ErrLevel, fmt.Sprintf(resources.NonExistentValue.String(), key))
	}

	return v, nil
}

// OverRide is used to signal to *shell.Shell.SetValue() that a
// value should be allowed to be overriden
type OverRide *override.OverRider

// SetValue sets the given value in the shell with the given key. To retrieve the value,
// use *Shell.GetValue(key string). To remove the value, use *Shell.RemoveValue(key string) *error.Escargot
func (s *Shell) SetValue(key string, value any, override OverRide) *err.EscargotError {
	s.Lock()
	defer s.Unlock()

	if key == "" {
		return build.BuildErr(resources.NoKeyProvidedSet, resources.ErrLevel, resources.NoKeyProvided.String())
	}

	if override == nil {
		return build.BuildErr(resources.OverRideWithoutOverRider, resources.ErrLevel, resources.OverRideNotAllowed.String())
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
		return build.BuildErr(resources.DeleteNonExistentValue, resources.ErrLevel, resources.NonExistentValue.String())
	}

	delete(s.values, key)

	vals := make(map[string]any)

	for k, v := range s.values {
		vals[k] = v
	}

	s.values = vals

	return nil
}

// CallBackX executes a callback.CallBackX function on the value returned from
// the given key. This method will error if a key is not provided or a value is not
// set at the given key
func (s *Shell) CallBackX(key string, cb callback.CallBackX, cbValOverRide OverRide) *Shell {

	if key == "" {
		return s.buildErrVal(resources.SetWithoutKey, resources.ErrLevel, resources.NoKeyProvided.String())
	}

	v, ok := s.values[key]

	if !ok {
		return s.buildErrVal(resources.AccessNonExistentValue, resources.ErrLevel, resources.NonExistentValue.String())
	}

	sh := Shell{
		values: make(map[string]any),
		err:    nil,
	}

	val, er := cb.CallBackX(key, v, cbValOverRide)

	for _, v := range val {
		sh.SetValue(key, v, cbValOverRide)
	}

	sh.err = er

	return &sh
}

// Range ranges over all the values in the shell and executes the given callback for each
// value
func (s *Shell) Range(cb callback.CallBackX, cbValOverRide OverRide) {
	var wg sync.WaitGroup

	for k, v := range s.values {
		v := v

		key := k

		wg.Add(1)

		go func(key string, v any) {
			defer wg.Done()
			s.Lock()
			defer s.Unlock()

			sh := Shell{
				values: make(map[string]any),
				err:    nil,
			}

			val, er := cb.CallBackX(key, v, cbValOverRide)

			for i, v := range val {
				sh.SetValue(fmt.Sprintf("cbValNum%d", i), v, cbValOverRide)
			}

			sh.err = er

			s.values[key] = sh
		}(key, v)
	}

	wg.Wait()
}

// RangeAtKeys works just like Range, but only executes the callback if
// a Shell value's key exists in the slice of keys provided
func (s *Shell) RangeAtKeys(cb callback.CallBackX, cbValOverRide OverRide, keys []string) {
	var wg sync.WaitGroup

	km := make(map[string]struct{}, len(keys))

	for _, key := range keys {
		km[key] = struct{}{}
	}

	for k, v := range s.values {
		v := v

		key := k

		if _, ok := km[k]; !ok {
			continue
		}

		wg.Add(1)

		go func(key string, v any) {
			defer wg.Done()
			s.Lock()
			defer s.Unlock()

			sh := Shell{
				values: make(map[string]any),
				err:    nil,
			}

			val, er := cb.CallBackX(key, v, cbValOverRide)

			for i, v := range val {
				sh.SetValue(fmt.Sprintf("cbValNum%d", i), v, cbValOverRide)
			}

			sh.err = er

			s.values[key] = sh
		}(key, v)
	}

	wg.Wait()
}
