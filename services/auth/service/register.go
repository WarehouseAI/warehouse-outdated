package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"warehouseai/auth/adapter"
	"warehouseai/auth/adapter/grpc/gen"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"

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

func Register(req *RegisterRequest, user adapter.UserGrpcInterface, picture dataservice.PictureInterface, logger *logrus.Logger) (*RegisterResponse, *e.ErrorResponse) {
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

	userId, gwErr := user.Create(context.Background(), &gen.CreateUserMsg{Firstname: req.Firstname, Lastname: req.Lastname, Username: req.Username, Password: string(hash), Picture: pictureUrl, Email: req.Email, ViaGoogle: req.ViaGoogle})

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Register user")
		return nil, gwErr
	}

	return &RegisterResponse{UserId: *userId}, nil
}
