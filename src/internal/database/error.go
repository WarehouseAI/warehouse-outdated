package database

type DBErrorType string

const (
	Exist    DBErrorType = "exist"
	NotFound DBErrorType = "not_found"
	Update   DBErrorType = "update"
	System   DBErrorType = "system"
)

type DBError struct {
	ErrorType DBErrorType
	Message   string
	Payload   string
}

func NewDBError(errorType DBErrorType, message string, payload string) *DBError {
	return &DBError{
		errorType,
		message,
		payload,
	}
}
