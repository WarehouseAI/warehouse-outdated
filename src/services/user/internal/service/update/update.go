package update

import (
	"fmt"
	"os"
	"sync"
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"
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
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UserUpdater interface {
	Update(id string, updatedFields interface{}) (*pg.User, *db.DBError)
}

func UpdateUser(request UpdateUserRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) (*pg.User, *httputils.ErrorResponse) {
	updatedUser, dbErr := userUpdater.Update(userId, request)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return updatedUser, nil
}

func UpdatePassword(request UpdatePasswordRequest, user *pg.User, userUpdater UserUpdater, logger *logrus.Logger) *httputils.ErrorResponse {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user password")
		return httputils.NewErrorResponse(httputils.BadRequest, err.Error())
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = string(hash)

	if _, dbErr := userUpdater.Update(user.ID.String(), request); dbErr != nil {
		fmt.Println("Invalid update")
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user password")
		return httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func UpdateEmail(wg *sync.WaitGroup, respch chan *httputils.ErrorResponse, request UpdateEmailRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) {
	if _, dbErr := userUpdater.Update(userId, request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update email")
		respch <- httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	} else {
		respch <- nil
	}

	wg.Done()
}

func SendUpdateNotification(wg *sync.WaitGroup, respch chan *httputils.ErrorResponse, logger *logrus.Logger, request UpdateEmailRequest) {
	message := mailutils.NewMessage(mailutils.EmailVerify, request.Email, fmt.Sprintf("%s/api/user/verify/%s", os.Getenv("DOMAIN_HOST"), request.VerificationCode))

	if mailErr := mailutils.SendEmail(message); mailErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": mailErr.Error()}).Info("Update email")
		respch <- httputils.NewErrorResponse(httputils.InternalError, mailErr.Error())
	} else {
		respch <- nil
	}

	wg.Done()
}
