package service

import (
	"context"
	"time"
	e "warehouseai/internal/errors"
	m "warehouseai/user/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(map[string]interface{}) (*m.User, *e.DBError)
	GetOneByPreload(map[string]interface{}, string) (*m.User, *e.DBError)
}

func GetByEmail(email string, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*m.User, *e.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy(map[string]interface{}{"email": email})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Email")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetById(id string, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*m.User, *e.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy(map[string]interface{}{"id": id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Id")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetUserFavoriteAi(userId string, userProvider UserProvider, logger *logrus.Logger) ([]uuid.UUID, *e.ErrorResponse) {
	user, dbErr := userProvider.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user favorite ai")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return user.FavoriteAi, nil
}
