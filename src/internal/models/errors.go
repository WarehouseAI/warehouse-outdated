package models

import "errors"

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
)

var InternalError = errors.New("Internal")
var ExistError = errors.New("Exist")
