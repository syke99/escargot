package build

import (
	err "github.com/syke99/escargot/error"
	"github.com/syke99/escargot/internal/resources"
)

func BuildErr(e resources.Err, l, m string) *err.EscargotError {
	escErr := err.EscargotError{
		Level: l,
		Msg:   m,
	}

	escErr.Err(e.Error())

	return &escErr
}

func BuildCustomErr(e error, l, m string) *err.EscargotError {
	escErr := err.EscargotError{
		Level: l,
		Msg:   m,
	}

	escErr.Err(e)

	return &escErr
}
