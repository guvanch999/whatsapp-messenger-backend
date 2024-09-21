package exceptions

import "errors"

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return "record not found"
}

func (e *NotFoundError) Is(err error) bool {
	var targetError *NotFoundError
	return errors.As(err, &targetError)
}
