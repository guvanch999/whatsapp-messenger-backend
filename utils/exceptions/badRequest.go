package exceptions

import (
	"errors"
	"fmt"
)

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Some params is invalid: %s", e.Message)
}

func (e *BadRequestError) Is(err error) bool {
	var targetError *BadRequestError
	return errors.As(err, &targetError)
}
