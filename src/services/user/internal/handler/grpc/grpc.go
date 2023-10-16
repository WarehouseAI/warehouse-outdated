package grpc

import (
	"context"
	"errors"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"
	"warehouse/src/internal/utils/grpcutils"
	"warehouse/src/services/user/internal/service/create"
	"warehouse/src/services/user/internal/service/get"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceProvider struct {
	gen.UnimplementedUserServiceServer
	userDatabase *pg.PostgresDatabase[pg.User]
	logger       *logrus.Logger
}

func NewUserPrivateAPI(userDatabase *pg.PostgresDatabase[pg.User], logger *logrus.Logger) *UserServiceProvider {
	return &UserServiceProvider{
		userDatabase: userDatabase,
		logger:       logger,
	}
}

func (pvd *UserServiceProvider) CreateUser(ctx context.Context, req *gen.CreateUserMsg) (*gen.CreateUserResponse, error) {
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, dto.BadRequestError.Error())
	}

	user, err := create.Create(req, pvd.userDatabase, pvd.logger, ctx)

	if err != nil && errors.Is(err, dto.ExistError) {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.CreateUserResponse{Id: grpcutils.UserToProto(user).Id}, nil
}

func (pvd *UserServiceProvider) GetUserByEmail(ctx context.Context, req *gen.GetUserByEmailMsg) (*gen.User, error) {
	if req == nil || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, dto.BadRequestError.Error())
	}

	user, err := get.GetByEmail(req, pvd.userDatabase, pvd.logger, ctx)

	if err != nil && errors.Is(err, dto.NotFoundError) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return grpcutils.UserToProto(user), nil
}

func (pvd *UserServiceProvider) GetUserById(ctx context.Context, req *gen.GetUserByIdMsg) (*gen.User, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, dto.BadRequestError.Error())
	}

	user, err := get.GetById(req, pvd.userDatabase, pvd.logger, ctx)

	if err != nil && errors.Is(err, dto.NotFoundError) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return grpcutils.UserToProto(user), nil
}
