package exceptions

import "errors"

type AccessDenied struct {
}

func (e *AccessDenied) Error() string {
	return "You dont have an access for this route"
}

func (e *AccessDenied) Is(err error) bool {
	var targetError *AccessDenied
	return errors.As(err, &targetError)
}
