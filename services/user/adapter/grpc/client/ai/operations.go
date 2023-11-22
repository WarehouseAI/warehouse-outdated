package ai

import (
	"context"
	"warehouseai/user/adapter/grpc/gen"
	e "warehouseai/user/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AiGrpcClient struct {
	conn *grpc.ClientConn
}

func NewAiGrpcClient(grpcUrl string) *AiGrpcClient {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return &AiGrpcClient{
		conn: conn,
	}
}

func (c *AiGrpcClient) GetById(aiId string) (string, *e.ErrorResponse) {
	client := gen.NewAiServiceClient(c.conn)
	resp, err := client.GetAiById(context.Background(), &gen.GetAiByIdMsg{Id: aiId})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return "", e.NewErrorResponse(e.HttpNotFound, s.Message())
		}

		return "", e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp.Id, nil
}
