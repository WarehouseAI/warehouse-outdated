package user

import (
	"context"
	"fmt"
	"warehouseai/auth/adapter/grpc/gen"
	e "warehouseai/auth/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type UserGrpcClient struct {
	conn gen.UserServiceClient
}

func NewUserGrpcClient(grpcUrl string) *UserGrpcClient {
	conn, err := grpc.Dial(grpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	client := gen.NewUserServiceClient(conn)

	return &UserGrpcClient{
		conn: client,
	}
}

func (c *UserGrpcClient) Create(ctx context.Context, userInfo *gen.CreateUserMsg) (string, *e.ErrorResponse) {
	resp, err := c.conn.CreateUser(ctx, userInfo)

	if err != nil {
		fmt.Println(err)
		s, _ := status.FromError(err)

		if s.Code() == codes.AlreadyExists {
			return "", e.NewErrorResponse(e.HttpAlreadyExist, s.Message())
		}

		return "", e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp.UserId, nil
}

func (c *UserGrpcClient) ResetPassword(ctx context.Context, resetPasswordRequest *gen.ResetPasswordRequest) (string, *e.ErrorResponse) {
	resp, err := c.conn.ResetPassword(ctx, resetPasswordRequest)

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.Aborted {
			return "", e.NewErrorResponse(e.HttpBadRequest, s.Message())
		}

		return "", e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp.UserId, nil
}

func (c *UserGrpcClient) GetByEmail(ctx context.Context, email string) (*gen.User, *e.ErrorResponse) {
	resp, err := c.conn.GetUserByEmail(ctx, &gen.GetUserByEmailMsg{Email: email})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, e.NewErrorResponse(e.HttpNotFound, s.Message())
		}

		return nil, e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp, nil
}

func (c *UserGrpcClient) UpdateVerificationStatus(ctx context.Context, userId string) (bool, *e.ErrorResponse) {
	resp, err := c.conn.UpdateVerificationStatus(ctx, &gen.UpdateVerificationStatusRequest{UserId: userId})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.InvalidArgument {
			return false, e.NewErrorResponse(e.HttpBadRequest, s.Message())
		}

		if s.Code() == codes.NotFound {
			return false, e.NewErrorResponse(e.HttpNotFound, s.Message())
		}

		return false, e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp.Verified, nil
}

func (c *UserGrpcClient) GetById(ctx context.Context, userId string) (*gen.User, *e.ErrorResponse) {
	resp, err := c.conn.GetUserById(ctx, &gen.GetUserByIdMsg{UserId: userId})

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, e.NewErrorResponse(e.HttpNotFound, s.Message())
		}

		return nil, e.NewErrorResponse(e.HttpInternalError, s.Message())
	}

	return resp, nil
}
