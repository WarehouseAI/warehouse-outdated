package server

import (
	"context"
	e "warehouseai/internal/errors"
	"warehouseai/internal/gen"
	"warehouseai/internal/utils"
	"warehouseai/user/dataservice"
	m "warehouseai/user/model"
	"warehouseai/user/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGrpcServer struct {
	gen.UnimplementedUserServiceServer
	db     dataservice.UserInterface
	logger *logrus.Logger
}

func NewUserGrpcServer(database dataservice.UserInterface, logger *logrus.Logger) *UserGrpcServer {
	return &UserGrpcServer{
		db:     database,
		logger: logger,
	}
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

	userId, err := service.Create(newUser, s.db, s.logger, ctx)

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

	if err := service.ResetUserPassword(resetPasswordRequest, req.UserId, s.db, s.logger); err != nil {
		return nil, status.Errorf(codes.Aborted, err.ErrorMessage)
	}

	return &gen.ResetPasswordResponse{UserId: req.UserId}, nil
}

func (s *UserGrpcServer) GetUserByEmail(ctx context.Context, req *gen.GetUserByEmailMsg) (*gen.User, error) {
	if req == nil || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	user, err := service.GetByEmail(req.Email, s.db, s.logger, ctx)

	if err != nil {
		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return utils.UserToProto(user), nil
}

func (s *UserGrpcServer) GetUserById(ctx context.Context, req *gen.GetUserByIdMsg) (*gen.User, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	user, err := service.GetById(req.UserId, s.db, s.logger, ctx)

	if err != nil {
		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return utils.UserToProto(user), nil
}
