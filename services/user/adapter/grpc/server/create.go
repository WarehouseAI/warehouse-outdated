package server

import (
	"context"
	"warehouseai/user/adapter/grpc/gen"
	errs "warehouseai/user/errors"
	"warehouseai/user/model"
	"warehouseai/user/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pvd *UserGrpcServer) CreateUser(ctx context.Context, req *gen.CreateUserMsg) (*gen.CreateUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	newUser := model.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Lastname,
		Password:  req.Password,
		Picture:   req.Picture,
		Email:     req.Email,
	}

	userId, err := service.Create(newUser, pvd.db, pvd.logger, ctx)

	if err != nil {
		if err.ErrorCode == errs.HttpAlreadyExist {
			return nil, status.Errorf(codes.AlreadyExists, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return &gen.CreateUserResponse{Id: *userId}, nil
}
