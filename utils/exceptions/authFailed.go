package exceptions

import "errors"

type AuthFailed struct {
	Message string
}

func (e *AuthFailed) Error() string {
	return "Authorization is failed"
}

func (e *AuthFailed) Is(err error) bool {
	var targetError *AuthFailed
	return errors.As(err, &targetError)
}
