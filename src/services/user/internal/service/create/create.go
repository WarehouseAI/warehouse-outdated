package create

import (
	"context"
	"time"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
	"warehouse/src/internal/utils/grpcutils"

	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(key string, payload interface{}) (*pg.User, error)
	Add(item *pg.User) error
}

func Create(userInfo *gen.CreateUserMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, error) {
	userEntity := grpcutils.UserPayloadToEntity(userInfo)

	existUser, _ := userProvider.GetOneBy("email", userEntity.Email)

	if existUser != nil {
		return nil, dto.ExistError
	}

	err := userProvider.Add(userEntity)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create user")
		return nil, dto.InternalError
	}

	return userEntity, nil
}
