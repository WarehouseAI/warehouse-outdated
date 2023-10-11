package login

import (
	"context"
	"errors"
	"time"
	"warehouse/gen"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/dto"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SessionCreator interface {
	Create(context.Context, string) (*r.Session, error)
}

type UserProvider interface {
	GetByEmail(context.Context, *gen.GetUserByEmailMsg) (*gen.User, error)
}

func Login(userInfo *m.LoginRequest, userProvider UserProvider, sessionCreator SessionCreator, logger *logrus.Logger, ctx context.Context) (*r.Session, error) {
	user, err := userProvider.GetByEmail(ctx, &gen.GetUserByEmailMsg{Email: userInfo.Email})

	if err != nil && errors.Is(err, status.Errorf(codes.NotFound, dto.NotFoundError.Error())) {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.NotFoundError
	}

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.InternalError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.BadRequestError
	}

	// Сохраняем сессию
	session, err := sessionCreator.Create(ctx, user.Id)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, dto.InternalError
	}

	return session, nil
}
