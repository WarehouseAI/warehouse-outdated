package adapter

import e "warehouseai/ai/errors"

type AuthGrpcInterface interface {
	Authenticate(sessionId string) (string, string, *e.HttpErrorResponse)
}

type UserGrpcInterface interface {
	GetFavorite(aiId string, userId string) (bool, *e.HttpErrorResponse)
}
