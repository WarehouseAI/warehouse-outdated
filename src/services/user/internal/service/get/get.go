package get

import (
	"context"
	"time"
	"warehouse/gen"
	"warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(key string, payload interface{}) (*pg.User, *database.DBError)
}

func GetByEmail(userInfo *gen.GetUserByEmailMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, *httputils.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy("email", userInfo.Email)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Email")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetById(userInfo *gen.GetUserByIdMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, *httputils.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy("id", userInfo.Id)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Id")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}
