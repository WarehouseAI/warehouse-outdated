package update

import (
	"fmt"
	"os"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
	u "warehouse/src/internal/utils"
	"warehouse/src/internal/utils/mailutils"

	"github.com/sirupsen/logrus"
)

type UpdateEmailRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"-"`
	Verified         bool   `json:"-"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
}

type UserUpdater interface {
	Update(id string, updatedFields interface{}) (*pg.User, error)
}

func UpdateUser(request UpdateUserRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) (*pg.User, error) {
	updatedUser, err := userUpdater.Update(userId, request)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user")
		return nil, dto.BadRequestError
	}

	return updatedUser, nil
}

func UpdateEmail(request UpdateEmailRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) error {
	key, err := u.GenerateKey(64)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update email")
		return dto.InternalError
	}

	request.Verified = false
	request.VerificationCode = key

	message := mailutils.NewMessage(mailutils.EmailVerify, request.Email, fmt.Sprintf("%s/api/user/verify/%s", os.Getenv("DOMAIN_HOST"), key))

	if err := mailutils.SendEmail(message); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update email")
		return dto.InternalError
	}

	if _, err := userUpdater.Update(userId, request); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update email")
		return dto.InternalError
	}

	return nil
}
