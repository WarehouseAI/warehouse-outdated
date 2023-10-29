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

type CommandProvider interface {
	GetOneBy(key string, value interface{}) (*pg.Command, *db.DBError)
}

func GetCommand(getRequest Request, commandProvider CommandProvider, logger *logrus.Logger) (*pg.Command, *httputils.ErrorResponse) {
	// TODO: Получать команду по двум ключам, а не по одному
	existCommand, dbErr := commandProvider.GetOneBy("name", getRequest.Name)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get command")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existCommand, nil
}
