package create

import (
	"context"
	"time"
	"warehouse/gen"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/grpcutils"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(key string, payload interface{}) (*pg.User, *db.DBError)
	Add(item *pg.User) *db.DBError
}

func Create(userInfo *gen.CreateUserMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, *httputils.ErrorResponse) {
	userEntity := grpcutils.UserPayloadToEntity(userInfo)

	err := userProvider.Add(userEntity)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Create user")
		return nil, httputils.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return userEntity, nil
}
