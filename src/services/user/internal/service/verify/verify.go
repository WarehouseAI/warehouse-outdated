package verify

import (
	"errors"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/sirupsen/logrus"
)

type Request struct {
	Verified         bool   `json:"verified"`
	VerificationCode string `json:"verification_code"`
}

type UserUpdater interface {
	Update(id string, updatedFields interface{}) (*pg.User, error)
}

func VerifyUserEmail(request Request, user *pg.User, userUpdater UserUpdater, logger *logrus.Logger) error {
	if user.Verified {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Already verified"}).Info("Verify user")
		return errors.New("User already verified")
	}

	if request.VerificationCode != user.VerificationCode {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dto.BadRequestError.Error()}).Info("Verify user")
		return dto.BadRequestError
	}

	request.VerificationCode = ""

	_, err := userUpdater.Update(user.ID.String(), request)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Verify user")
		return dto.InternalError
	}

	return nil
}
