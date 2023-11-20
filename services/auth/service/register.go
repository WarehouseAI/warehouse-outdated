package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"warehouseai/auth/adapter"
	"warehouseai/auth/adapter/grpc/gen"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Picture   *multipart.FileHeader
	Email     string `json:"email"`
	ViaGoogle bool   `json:"via_google"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
}

type RegisterVerifyRequest struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type RegisterVerifyResponse struct {
	Verified bool `json:"verified"`
}

func Register(
	req *RegisterRequest,
	user adapter.UserGrpcInterface,
	verificationToken dataservice.VerificationTokenInterface,
	picture dataservice.PictureInterface,
	mail adapter.MailProducerInterface,
	logger *logrus.Logger,
) (*RegisterResponse, *e.ErrorResponse) {

	if len(req.Password) > 72 {
		return nil, e.NewErrorResponse(e.HttpBadRequest, "Password is too long.")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	var pictureUrl string

	if req.Picture != nil {
		picturePayload, err := req.Picture.Open()
		if err != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register user")
			return nil, e.NewErrorResponse(e.HttpInternalError, "Can't read provided avatar.")
		}

		url, fileErr := picture.UploadFile(picturePayload, fmt.Sprintf("%s_avatar%s", req.Username, filepath.Ext(req.Picture.Filename)))
		if fileErr != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": fileErr.Error()}).Info("Register user")
			return nil, e.NewErrorResponse(e.HttpInternalError, "Can't upload user avatar.")
		}

		defer picturePayload.Close()
		pictureUrl = url
	} else {
		pictureUrl = ""
	}

	// Move token create here to avoid create-missmatch with user on other service
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

	// Create user
	userId, gwErr := user.Create(context.Background(), &gen.CreateUserMsg{Firstname: req.Firstname, Lastname: req.Lastname, Username: req.Username, Password: string(hash), Picture: pictureUrl, Email: req.Email, ViaGoogle: req.ViaGoogle})

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Register user")
		return nil, gwErr
	}

	// Store verification token
	verificationTokenItem := m.VerificationToken{
		ID:        uuid.Must(uuid.NewV4()),
		UserId:    userId,
		Token:     string(tokenHash),
		ExpiresAt: time.Now().Add(time.Minute * 10),
		CreatedAt: time.Now(),
	}

	if err := verificationToken.Create(&verificationTokenItem); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Register user")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, "Failed to save the verification code")
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

	if err := mail.SendEmail(message); err != nil {
		fmt.Println(err)
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Send email")
		return nil, e.NewErrorResponse(e.HttpInternalError, "Failed to send email.")
	}

	return &RegisterResponse{UserId: userId}, nil
}

func RegisterVerify(
	request RegisterVerifyRequest,
	user adapter.UserGrpcInterface,
	verificationToken dataservice.VerificationTokenInterface,
	logger *logrus.Logger,
) (*RegisterVerifyResponse, *e.ErrorResponse) {
	existVerificationToken, dbErr := verificationToken.Get(map[string]interface{}{"user_id": request.UserId})

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Register verify user")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existVerificationToken.Token), []byte(request.Token)); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Register verify user")
		return nil, e.NewErrorResponse(e.HttpBadRequest, "Invalid register verification key")
	}

	verified, gwErr := user.UpdateVerificationStatus(context.Background(), request.UserId)

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Register verify user")
		return nil, gwErr
	}

	if err := verificationToken.Delete(map[string]interface{}{"id": existVerificationToken.ID}); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Register verify user")
		return nil, e.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return &RegisterVerifyResponse{Verified: verified}, nil
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
