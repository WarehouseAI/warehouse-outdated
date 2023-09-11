package service

import (
	"context"
	"errors"
	"time"
	"warehouse/gen"
	im "warehouse/src/internal/models"
	d "warehouse/src/services/auth/internal/datastore"
	gw "warehouse/src/services/auth/internal/gateway"
	"warehouse/src/services/auth/pkg/model"
	m "warehouse/src/services/auth/pkg/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceConfig struct {
	operations d.SessionDatabaseOperations
	logger     *logrus.Logger
}

func NewAuthService(operations d.SessionDatabaseOperations, logger *logrus.Logger) m.AuthService {
	return &AuthServiceConfig{
		operations: operations,
		logger:     logger,
	}
}

// func (s *AuthService) Refresh(ctx context.Context) (m.TokenPair, error) {

// }

func (cfg *AuthServiceConfig) Login(ctx context.Context, userInfo *model.LoginRequest) (*m.Session, error) {
	user, err := gw.GetUser(ctx, &gen.GetUserRequest{Email: userInfo.Email})

	if err != nil && errors.Is(err, status.Errorf(codes.NotFound, im.NotFoundError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, im.NotFoundError
	}

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, im.InternalError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.User.Password), []byte(userInfo.Password)); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, im.BadRequestError
	}

	// Сохраняем сессию
	session, err := cfg.operations.CreateSession(ctx, user.User.Id)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, im.InternalError
	}

	return session, nil
}

func (cfg *AuthServiceConfig) Register(ctx context.Context, userInfo *gen.CreateUserRequest) (*model.UserIdResponse, error) {
	if len(userInfo.Password) > 72 {
		return nil, errors.New("Password is too long")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	userInfo.Password = string(hash)
	userId, err := gw.CreateUser(ctx, userInfo)

	if err != nil && errors.Is(err, status.Errorf(codes.AlreadyExists, im.ExistError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, im.ExistError
	}

	if err != nil && errors.Is(err, status.Errorf(codes.InvalidArgument, im.BadRequestError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, im.BadRequestError
	}

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, im.InternalError
	}

	return userId, nil
}
