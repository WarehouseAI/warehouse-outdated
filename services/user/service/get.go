package service

import (
	"context"
	"time"
	errs "warehouseai/user/errors"
	"warehouseai/user/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(map[string]interface{}) (*model.User, *errs.DBError)
	GetOneByPreload(map[string]interface{}, string) (*model.User, *errs.DBError)
}

func GetByEmail(email string, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*model.User, *errs.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy(map[string]interface{}{"email": email})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Email")
		return nil, errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetById(id string, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*model.User, *errs.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy(map[string]interface{}{"id": id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Id")
		return nil, errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetUserFavoriteAi(userId string, userProvider UserProvider, logger *logrus.Logger) ([]uuid.UUID, *errs.ErrorResponse) {
	user, dbErr := userProvider.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user favorite ai")
		return nil, errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return user.FavoriteAi, nil
}
