package service

import (
	"context"
	"time"
	e "warehouseai/internal/errors"
	m "warehouseai/user/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type UserCreator interface {
	Add(item *m.User) *e.DBError
}

func Create(userInfo m.User, userCreator UserCreator, logger *logrus.Logger, ctx context.Context) (*string, *e.ErrorResponse) {
	userInfo.ID = uuid.Must(uuid.NewV4())

	if err := userCreator.Add(&userInfo); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Create user")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	plainId := userInfo.ID.String()

	return &plainId, nil
}
