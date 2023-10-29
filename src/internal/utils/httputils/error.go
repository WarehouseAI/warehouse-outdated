package httputils

import (
	"fmt"
	db "warehouse/src/internal/database"
)

type ErrorResponseType string

const (
	ServerError ErrorResponseType = "Server Error: %s"
	Abort       ErrorResponseType = "Abort. %s"
)

type (
	ErrorResponse struct {
		ErrorType    ErrorResponseType `json:"error_type"`
		ErrorMessage string            `json:"error_message"`
	}
)

var errorMapper = map[db.DBErrorType]ErrorResponseType{
	db.Exist:    Abort,
	db.NotFound: Abort,
	db.System:   ServerError,
}

func NewErrorResponseFromDBError(errorType db.DBErrorType, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorType:    errorMapper[errorType],
		ErrorMessage: fmt.Sprintf(string(errorType), message),
	}
}

func NewErrorResponse(errorType ErrorResponseType, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorType:    errorType,
		ErrorMessage: fmt.Sprintf(string(errorType), message),
	}
}
