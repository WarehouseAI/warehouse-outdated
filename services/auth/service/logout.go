package service

import (
	"context"
	"time"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"

	"github.com/sirupsen/logrus"
)

func Logout(sessionId string, session dataservice.SessionInterface, logger *logrus.Logger) *e.ErrorResponse {
	if err := session.Delete(context.Background(), sessionId); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Logout user")
		return e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
