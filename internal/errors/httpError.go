package errors

import (
	"github.com/gofiber/fiber/v2"
)

const (
	HttpInternalError int = fiber.StatusInternalServerError
	HttpAlreadyExist  int = fiber.StatusConflict
	HttpNotFound      int = fiber.StatusNotFound
	HttpBadRequest    int = fiber.StatusBadRequest
	HttpForbidden     int = fiber.StatusForbidden
	HttpUnauthorized  int = fiber.StatusUnauthorized
)

type (
	ErrorResponse struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"err_msg"`
	}
)

var dbErrorMapper = map[DBErrorType]int{
	DbExist:    HttpAlreadyExist,
	DbNotFound: HttpNotFound,
	DbUpdate:   HttpBadRequest,
	DbSystem:   HttpInternalError,
}

func NewErrorResponseFromDBError(errorType DBErrorType, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    dbErrorMapper[errorType],
		ErrorMessage: message,
	}
}

func NewErrorResponse(errorType int, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    errorType,
		ErrorMessage: message,
	}
}
