package httputils

import (
	db "warehouse/src/internal/database"
	"warehouse/src/internal/s3"

	"github.com/gofiber/fiber/v2"
)

const (
	InternalError int = fiber.StatusInternalServerError
	AlreadyExist  int = fiber.StatusConflict
	NotFound      int = fiber.StatusNotFound
	BadRequest    int = fiber.StatusBadRequest
	Forbidden     int = fiber.StatusForbidden
	Unauthorized  int = fiber.StatusUnauthorized
)

type (
	ErrorResponse struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"err_msg"`
	}
)

var dbErrorMapper = map[db.DBErrorType]int{
	db.Exist:    AlreadyExist,
	db.NotFound: NotFound,
	db.Update:   BadRequest,
	db.System:   InternalError,
}

var s3ErrorMapper = map[s3.S3ErrorType]int{
	s3.System: InternalError,
}

func NewErrorResponseFromDBError(errorType db.DBErrorType, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    dbErrorMapper[errorType],
		ErrorMessage: message,
	}
}

func NewErrorResponseFromS3Error(errorType s3.S3ErrorType, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    s3ErrorMapper[errorType],
		ErrorMessage: message,
	}
}

func NewErrorResponse(errorType int, message string) *ErrorResponse {
	return &ErrorResponse{
		ErrorCode:    errorType,
		ErrorMessage: message,
	}
}
