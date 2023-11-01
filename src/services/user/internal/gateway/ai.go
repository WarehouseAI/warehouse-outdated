package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpcutils"
	"warehouse/src/internal/utils/httputils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AiGrpcConnection struct {
	grpcUrl string
}

func NewAiGrpcConnection(grpcUrl string) *AiGrpcConnection {
	return &AiGrpcConnection{
		grpcUrl: grpcUrl,
	}
}

func (c AiGrpcConnection) GetById(ctx context.Context, aiInfo *gen.GetAiByIdMsg) (*gen.AI, *httputils.ErrorResponse) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	defer conn.Close()

	client := gen.NewAiServiceClient(conn)
	resp, err := client.GetAiById(ctx, aiInfo)

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, httputils.NewErrorResponse(httputils.NotFound, s.Message())
		}

		return nil, httputils.NewErrorResponse(httputils.InternalError, s.Message())
	}

	return resp, nil
}
