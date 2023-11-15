package auth

import (
	"context"
	"warehouseai/user/adapter/grpc/gen"
	e "warehouseai/user/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthGrpcClient struct {
	conn *grpc.ClientConn
}

func NewAuthGrpcClient(grpcUrl string) *AuthGrpcClient {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	return &AuthGrpcClient{
		conn: conn,
	}
}

func (c *AuthGrpcClient) Authenticate(sessionId string) (*string, *string, *e.ErrorResponse) {
	client := gen.NewAuthServiceClient(c.conn)
	resp, err := client.Authenticate(context.Background(), &gen.AuthenticationRequest{SessionId: sessionId})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, nil, e.NewErrorResponse(e.HttpNotFound, s.Message())
		}

		return nil, nil, e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return &resp.UserId, &resp.SessionId, nil
}
