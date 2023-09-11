package grpc

import (
	"context"
	"errors"
	"warehouse/gen"
	im "warehouse/src/internal/models"
	m "warehouse/src/services/user/pkg/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserPrivateAPI struct {
	gen.UnimplementedUserServiceServer
	svc m.UserService
}

func NewUserPrivateAPI(svc m.UserService) *UserPrivateAPI {
	return &UserPrivateAPI{
		svc: svc,
	}
}

func (api *UserPrivateAPI) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, im.BadRequestError.Error())
	}

	user, err := api.svc.Create(ctx, req)

	if err != nil && errors.Is(err, im.ExistError) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.CreateUserResponse{Id: m.UserToProto(user).Id}, nil
}

func (api *UserPrivateAPI) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, im.BadRequestError.Error())
	}

	user, err := api.svc.Get(ctx, req)

	if err != nil && errors.Is(err, im.NotFoundError) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetUserResponse{User: m.UserToProto(user)}, nil
}
