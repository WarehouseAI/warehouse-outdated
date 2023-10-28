package update

import (
	"fmt"
	"os"
	"sync"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
	"warehouse/src/internal/utils/mailutils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UpdateEmailRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"-"`
	Verified         bool   `json:"-"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
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

func UpdatePassword(request UpdatePasswordRequest, user *pg.User, userUpdater UserUpdater, logger *logrus.Logger) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		fmt.Println("Invalid pass")
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user password")
		return dto.BadRequestError
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = string(hash)

	_, err := userUpdater.Update(user.ID.String(), request)

	if err != nil {
		fmt.Println("Invalid update")
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user password")
		return dto.InternalError
	}

	return nil
}

func UpdateEmail(wg *sync.WaitGroup, respch chan error, request UpdateEmailRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) {
	if _, err := userUpdater.Update(userId, request); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update email")
		respch <- err
	} else {
		respch <- nil
	}

	wg.Done()
}

func SendUpdateNotification(wg *sync.WaitGroup, respch chan error, logger *logrus.Logger, request UpdateEmailRequest) {
	message := mailutils.NewMessage(mailutils.EmailVerify, request.Email, fmt.Sprintf("%s/api/user/verify/%s", os.Getenv("DOMAIN_HOST"), request.VerificationCode))

	if err := mailutils.SendEmail(message); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update email")
		respch <- err
	} else {
		respch <- nil
	}

	wg.Done()
}
