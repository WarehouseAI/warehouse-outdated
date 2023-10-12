package get

import (
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/sirupsen/logrus"
)

type Request struct {
	ID string `json:"id"`
}

type AIProvider interface {
	GetOneBy(key string, value interface{}) (*pg.AI, error)
}

func GetByID(getRequest Request, aiProvider AIProvider, logger *logrus.Logger) (*pg.AI, error) {
	existAI, err := aiProvider.GetOneBy("id", getRequest.ID)

	if existAI == nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get AI")
		return nil, dto.NotFoundError
	}

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create AI")
		return nil, dto.InternalError
	}

	return existAI, nil
}
