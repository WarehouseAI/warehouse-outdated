package service

import (
	"context"
	"time"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/sirupsen/logrus"
)

func Authenticate(sessionId string, session dataservice.SessionInterface, logger *logrus.Logger) (*string, *m.Session, *e.ErrorResponse) {
	userId, newSession, err := session.Update(context.Background(), sessionId)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Authenticate")
		return nil, nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return userId, newSession, nil
}
