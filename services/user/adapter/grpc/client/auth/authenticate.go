package auth

import (
	"context"
	"warehouseai/user/adapter/grpc/client"
	"warehouseai/user/adapter/grpc/gen"
	errs "warehouseai/user/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Authenticate(consumer *client.GrpcClient, sessionId string) (*string, *errs.ErrorResponse) {
	conn, err := consumer.Connect()

	if err != nil {
		return nil, errs.NewErrorResponse(errs.HttpInternalError, err.Error())
	}

	defer conn.Close()

	client := gen.NewAuthServiceClient(conn)
	resp, err := client.Authenticate(context.Background(), &gen.AuthenticationRequest{SessionId: sessionId})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, errs.NewErrorResponse(errs.HttpNotFound, s.Message())
		}

		return nil, errs.NewErrorResponse(errs.HttpInternalError, s.Message())
	}

	return &resp.UserId, nil
}
