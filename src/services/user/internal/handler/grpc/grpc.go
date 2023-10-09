package grpc

import (
	"context"
	"errors"
	"warehouse/gen"
	"warehouse/src/internal/dto"
	"warehouse/src/internal/utils/mapper"
	svc "warehouse/src/services/user/internal/service/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserPrivateAPI struct {
	gen.UnimplementedUserServiceServer
	svc svc.UserService
}

func NewUserPrivateAPI(svc svc.UserService) *UserPrivateAPI {
	return &UserPrivateAPI{
		svc: svc,
	}
}

func (api *UserPrivateAPI) CreateUser(ctx context.Context, req *gen.CreateUserMsg) (*gen.CreateUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, dto.BadRequestError.Error())
	}

	user, err := api.svc.Create(ctx, req)

	if err != nil && errors.Is(err, dto.ExistError) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.CreateUserResponse{Id: mapper.UserToProto(user).Id}, nil
}

func (api *UserPrivateAPI) GetUserByEmail(ctx context.Context, req *gen.GetUserByEmailMsg) (*gen.User, error) {
	if req == nil || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, dto.BadRequestError.Error())
	}

	user, err := api.svc.GetByEmail(ctx, req)

	if err != nil && errors.Is(err, dto.NotFoundError) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return mapper.UserToProto(user), nil
}

func (api *UserPrivateAPI) GetUserById(ctx context.Context, req *gen.GetUserByIdMsg) (*gen.User, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, dto.BadRequestError.Error())
	}

	user, err := api.svc.GetById(ctx, req)

	if err != nil && errors.Is(err, dto.NotFoundError) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return mapper.UserToProto(user), nil
}
