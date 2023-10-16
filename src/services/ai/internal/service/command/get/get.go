package get

import (
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Request struct {
	AiID uuid.UUID `json:"ai_id"`
	Name string    `json:"name"`
}

type CommandProvider interface {
	GetOneBy(key string, value interface{}) (*pg.Command, error)
}

func GetCommand(getRequest Request, commandProvider CommandProvider, logger *logrus.Logger) (*pg.Command, error) {
	// TODO: Получать команду по двум ключам, а не по одному
	existCommand, err := commandProvider.GetOneBy("name", getRequest.Name)

	if existCommand == nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get command")
		return nil, nil
	}

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get command")
		return nil, dto.InternalError
	}

	return existCommand, nil
}
