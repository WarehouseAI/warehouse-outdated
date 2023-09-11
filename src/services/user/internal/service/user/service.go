package service

import (
	"context"
	"time"

	"warehouse/gen"
	im "warehouse/src/internal/models"
	d "warehouse/src/services/user/internal/datastore"
	m "warehouse/src/services/user/pkg/model"

	"github.com/sirupsen/logrus"
)

type UserServiceConfig struct {
	operations d.UserDatabaseOperations
	logger     *logrus.Logger
}

func NewUserService(operations d.UserDatabaseOperations, logger *logrus.Logger) m.UserService {
	return &UserServiceConfig{
		operations: operations,
		logger:     logger,
	}
}

func (cfg *UserServiceConfig) Create(ctx context.Context, userInfo *gen.CreateUserRequest) (*d.User, error) {
	userEntity := m.UserPayloadToEntity(userInfo)

	existUser, _ := cfg.operations.GetOneBy("email", userEntity.Email)

	if existUser != nil {
		return nil, im.ExistError
	}

	err := cfg.operations.Add(userEntity)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add user")
		return nil, im.InternalError
	}

	return userEntity, nil
}

func (cfg *UserServiceConfig) Get(ctx context.Context, userInfo *gen.GetUserRequest) (*d.User, error) {
	existUser, err := cfg.operations.GetOneBy("email", userInfo.Email)

	if err != nil {
		return nil, im.InternalError
	}

	if existUser == nil {
		return nil, im.NotFoundError
	}

	return existUser, nil
}
