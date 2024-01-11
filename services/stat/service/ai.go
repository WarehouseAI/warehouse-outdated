package service

import (
	"time"
	d "warehouseai/stat/dataservice"
	e "warehouseai/stat/errors"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetNumOfAiUses(id uuid.UUID, ai d.AiStatInterface, logger *logrus.Logger) (uint, *e.ErrorResponse) {
	num, dbErr := ai.GetNumOfAiUses(id)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get ai's stat")
		return 0, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return num, nil
}
