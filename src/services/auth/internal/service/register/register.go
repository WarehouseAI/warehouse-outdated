package register

import (
	"context"
	"errors"
	"time"
	"warehouse/gen"
	"warehouse/src/internal/dto"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Response struct {
	ID string `json:"id"`
}

type UserCreator interface {
	Create(ctx context.Context, userInfo *gen.CreateUserMsg) (*Response, error)
}

func Register(userInfo *gen.CreateUserMsg, userCreator UserCreator, logger *logrus.Logger, ctx context.Context) (*Response, error) {
	if len(userInfo.Password) > 72 {
		return nil, errors.New("Password is too long")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	userInfo.Password = string(hash)
	userId, err := userCreator.Create(ctx, userInfo)

	if err != nil && errors.Is(err, status.Errorf(codes.AlreadyExists, dto.ExistError.Error())) {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.ExistError
	}

	if err != nil && errors.Is(err, status.Errorf(codes.InvalidArgument, dto.BadRequestError.Error())) {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.BadRequestError
	}

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, dto.InternalError
	}

	return userId, nil
}
