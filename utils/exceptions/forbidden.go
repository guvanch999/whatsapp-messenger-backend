package exceptions

import "errors"

type Forbidden struct {
	Message string
}

func (e *Forbidden) Error() string {
	return "Route is forbidden"
}

func (e *Forbidden) Is(err error) bool {
	var targetError *Forbidden
	return errors.As(err, &targetError)
}
