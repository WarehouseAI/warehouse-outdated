package user

import (
	"context"
	"warehouseai/ai/adapter/grpc/gen"
	e "warehouseai/ai/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type UserGrpcClient struct {
	conn *grpc.ClientConn
}

func NewUserGrpcClient(grpcUrl string) *UserGrpcClient {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return &UserGrpcClient{
		conn: conn,
	}
}

func (s *UserGrpcClient) GetFavorite(aiId string, userId string) (bool, *e.ErrorResponse) {
	client := gen.NewUserServiceClient(s.conn)

	if _, err := client.GetFavorite(context.Background(), &gen.GetFavoriteRequest{UserId: userId, AiId: aiId}); err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return false, nil
		}

		return false, e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return true, nil
}
