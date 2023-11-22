package service

import (
	"time"
	"warehouseai/user/adapter"
	"warehouseai/user/dataservice"
	e "warehouseai/user/errors"
	m "warehouseai/user/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type AddFavoriteRequest struct {
	AiId string `json:"ai_id"`
}

type RemoveFavoriteRequest struct {
	AiId string `json:"ai_id"`
}

type GetFavoritesRequest struct {
	UserId string `json:"user_id"`
}

func AddFavorite(userId string, request *AddFavoriteRequest, favorites dataservice.FavoritesInterface, ai adapter.AiGrpcInterface, logger *logrus.Logger) *e.ErrorResponse {
	existAi, gwErr := ai.GetById(request.AiId)

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Add favorite")
		return gwErr
	}

	newFavorite := m.UserFavorites{
		AiId:   uuid.FromStringOrNil(existAi),
		UserId: uuid.FromStringOrNil(userId),
	}

	if err := favorites.Add(&newFavorite); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Add favorite")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}

func GetFavorites(request *GetFavoritesRequest, favorites dataservice.FavoritesInterface, logger *logrus.Logger) (*[]m.UserFavorites, *e.ErrorResponse) {
	userFavorites, err := favorites.GetUserFavorites(request.UserId)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Add favorite")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return userFavorites, nil
}

func GetFavorite(userId string, aiId string, favorites dataservice.FavoritesInterface, logger *logrus.Logger) (*m.UserFavorites, *e.ErrorResponse) {
	existFavorite, dbErr := favorites.GetFavorite(userId, aiId)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user favorite ai")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existFavorite, nil
}

func RemoveFavorite(userId string, request *RemoveFavoriteRequest, favorites dataservice.FavoritesInterface, logger *logrus.Logger) *e.ErrorResponse {
	if err := favorites.Delete(userId, request.AiId); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Get user by Id")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
