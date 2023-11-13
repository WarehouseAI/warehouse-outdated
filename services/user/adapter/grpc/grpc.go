package grpc

import e "warehouseai/internal/errors"

type AuthGrpcInterface interface {
	Authenticate(sessionId string) (*string, *e.ErrorResponse)
}
