package exceptions

import "errors"

type ResponseError struct {
	Url     string
	Message string
}

func (e *ResponseError) Error() string {
	return "Request to service is failed"
}

func (e *ResponseError) Is(err error) bool {
	var targetError *ResponseError
	return errors.As(err, &targetError)
}
