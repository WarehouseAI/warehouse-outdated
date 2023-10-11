package service

import (
	"context"
	"errors"
	"time"
	"warehouse/gen"
	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"

	"warehouse/src/internal/dto"
	gw "warehouse/src/services/auth/internal/gateway"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Пофиксить разные пакеты (m и gen и dbm)
type AuthService interface {
	Login(context.Context, *m.LoginRequest) (*dbm.Session, error)
	Register(context.Context, *gen.CreateUserMsg) (*m.RegisterResponse, error)
	Logout(context.Context, string) error
}

type AuthServiceConfig struct {
	database *redis.Client
	logger   *logrus.Logger
}

func NewAuthService(database *redis.Client, logger *logrus.Logger) AuthService {
	return &AuthServiceConfig{
		database: database,
		logger:   logger,
	}
}

func (cfg *AuthServiceConfig) Logout(ctx context.Context, sessionId string) error {
	sessionOperations := dbo.NewSessionOperations(cfg.database)

	if err := sessionOperations.DeleteSession(ctx, sessionId); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Logout user")
		return err
	}

	return nil
}

func (cfg *AuthServiceConfig) Login(ctx context.Context, userInfo *m.LoginRequest) (*dbm.Session, error) {
	sessionOperations := dbo.NewSessionOperations(cfg.database)
	user, err := gw.GetUserByEmail(ctx, &gen.GetUserByEmailMsg{Email: userInfo.Email})

	if err != nil && errors.Is(err, status.Errorf(codes.NotFound, dto.NotFoundError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.NotFoundError
	}

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.InternalError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.BadRequestError
	}

	// Сохраняем сессию
	session, err := sessionOperations.CreateSession(ctx, user.Id)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.InternalError
	}

	return session, nil
}

func (cfg *AuthServiceConfig) Register(ctx context.Context, userInfo *gen.CreateUserMsg) (*m.RegisterResponse, error) {
	if len(userInfo.Password) > 72 {
		return nil, errors.New("Password is too long")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	userInfo.Password = string(hash)
	userId, err := gw.CreateUser(ctx, userInfo)

	if err != nil && errors.Is(err, status.Errorf(codes.AlreadyExists, dto.ExistError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.ExistError
	}

	if err != nil && errors.Is(err, status.Errorf(codes.InvalidArgument, dto.BadRequestError.Error())) {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.BadRequestError
	}

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.InternalError
	}

	return userId, nil
}
