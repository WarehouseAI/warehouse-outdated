package grpc

import (
	"context"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/grpcutils"
	"warehouse/src/internal/utils/httputils"
	"warehouse/src/services/ai/internal/service/ai/get"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AiServiceProvider struct {
	gen.UnimplementedAiServiceServer
	aiDatabase *pg.PostgresDatabase[pg.AI]
	logger     *logrus.Logger
}

func NewAiPrivateAPI(aiDatabase *pg.PostgresDatabase[pg.AI], logger *logrus.Logger) *AiServiceProvider {
	return &AiServiceProvider{
		aiDatabase: aiDatabase,
		logger:     logger,
	}
}

func (pvd *AiServiceProvider) GetAiById(ctx context.Context, req *gen.GetAiByIdMsg) (*gen.AI, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty request data")
	}

	ai, err := get.GetLoadedAiByID(req, pvd.aiDatabase, pvd.logger)

	if err != nil {
		if err.ErrorCode == httputils.NotFound {
			return nil, status.Errorf(codes.NotFound, err.ErrorMessage)
		}

		return nil, status.Errorf(codes.Internal, err.ErrorMessage)
	}

	return grpcutils.AiToProto(ai), nil
}
