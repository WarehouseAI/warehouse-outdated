package login

import (
	"context"
	"time"
	"warehouse/gen"
	db "warehouse/src/internal/database"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionCreator interface {
	Create(context.Context, string) (*r.Session, *db.DBError)
}

type UserProvider interface {
	GetByEmail(context.Context, *gen.GetUserByEmailMsg) (*gen.User, *httputils.ErrorResponse)
}

func Login(userInfo *Request, userProvider UserProvider, sessionCreator SessionCreator, logger *logrus.Logger, ctx context.Context) (*r.Session, *httputils.ErrorResponse) {
	user, gwErr := userProvider.GetByEmail(ctx, &gen.GetUserByEmailMsg{Email: userInfo.Email})

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Login user")
		return nil, gwErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInfo.Password)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, httputils.NewErrorResponse(httputils.BadRequest, "Invalid credentials")
	}

	// Сохраняем сессию
	session, sessErr := sessionCreator.Create(ctx, user.Id)

	if sessErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": sessErr}).Info("Login user")
		return nil, httputils.NewErrorResponseFromDBError(sessErr.ErrorType, sessErr.Message)
	}

	return session, nil
}
