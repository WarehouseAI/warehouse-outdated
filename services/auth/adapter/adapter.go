package adapter

import (
	"context"
	"warehouseai/auth/adapter/grpc/gen"
	e "warehouseai/auth/errors"
	"warehouseai/auth/model"
)

type UserGrpcInterface interface {
	Create(ctx context.Context, userInfo *gen.CreateUserMsg) (string, *e.ErrorResponse)
	ResetPassword(ctx context.Context, resetPasswordRequest *gen.ResetPasswordRequest) (string, *e.ErrorResponse)
	GetByEmail(ctx context.Context, email string) (*gen.User, *e.ErrorResponse)
	GetById(ctx context.Context, userId string) (*gen.User, *e.ErrorResponse)
	UpdateVerificationStatus(ctx context.Context, userId string) (bool, *e.ErrorResponse)
}

type BrokerInterface interface {
	SendEmail(email model.Email) error
	SendTokenReject(userId string) error
}
