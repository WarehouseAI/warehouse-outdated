package ai

import (
	"time"
	"warehouseai/ai/adapter"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/sirupsen/logrus"
)

type GetAiResponse struct {
	m.AiProduct
	IsFavorite bool `json:"is_favorite"`
}

func GetById(id string, ai dataservice.AiInterface, logger *logrus.Logger) (*m.AiProduct, *e.HttpErrorResponse) {
	existAI, dbErr := ai.Get(map[string]interface{}{"id": id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}

func GetManyById(ids []string, ai dataservice.AiInterface, logger *logrus.Logger) (*[]m.AiProduct, *e.HttpErrorResponse) {
	existAis, dbErr := ai.GetMany(ids)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get many AIs")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAis, nil
}

func GetLike(field string, value string, ai dataservice.AiInterface, logger *logrus.Logger) (*[]m.AiProduct, *e.HttpErrorResponse) {
	if field == "auth_scheme" || field == "api_key" {
		return nil, e.NewErrorResponse(e.HttpBadRequest, "Invalid parameters")
	}

	if len(value) < 3 {
		return nil, e.NewErrorResponse(e.HttpBadRequest, "Too small value, provide value larger or equals 3.")
	}

	existAis, dbErr := ai.GetLike(field, "%"+value+"%")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI like")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAis, nil
}

func GetByIdPreload(id string, ai dataservice.AiInterface, logger *logrus.Logger) (*GetAiResponse, *e.HttpErrorResponse) {
	existAI, dbErr := ai.GetWithPreload(map[string]interface{}{"id": id}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI with Preload")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &GetAiResponse{*existAI, false}, nil
}

func GetByIdPreloadAuthed(userId string, aiId string, ai dataservice.AiInterface, user adapter.UserGrpcInterface, logger *logrus.Logger) (*GetAiResponse, *e.HttpErrorResponse) {
	existAI, dbErr := ai.GetWithPreload(map[string]interface{}{"id": aiId}, "Commands")
	isAiFavorite, gwErr := user.GetFavorite(aiId, userId)

	if gwErr != nil {
		return nil, gwErr
	}

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI with Preload")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &GetAiResponse{*existAI, isAiFavorite}, nil
}
