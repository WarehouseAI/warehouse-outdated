package login

import (
	"context"
	"net/mail"
	"time"
	"warehouseai/auth/adapter"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"
	"warehouseai/auth/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId string `json:"user_id"`
}

func validateLoginRequest(req *LoginRequest) *e.ErrorResponse {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return e.NewErrorResponse(e.HttpBadRequest, "Invalid email address")
	}

	return nil
}

func Login(req *LoginRequest, user adapter.UserGrpcInterface, session dataservice.SessionInterface, logger *logrus.Logger) (*LoginResponse, *model.Session, *e.ErrorResponse) {
	if err := validateLoginRequest(req); err != nil {
		return nil, nil, err
	}

	existUser, gwErr := user.GetByEmail(context.Background(), req.Email)

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Login user")
		return nil, nil, gwErr
	}

	if !existUser.Verified {
		return nil, nil, e.NewErrorResponse(e.HttpForbidden, "Verify your email first")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(req.Password)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Login user")
		return nil, nil, e.NewErrorResponse(e.HttpBadRequest, "Invalid credentials")
	}

	newSession, sessErr := session.Create(context.Background(), existUser.Id)

	if sessErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": sessErr}).Info("Login user")
		return nil, nil, e.NewErrorResponseFromDBError(sessErr.ErrorType, sessErr.Message)
	}

	return &LoginResponse{UserId: existUser.Id}, newSession, nil
}
