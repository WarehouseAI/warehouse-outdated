package server

import (
	"context"
	"warehouseai/user/adapter/grpc/gen"
	"warehouseai/user/adapter/grpc/mapper"
	"warehouseai/user/dataservice"
	e "warehouseai/user/errors"
	m "warehouseai/user/model"
	"warehouseai/user/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGrpcServer struct {
	gen.UnimplementedUserServiceServer
	DB     dataservice.UserInterface
	Logger *logrus.Logger
}

func (s *UserGrpcServer) CreateUser(ctx context.Context, req *gen.CreateUserMsg) (*gen.CreateUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	newUser := m.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Lastname,
		Password:  req.Password,
		Picture:   req.Picture,
		Email:     req.Email,
	}

	userId, err := service.Create(newUser, s.DB, s.Logger)

	if err != nil {
		if err.ErrorCode == e.HttpAlreadyExist {
			return nil, status.Errorf(codes.AlreadyExists, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return &gen.CreateUserResponse{UserId: *userId}, nil
}

func (s *UserGrpcServer) ResetPassword(ctx context.Context, req *gen.ResetPasswordRequest) (*gen.ResetPasswordResponse, error) {
	if req == nil || req.UserId == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	resetPasswordRequest := service.ResetUserPasswordRequest{
		Password: req.Password,
	}

	if err := service.ResetUserPassword(resetPasswordRequest, req.UserId, s.DB, s.Logger); err != nil {
		return nil, status.Errorf(codes.Aborted, err.ErrorMessage)
	}

	return &gen.ResetPasswordResponse{UserId: req.UserId}, nil
}

func (s *UserGrpcServer) GetUserByEmail(ctx context.Context, req *gen.GetUserByEmailMsg) (*gen.User, error) {
	if req == nil || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	user, err := service.GetByEmail(req.Email, s.DB, s.Logger)

	if err != nil {
		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return mapper.UserToProto(user), nil
}

func (s *UserGrpcServer) GetUserById(ctx context.Context, req *gen.GetUserByIdMsg) (*gen.User, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	user, err := service.GetById(req.UserId, s.DB, s.Logger)

	if err != nil {
		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return mapper.UserToProto(user), nil
}

func (s *UserGrpcServer) UpdateVerificationStatus(ctx context.Context, req *gen.UpdateVerificationStatusRequest) (*gen.UpdateVerificationStatusResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	if err := service.UpdateUserVerification(req.UserId, s.DB, s.Logger); err != nil {
		if err.ErrorCode == e.HttpBadRequest {
			return nil, status.Errorf(codes.InvalidArgument, err.ErrorMessage)
		}

		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return &gen.UpdateVerificationStatusResponse{Verified: true}, nil
}