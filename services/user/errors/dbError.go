package errors

type DBErrorType string

const (
	DbExist    DBErrorType = "exist"
	DbNotFound DBErrorType = "not_found"
	DbUpdate   DBErrorType = "update"
	DbSystem   DBErrorType = "system"
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
