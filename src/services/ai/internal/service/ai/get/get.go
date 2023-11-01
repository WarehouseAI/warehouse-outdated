package get

import (
	"time"
	"warehouse/gen"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type Request struct {
	ID string `json:"id"`
}

type AiProvider interface {
	GetOneBy(map[string]interface{}) (*pg.AI, *db.DBError)
	GetOneByPreload(map[string]interface{}, string) (*pg.AI, *db.DBError)
}

func GetAiByID(getRequest Request, aiProvider AiProvider, logger *logrus.Logger) (*pg.AI, *httputils.ErrorResponse) {
	existAI, dbErr := aiProvider.GetOneBy(map[string]interface{}{"id": getRequest.ID})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}

func GetLoadedAiByID(request *gen.GetAiByIdMsg, aiProvider AiProvider, logger *logrus.Logger) (*pg.AI, *httputils.ErrorResponse) {
	existAI, dbErr := aiProvider.GetOneByPreload(map[string]interface{}{"id": request.Id}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existAI, nil
}
