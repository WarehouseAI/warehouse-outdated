package service

import (
	"context"
	"time"
	errs "warehouseai/user/errors"
	"warehouseai/user/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type UserCreator interface {
	Add(item *model.User) *errs.DBError
}

func Create(userInfo model.User, userCreator UserCreator, logger *logrus.Logger, ctx context.Context) (*string, *errs.ErrorResponse) {
	userInfo.ID = uuid.Must(uuid.NewV4())

	if err := userCreator.Add(&userInfo); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Create user")
		return nil, errs.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	plainId := userInfo.ID.String()

	return &plainId, nil
}
