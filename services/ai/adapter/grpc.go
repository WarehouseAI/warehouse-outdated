package adapter

import e "warehouseai/ai/errors"

type AuthGrpcInterface interface {
	Authenticate(sessionId string) (string, string, *e.ErrorResponse)
}

type UserGrpcInterface interface {
	GetFavorite(aiId string, userId string) (bool, *e.ErrorResponse)
}
