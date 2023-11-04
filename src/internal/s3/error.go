package s3

type S3ErrorType string

const (
	System S3ErrorType = "system"
)

type S3Error struct {
	ErrorType S3ErrorType
	Message   string
	Payload   string
}

func NewS3Error(errorType S3ErrorType, message string, payload string) *S3Error {
	return &S3Error{
		errorType,
		message,
		payload,
	}
}
