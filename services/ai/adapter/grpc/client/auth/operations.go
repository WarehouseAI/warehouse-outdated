package auth

import (
	"context"
	"warehouseai/ai/adapter/grpc/gen"
	e "warehouseai/ai/errors"

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

func (c *AuthGrpcClient) Authenticate(sessionId string) (string, string, *e.HttpErrorResponse) {
	client := gen.NewAuthServiceClient(c.conn)
	resp, err := client.Authenticate(context.Background(), &gen.AuthenticationRequest{SessionId: sessionId})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.Aborted {
			return "", "", e.NewErrorResponse(e.HttpUnauthorized, s.Message())
		}

		return "", "", e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp.UserId, resp.SessionId, nil
}
