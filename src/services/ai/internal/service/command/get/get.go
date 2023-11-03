package get

import (
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Request struct {
	AiID uuid.UUID `json:"ai_id"`
	Name string    `json:"name"`
}

type Response struct {
	AI         *pg.AI
	Payload    pg.Command
	ApiKey     string
	AuthScheme pg.AuthScheme
}

type AiProvider interface {
	GetOneByPreload(map[string]interface{}, string) (*pg.AI, *db.DBError)
}

func GetCommand(getRequest Request, aiProvider AiProvider, logger *logrus.Logger) (*Response, *httputils.ErrorResponse) {
	existAI, dbErr := aiProvider.GetOneByPreload(map[string]interface{}{"id": getRequest.AiID}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get command")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	for i := 0; i <= len(existAI.Commands); i++ {
		if existAI.Commands[i].Name == getRequest.Name {
			return &Response{existAI, existAI.Commands[i], existAI.ApiKey, existAI.AuthScheme}, nil
		}
	}

	return nil, httputils.NewErrorResponse(httputils.NotFound, "Command not found.")
}
