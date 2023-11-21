package ai

import (
	"time"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/sirupsen/logrus"
)

func GetById(id string, ai dataservice.AiInterface, logger *logrus.Logger) (*m.AI, *e.ErrorResponse) {
	existAI, dbErr := ai.Get(map[string]interface{}{"id": id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}

func GetManyById(ids []string, ai dataservice.AiInterface, logger *logrus.Logger) (*[]m.AI, *e.ErrorResponse) {
	existAis, dbErr := ai.GetMany(ids)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AIs")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAis, nil
}

func GetByIdPreload(id string, ai dataservice.AiInterface, logger *logrus.Logger) (*m.AI, *e.ErrorResponse) {
	existAI, dbErr := ai.GetWithPreload(map[string]interface{}{"id": id}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}
