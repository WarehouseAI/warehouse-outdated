package grpc

import e "warehouseai/user/errors"

type AuthGrpcInterface interface {
	Authenticate(sessionId string) (*string, *string, *e.ErrorResponse)
}
