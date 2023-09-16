package service

import (
	"context"
	"time"

	"warehouse/gen"

	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"
	m "warehouse/src/services/user/pkg/models"

	"github.com/sirupsen/logrus"
)

type UserService interface {
	Create(context.Context, *gen.CreateUserRequest) (*dbm.User, error)
	Get(context.Context, *gen.GetUserRequest) (*dbm.User, error)
}

type UserServiceConfig struct {
	operations dbo.UserDatabaseOperations
	logger     *logrus.Logger
}

func NewUserService(operations dbo.UserDatabaseOperations, logger *logrus.Logger) UserService {
	return &UserServiceConfig{
		operations: operations,
		logger:     logger,
	}
}

func (cfg *UserServiceConfig) Create(ctx context.Context, userInfo *gen.CreateUserRequest) (*dbm.User, error) {
	userEntity := m.UserPayloadToEntity(userInfo)

	existUser, _ := cfg.operations.GetOneBy("email", userEntity.Email)

	if existUser != nil {
		return nil, dto.ExistError
	}

	err := cfg.operations.Add(userEntity)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add user")
		return nil, dto.InternalError
	}

	return userEntity, nil
}

func (cfg *UserServiceConfig) Get(ctx context.Context, userInfo *gen.GetUserRequest) (*dbm.User, error) {
	existUser, err := cfg.operations.GetOneBy("email", userInfo.Email)

	if err != nil {
		return nil, dto.InternalError
	}

	if existUser == nil {
		return nil, dto.NotFoundError
	}

	return existUser, nil
}
