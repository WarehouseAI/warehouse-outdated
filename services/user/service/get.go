package service

import (
	"time"
	d "warehouseai/user/dataservice"
	e "warehouseai/user/errors"
	m "warehouseai/user/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

func GetByEmail(email string, user d.UserInterface, logger *logrus.Logger) (*m.User, *e.ErrorResponse) {
	existUser, dbErr := user.GetOneBy(map[string]interface{}{"email": email})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Email")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetById(id string, user d.UserInterface, logger *logrus.Logger) (*m.User, *e.ErrorResponse) {
	existUser, dbErr := user.GetOneBy(map[string]interface{}{"id": id})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user by Id")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser, nil
}

func GetUserFavoriteAi(userId string, user d.UserInterface, logger *logrus.Logger) ([]uuid.UUID, *e.ErrorResponse) {
	existUser, dbErr := user.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get user favorite ai")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return existUser.FavoriteAi, nil
}
