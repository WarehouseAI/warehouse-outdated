package register

import (
	"context"
	"errors"
	"time"
	"warehouse/gen"
	"warehouse/src/internal/dto"
	m "warehouse/src/services/auth/pkg/models"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserProvider interface {
	Create(ctx context.Context, userInfo *gen.CreateUserMsg) (*m.RegisterResponse, error)
}

func Register(userInfo *gen.CreateUserMsg) (*m.RegisterResponse, error) {
	if len(userInfo.Password) > 72 {
		return nil, errors.New("Password is too long")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	userInfo.Password = string(hash)
	userId, err := gw.CreateUser(ctx, userInfo)

	if err != nil && errors.Is(err, status.Errorf(codes.AlreadyExists, dto.ExistError.Error())) {
		pvd.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.ExistError
	}

	if err != nil && errors.Is(err, status.Errorf(codes.InvalidArgument, dto.BadRequestError.Error())) {
		pvd.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.BadRequestError
	}

	if err != nil {
		pvd.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.InternalError
	}

	return userId, nil
}
