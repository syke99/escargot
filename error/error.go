package error

// EscargotError holds the error level, message, and error. EscargotError can also
// wrap other EscargotErrors with the use of *EscargotError.Wrap(err *EscargotError)
type EscargotError struct {
	Level      string
	Msg        string
	err        error
	wrappedErr *EscargotError
}

// Error returns the string representation of the stored error value
func (e EscargotError) Error() string {
	return e.err.Error()
}

// Unwrap unwraps the error stored in the EscargotError
func (e EscargotError) Unwrap() error {
	return e.err
}

// Err returns the error value stored inside the EscargotError this method is called on
func (e EscargotError) Err(err error) {
	e.err = err
}

// UnwrapError unwraps the wrapped EscargotError if one is wrapped, otherwise nil is returned
func (e EscargotError) UnwrapError() *EscargotError {
	if e.wrappedErr == nil {
		return nil
	}
	return e.wrappedErr
}

// Wrap wraps the provided EscargotError inside the EscargotError this method is called on
func (e EscargotError) Wrap(err *EscargotError) {
	e.wrappedErr = err
}
