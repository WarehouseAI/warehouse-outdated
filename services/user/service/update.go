package service

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	errs "warehouseai/user/errors"
	"warehouseai/user/model"

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
	RawUpdate(string, interface{}) (*model.User, *errs.DBError)
}

// TODO: добавить обновление избранного

func UpdateUserVerification(request UpdateVerificationRequest, user *model.User, updater UserUpdater, logger *logrus.Logger) *errs.ErrorResponse {
	if user.Verified {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Already verified"}).Info("Verify user")
		return errs.NewErrorResponse(errs.HttpBadRequest, "User already verified")
	}

	if request.VerificationCode != user.VerificationCode {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": "Invalid verification code"}).Info("Verify user")
		return errs.NewErrorResponse(errs.HttpBadRequest, "Invalid verification code")
	}

	request.VerificationCode = nil

	if _, dbErr := updater.RawUpdate(user.ID.String(), request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Verify user")
		return errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func UpdateUserPersonalData(request UpdatePersonalDataRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) (*model.User, *errs.ErrorResponse) {
	updatedUser, dbErr := userUpdater.RawUpdate(userId, request)

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user")
		return nil, errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return updatedUser, nil
}

func UpdateUserPassword(request UpdateUserPasswordRequest, user *model.User, userUpdater UserUpdater, logger *logrus.Logger) *errs.ErrorResponse {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Update user password")
		return errs.NewErrorResponse(errs.HttpBadRequest, err.Error())
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	request.Password = string(hash)

	if _, dbErr := userUpdater.RawUpdate(user.ID.String(), request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update user password")
		return errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}

func UpdateUserEmail(request UpdateUserEmailRequest, userId string, userUpdater UserUpdater, logger *logrus.Logger) *errs.ErrorResponse {
	key, err := generateKey(64)

	if err != nil {
		return errs.NewErrorResponse(errs.HttpInternalError, err.Error())
	}

	request.Verified = false
	request.VerificationCode = key

	if _, dbErr := userUpdater.RawUpdate(userId, request); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Update email")
		return errs.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	// TODO: добавить отправление уведомления на email

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
