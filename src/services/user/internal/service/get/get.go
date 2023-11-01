package get

import (
	"context"
	"time"
	"warehouse/gen"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type GetFavoriteAIResponse struct {
	FavoriteAi []*pg.AI `json:"favorite_ai"`
}

type UserProvider interface {
	GetOneBy(map[string]interface{}) (*pg.User, *db.DBError)
	GetOneByPreload(map[string]interface{}, string) (*pg.User, *db.DBError)
}

func GetByEmail(userInfo *gen.GetUserByEmailMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, *httputils.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy(map[string]interface{}{"email": userInfo.Email})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Email")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetById(userInfo *gen.GetUserByIdMsg, userProvider UserProvider, logger *logrus.Logger, ctx context.Context) (*pg.User, *httputils.ErrorResponse) {
	existUser, dbErr := userProvider.GetOneBy(map[string]interface{}{"id": userInfo.Id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Id")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetUserFavoriteAi(userId string, userProvider UserProvider, logger *logrus.Logger) (*GetFavoriteAIResponse, *httputils.ErrorResponse) {
	user, dbErr := userProvider.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user favorite ai")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &GetFavoriteAIResponse{user.FavoriteAi}, nil
}
