package service

import (
	"context"
	"time"

	"warehouse/gen"

	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"
	"warehouse/src/internal/utils/mapper"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService interface {
	Create(context.Context, *gen.CreateUserMsg) (*dbm.User, error)
	GetByEmail(context.Context, *gen.GetUserByEmailMsg) (*dbm.User, error)
	GetById(context.Context, *gen.GetUserByIdMsg) (*dbm.User, error)
}

type UserServiceConfig struct {
	database *gorm.DB
	logger   *logrus.Logger
}

func NewUserService(database *gorm.DB, logger *logrus.Logger) UserService {
	return &UserServiceConfig{
		database: database,
		logger:   logger,
	}
}

func (cfg *UserServiceConfig) Create(ctx context.Context, userInfo *gen.CreateUserMsg) (*dbm.User, error) {
	userOperations := dbo.NewUserOperations(cfg.database)
	userEntity := mapper.UserPayloadToEntity(userInfo)

	existUser, _ := userOperations.GetOneBy("email", userEntity.Email)

	if existUser != nil {
		return nil, dto.ExistError
	}

	err := userOperations.Add(userEntity)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add user")
		return nil, dto.InternalError
	}

	return userEntity, nil
}

func (cfg *UserServiceConfig) GetByEmail(ctx context.Context, userInfo *gen.GetUserByEmailMsg) (*dbm.User, error) {
	userOperations := dbo.NewUserOperations(cfg.database)
	existUser, err := userOperations.GetOneBy("email", userInfo.Email)

	if err != nil {
		return nil, dto.InternalError
	}

	if existUser == nil {
		return nil, dto.NotFoundError
	}

	return existUser, nil
}

func (cfg *UserServiceConfig) GetById(ctx context.Context, userInfo *gen.GetUserByIdMsg) (*dbm.User, error) {
	userOperations := dbo.NewUserOperations(cfg.database)
	existUser, err := userOperations.GetOneBy("id", userInfo.Id)

	if err != nil {
		return nil, dto.InternalError
	}

	if existUser == nil {
		return nil, dto.NotFoundError
	}

	return existUser, nil
}
