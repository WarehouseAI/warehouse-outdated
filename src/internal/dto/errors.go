package dto

import "errors"

type (
	ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

var InternalError = errors.New("Internal Server Error.")
var ExistError = errors.New("Entity already exist.")
var NotFoundError = errors.New("Entity not found.")
var BadRequestError = errors.New("Invalid or empty request body.")
