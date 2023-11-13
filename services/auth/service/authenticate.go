package service

import (
	"context"
	"time"
	"warehouseai/auth/dataservice"
	e "warehouseai/internal/errors"

	"github.com/sirupsen/logrus"
)

func Authenticate(sessionId string, session dataservice.SessionInterface, logger *logrus.Logger) (*string, *e.ErrorResponse) {
	newSession, err := session.Update(context.Background(), sessionId)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Authenticate")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return &newSession.ID, nil
}
