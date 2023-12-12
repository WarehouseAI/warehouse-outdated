package register

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/mail"
	"os"
	"time"
	"warehouseai/auth/adapter"
	"warehouseai/auth/adapter/grpc/gen"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Image     string `json:"image"`
	Email     string `json:"email"`
	ViaGoogle bool   `json:"via_google"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
}

func generateToken(length int) (string, error) {
	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(randomBytes)
	key = key[:length]

	return key, nil
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	return string(hash)
}

func validateRegisterRequest(req *RegisterRequest) *e.ErrorResponse {
	if len(req.Password) > 72 {
		return e.NewErrorResponse(e.HttpBadRequest, "Password is too long")
	}

	if len(req.Password) < 8 {
		return e.NewErrorResponse(e.HttpBadRequest, "Password is too short")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return e.NewErrorResponse(e.HttpBadRequest, "The provided string is not email")
	}

	return nil
}

func Register(
	req *RegisterRequest,
	user adapter.UserGrpcInterface,
	tokenRepository dataservice.VerificationTokenInterface,
	broker adapter.BrokerInterface,
	logger *logrus.Logger,
) (*RegisterResponse, *e.ErrorResponse) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Create user
	userId, gwErr := user.Create(context.Background(), &gen.CreateUserMsg{Firstname: req.Firstname, Lastname: req.Lastname, Username: req.Username, Password: hashPassword(req.Password), Picture: req.Image, Email: req.Email, ViaGoogle: req.ViaGoogle})

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Register user")
		return nil, gwErr
	}

	token, err := generateToken(12)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, e.NewErrorResponse(e.HttpInternalError, "Failed to create the verification code")
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), 12)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
		return nil, e.NewErrorResponse(e.HttpInternalError, "Failed to encrypt the verification code")
	}

	// Store verification token
	verificationTokenItem := m.VerificationToken{
		UserId:    userId,
		Token:     string(tokenHash),
		ExpiresAt: time.Now().Add(time.Minute * 10),
		CreatedAt: time.Now(),
	}

	if err := tokenRepository.Create(&verificationTokenItem); err != nil {
		if err := broker.SendTokenReject(userId); err != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
			return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
		}

		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Register user")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	message := m.Email{
		To:      req.Email,
		Subject: "Подтверждение электронной почты",
		Message: fmt.Sprintf(`
      Здравствуйте, %s!
      
      Для завершения регистрации перейдите, пожалуйста, по ссылке:
      %s
      
      Если вы не указывали эту электронную почту - проигнорируйте данное письмо.
      
      WarehouseAI Team
      `, req.Firstname, fmt.Sprintf("%s/register/confirm?user=%s&token=%s", os.Getenv("DOMAIN"), userId, token)),
	}

	if err := broker.SendEmail(message); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Send email")
		return nil, e.NewErrorResponse(e.HttpInternalError, "Failed to send email.")
	}

	return &RegisterResponse{UserId: userId}, nil
}
