package service

import (
	"time"
	d "warehouseai/stat/dataservice"
	e "warehouseai/stat/errors"

	"github.com/sirupsen/logrus"
)

func GetNumOfUsers(db d.UsersStatInterface, logger *logrus.Logger) (uint, *e.ErrorResponse) {
	num, dbErr := db.GetNumOfUsers()

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get numbers of users")
		return 0, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return num, nil
}

func GetNumOfDevelopers(db d.UsersStatInterface, logger *logrus.Logger) (uint, *e.ErrorResponse) {
	num, dbErr := db.GetNumOfDevelopers()

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get numbers of developers")
		return 0, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return num, nil
}
