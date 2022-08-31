package helper

import "strings"

type (
	ResponJSON[T any] struct {
		Status StatusStruct `json:"status"`
		Data   T            `json:"data"`
	}

	StatusStruct struct {
		Code    string `json:"code"`
		Message string `json:"Message"`
	}
)

type EmptyObj struct{}

func BuildSuccessResponse[T any](code, message string, data T) ResponJSON[T] {
	res := ResponJSON[T]{
		Status: StatusStruct{
			Code:    code,
			Message: message,
		},
		Data: data,
	}

	return res
}

func BuildErrorResponse(message string, err string) ResponJSON[[]string] {
	splittedError := strings.Split(err, "\n")
	res := ResponJSON[[]string]{
		Status: StatusStruct{
			Code:    "0000",
			Message: message,
		},
		Data: splittedError,
	}

	return res
}
