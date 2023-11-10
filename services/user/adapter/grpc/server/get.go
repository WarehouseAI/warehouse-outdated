package server

import (
	"context"
	models "warehouseai/user/adapter/grpc"
	"warehouseai/user/adapter/grpc/gen"
	errs "warehouseai/user/errors"
	"warehouseai/user/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (pvd *UserGrpcProducer) GetUserByEmail(ctx context.Context, req *gen.GetUserByEmailMsg) (*gen.User, error) {
	if req == nil || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	user, err := service.GetByEmail(req.Email, pvd.db, pvd.logger, ctx)

	if err != nil {
		if err.ErrorCode == errs.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return models.UserToProto(user), nil
}

func (pvd *UserGrpcProducer) GetUserById(ctx context.Context, req *gen.GetUserByIdMsg) (*gen.User, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	user, err := service.GetById(req.Id, pvd.db, pvd.logger, ctx)

	if err != nil {
		if err.ErrorCode == errs.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return models.UserToProto(user), nil
}
