package service

import (
	"time"
	d "warehouseai/user/dataservice"
	e "warehouseai/user/errors"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type ResetUserPasswordRequest struct {
	Password string `json:"password"`
}

func ResetUserPassword(request ResetUserPasswordRequest, userId string, user d.UserInterface, logger *logrus.Logger) *e.ErrorResponse {
	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)

	if dbErr := user.Update(userId, map[string]interface{}{"password": string(hash)}); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user password")
		return e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}
