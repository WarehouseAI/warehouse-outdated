package adapter

import e "warehouseai/user/errors"

type AuthGrpcInterface interface {
	Authenticate(sessionId string) (string, string, *e.ErrorResponse)
}

type AiGrpcInterface interface {
	GetById(aiId string) (string, *e.ErrorResponse)
}
