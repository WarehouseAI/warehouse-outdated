package service

import (
	"context"
	"time"
	"warehouseai/auth/adapter"
	e "warehouseai/internal/errors"
	"warehouseai/internal/gen"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Picture   string `json:"picture"`
	Email     string `json:"email"`
	ViaGoogle bool   `json:"via_google"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
}

func Register(req *RegisterRequest, user adapter.UserGrpcInterface, logger *logrus.Logger, ctx context.Context) (*RegisterResponse, *e.ErrorResponse) {
	if len(req.Password) > 72 {
		return nil, e.NewErrorResponse(e.HttpBadRequest, "Password is too long.")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	userId, gwErr := user.Create(ctx, &gen.CreateUserMsg{Firstname: req.Firstname, Lastname: req.Lastname, Username: req.Username, Password: string(hash), Picture: req.Picture, Email: req.Email, ViaGoogle: req.ViaGoogle})

	if gwErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": gwErr.ErrorMessage}).Info("Register user")
		return nil, gwErr
	}

	return &RegisterResponse{UserId: *userId}, nil
}
