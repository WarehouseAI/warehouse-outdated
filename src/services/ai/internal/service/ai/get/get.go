package get

import (
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type AiProvider interface {
	GetOneBy(map[string]interface{}) (*pg.AI, *db.DBError)
	GetOneByPreload(map[string]interface{}, string) (*pg.AI, *db.DBError)
}

func GetAiByID(id string, aiProvider AiProvider, logger *logrus.Logger) (*pg.AI, *httputils.ErrorResponse) {
	existAI, dbErr := aiProvider.GetOneBy(map[string]interface{}{"id": id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}

func GetLoadedAiByID(id string, aiProvider AiProvider, logger *logrus.Logger) (*pg.AI, *httputils.ErrorResponse) {
	existAI, dbErr := aiProvider.GetOneByPreload(map[string]interface{}{"id": id}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}
