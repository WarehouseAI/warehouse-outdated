package server

import (
	"context"
	"warehouseai/auth/adapter/grpc/gen"
	"warehouseai/auth/dataservice"
	e "warehouseai/auth/errors"
	"warehouseai/auth/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcServer struct {
	gen.UnimplementedAuthServiceServer
	db     dataservice.SessionInterface
	logger *logrus.Logger
}

func NewAuthGrpcServer(database dataservice.SessionInterface, logger *logrus.Logger) *AuthGrpcServer {
	return &AuthGrpcServer{
		db:     database,
		logger: logger,
	}
}

func (s *AuthGrpcServer) Authenticate(ctx context.Context, req *gen.AuthenticationRequest) (*gen.AuthenticationResponse, error) {
	if req == nil || req.SessionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty request data")
	}

	userId, session, err := service.Authenticate(req.SessionId, s.db, s.logger)

	if err != nil {
		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Error(codes.Internal, err.ErrorMessage)
	}

	return &gen.AuthenticationResponse{UserId: *userId, SessionId: session.ID}, nil
}
