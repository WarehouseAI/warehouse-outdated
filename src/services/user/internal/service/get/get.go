package get

import (
	"context"
	"time"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(key string, payload interface{}) (*pg.User, error)
}

func GetByEmail(userInfo *gen.GetUserByEmailMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, error) {
	existUser, err := userProvider.GetOneBy("email", userInfo.Email)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get user by Email")
		return nil, dto.InternalError
	}

	if existUser == nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get user by Email")
		return nil, dto.NotFoundError
	}

	return existUser, nil
}

func GetById(userInfo *gen.GetUserByIdMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, error) {
	existUser, err := userProvider.GetOneBy("id", userInfo.Id)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get user by Id")
		return nil, dto.InternalError
	}

	if existUser == nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Get user by Id")
		return nil, dto.NotFoundError
	}

	return existUser, nil
}
