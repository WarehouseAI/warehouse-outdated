package server

import (
	"context"
	"warehouseai/ai/adapter/grpc/gen"
	"warehouseai/ai/adapter/grpc/mapper"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	"warehouseai/ai/service/ai"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AiGrpcServer struct {
	gen.UnimplementedAiServiceServer
	DB     dataservice.AiInterface
	Logger *logrus.Logger
}

func (s *AiGrpcServer) GetAiById(ctx context.Context, req *gen.GetAiByIdMsg) (*gen.AI, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	ai, err := ai.GetByIdPreload(req.Id, s.DB, s.Logger)

	if err != nil {
		if err.ErrorCode == e.HttpNotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return mapper.AiToProto(&ai.AiProduct), nil
}
