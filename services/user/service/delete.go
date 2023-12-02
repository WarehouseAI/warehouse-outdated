package service

import (
	"time"
	"warehouseai/user/dataservice"
	e "warehouseai/user/errors"

	"github.com/sirupsen/logrus"
)

func Delete(userId string, user dataservice.UserInterface, logger *logrus.Logger) *e.ErrorResponse {
	if err := user.Delete(map[string]interface{}{"id": userId}); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Delete user")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
