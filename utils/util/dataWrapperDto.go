package util

// DataWrapperDto is a generic struct to wrap data in a response, used for swagger documentation
type DataWrapperDto[T any] struct {
	Data T `json:"data"`
}

// MessageWrapperDto wraps a message in a response, used for swagger documentation
type MessageWrapperDto struct {
	Message string `json:"message"`
}

type ListDataWrapperDto[T any] struct {
	Data struct {
		List T `json:"list"`
	} `json:"data"`
}

type ListMessageDataWrapperDto[T any] struct {
	Data struct {
		Message string `json:"message"`
		List    T      `json:"list"`
	} `json:"data"`
}
