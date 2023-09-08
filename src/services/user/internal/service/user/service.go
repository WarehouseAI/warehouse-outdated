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
	operations d.Operations
	logger     *logrus.Logger
}

func NewUserService(operations d.Operations, logger *logrus.Logger) m.UserService {
	return &UserServiceConfig{
		operations: operations,
		logger:     logger,
	}
}

func (cfg *UserServiceConfig) Create(ctx context.Context, userInfo *gen.CreateUserRequest) (*d.User, error) {
	user := m.UserPayloadToEntity(userInfo)

	existUser, _ := cfg.operations.GetOneBy("email", user.Email)

	if existUser != nil {
		return nil, im.ExistError
	}

	err := cfg.operations.Add(user)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add user")
		return nil, im.InternalError
	}

	return user, nil
}
