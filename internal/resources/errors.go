package resources

import "errors"

const ErrLevel = "Error"
const CancelLevel = "Cancel"

type Err int

const (
	SetWithoutKey Err = iota
	AccessNonExistentValue
	NoKeyProvidedSet
	DeleteNonExistentValue
	OverRideWithoutOverRider
	ContextCancel
)

func (e Err) Error() error {
	return [...]error{
		errors.New("attempt to set value with non-existent key"),
		errors.New("attempt to access non-existent value"),
		errors.New("no key provided to set value with"),
		errors.New("attempt to delete non-existent value"),
		errors.New("attempt to override value without explicit OverRide provided"),
		errors.New("context cancel signal received"),
	}[e]
}

type ErrMsg int

const (
	NoKeyProvided ErrMsg = iota
	NonExistentValue
	OverRideNotAllowed
	CanceledContext
)

func (em ErrMsg) String() string {
	return [...]string{
		"no key provided",
		"value with key %s does not exist",
		"attempt to override an already existing key/value pair with nil OverRide",
		"context canceled",
	}[em]
}
