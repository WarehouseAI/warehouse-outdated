package service

import (
	"context"
	"errors"
	"time"
	"warehouse/gen"
	im "warehouse/src/internal/models"
	gw "warehouse/src/services/auth/internal/gateway"
	"warehouse/src/services/auth/pkg/model"
	m "warehouse/src/services/auth/pkg/model"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceConfig struct {
	rClient *redis.Client
	logger  *logrus.Logger
}

func NewAuthService(rClient *redis.Client, logger *logrus.Logger) m.AuthService {
	return &AuthServiceConfig{
		rClient: rClient,
		logger:  logger,
	}
}

// func (s *AuthService) Refresh(ctx context.Context) (m.TokenPair, error) {

// }

// func (s *AuthService) Login(ctx context.Context) (m.TokenPair, error) {

// }

func (cfg *AuthServiceConfig) Register(ctx context.Context, userInfo *gen.CreateUserRequest) (*model.UserIdResponse, error) {
	userId, err := gw.CreateUser(ctx, userInfo)

	if err != nil && errors.Is(err, status.Errorf(codes.AlreadyExists, im.ExistError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, im.ExistError
	} else if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, im.InternalError
	}

	return userId, nil
}
