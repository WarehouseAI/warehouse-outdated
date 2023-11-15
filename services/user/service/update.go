package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"
	"warehouseai/user/adapter"
	d "warehouseai/user/dataservice"
	e "warehouseai/user/errors"
	"warehouseai/user/model"
	m "warehouseai/user/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UpdateVerificationRequest struct {
	Verified         bool    `json:"verified"`
	VerificationCode *string `json:"verification_code"`
}

type UpdatePersonalDataRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
}

type UpdateUserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

type UpdateUserEmailRequest struct {
	Email            string `json:"email"`
	VerificationCode string `json:"-"`
	Verified         bool   `json:"-"`
}

type UserUpdater interface {
	RawUpdate(string, interface{}) (*m.User, *e.DBError)
}

// TODO: добавить обновление избранного

func UpdateUserVerification(request UpdateVerificationRequest, existUser *m.User, user d.UserInterface, logger *logrus.Logger) *e.ErrorResponse {
	if existUser.Verified {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Already verified"}).Info("Verify user")
		return e.NewErrorResponse(e.HttpBadRequest, "User already verified")
	}

	if request.VerificationCode != existUser.VerificationCode {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Invalid verification code"}).Info("Verify user")
		return e.NewErrorResponse(e.HttpBadRequest, "Invalid verification code")
	}

	request.VerificationCode = nil

	if _, dbErr := user.RawUpdate(existUser.ID.String(), request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Verify user")
		return e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func UpdateUserPersonalData(request UpdatePersonalDataRequest, userId string, user d.UserInterface, logger *logrus.Logger) (*m.User, *e.ErrorResponse) {
	updatedUser, dbErr := user.RawUpdate(userId, request)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return updatedUser, nil
}

func UpdateUserPassword(request UpdateUserPasswordRequest, existUser *m.User, user d.UserInterface, logger *logrus.Logger) *e.ErrorResponse {
	if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(request.OldPassword)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user password")
		return e.NewErrorResponse(e.HttpBadRequest, err.Error())
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = string(hash)

	if _, dbErr := user.RawUpdate(existUser.ID.String(), request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user password")
		return e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func UpdateUserEmail(request UpdateUserEmailRequest, userId string, user d.UserInterface, mail adapter.MailProducerInterface, logger *logrus.Logger) *e.ErrorResponse {
	key, err := generateKey(64)

	if err != nil {
		return e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	// TODO: переработать систему подтверждения. Производить изменения после перехода пользователем по ссылке.
	// Можно использовать отдельную таблицу для этого.
	request.Verified = false
	request.VerificationCode = key

	existUser, dbErr := user.RawUpdate(userId, request)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update email")
		return e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	message := model.Email{
		To:      request.Email,
		Subject: "Изменение электронной почты",
		Message: fmt.Sprintf(`
      Здравствуйте, %s!
      
      Мы получили запрос на изменение электронной почты от аккаунта %s.
      Перейдите по этой ссылке для подтверждения: %s
      
      Если изменение электронной почты больше не требуется, или вы не делали этот запрос - проигнорируйте данное письмо.
      
      WarehouseAI Team
      `, existUser.Firstname, existUser.Username, fmt.Sprintf("%s/api/user/verify/%s", os.Getenv("DOMAIN"), key)),
	}

	if err := mail.SendEmail(message); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Send email")
		return e.NewErrorResponse(e.HttpInternalError, "Failed to send email.")
	}

	return nil
}

func generateKey(length int) (string, error) {
	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(randomBytes)
	key = key[:length]

	return key, nil
}
