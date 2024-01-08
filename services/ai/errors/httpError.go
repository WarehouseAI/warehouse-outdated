package errors

import (
	"github.com/gofiber/fiber/v2"
)

const (
	HttpInternalError       int = fiber.StatusInternalServerError
	HttpAlreadyExist        int = fiber.StatusConflict
	HttpNotFound            int = fiber.StatusNotFound
	HttpBadRequest          int = fiber.StatusBadRequest
	HttpForbidden           int = fiber.StatusForbidden
	HttpUnauthorized        int = fiber.StatusUnauthorized
	HttpTimeout             int = fiber.StatusGatewayTimeout
	HttpUnprocessableEntity int = fiber.StatusUnprocessableEntity
)

type (
	HttpErrorResponse struct {
		ErrorCode    int      `json:"error_code"`
		ErrorMessage []string `json:"err_msg"`
	}
)

var dbErrorMapper = map[DBErrorType]int{
	DbExist:    HttpAlreadyExist,
	DbNotFound: HttpNotFound,
	DbUpdate:   HttpBadRequest,
	DbSystem:   HttpInternalError,
}

func NewErrorResponseFromDBError(errorType DBErrorType, message string) *HttpErrorResponse {
	return &HttpErrorResponse{
		ErrorCode:    dbErrorMapper[errorType],
		ErrorMessage: []string{message},
	}
}

func NewErrorResponse(errorType int, message string) *HttpErrorResponse {
	return &HttpErrorResponse{
		ErrorCode:    errorType,
		ErrorMessage: []string{message},
	}
}

func NewErrorResponseMultiple(errorType int, messages []string) *HttpErrorResponse {
	return &HttpErrorResponse{
		ErrorCode:    errorType,
		ErrorMessage: messages,
	}
}
