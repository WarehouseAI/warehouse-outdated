package models

import "errors"

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
)

var InternalError = errors.New("Internal Server Error.")
var ExistError = errors.New("Entity already exist.")
