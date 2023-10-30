package register

import (
	"context"
	"time"
	"warehouse/gen"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	ID string `json:"id"`
}

type UserCreator interface {
	Create(ctx context.Context, userInfo *gen.CreateUserMsg) (*Response, *httputils.ErrorResponse)
}

func Register(userInfo *gen.CreateUserMsg, userCreator UserCreator, logger *logrus.Logger, ctx context.Context) (*Response, *httputils.ErrorResponse) {
	if len(userInfo.Password) > 72 {
		return nil, httputils.NewErrorResponse(httputils.BadRequest, "Password is too long.")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	userInfo.Password = string(hash)
	userId, gwErr := userCreator.Create(ctx, userInfo)

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Register user")
		return nil, gwErr
	}

	return userId, nil
}
