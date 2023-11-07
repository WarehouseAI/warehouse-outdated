package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpcutils"
	"warehouse/src/internal/utils/httputils"
	"warehouse/src/services/auth/internal/service/register"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGrpcConnection struct {
	grpcUrl string
}

func NewUserGrpcConnection(grpcUrl string) *UserGrpcConnection {
	return &UserGrpcConnection{
		grpcUrl: grpcUrl,
	}
}

func (c UserGrpcConnection) Create(ctx context.Context, userInfo *gen.CreateUserMsg) (*register.Response, *httputils.ErrorResponse) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.CreateUser(ctx, userInfo)

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.AlreadyExists {
			return nil, httputils.NewErrorResponse(httputils.AlreadyExist, s.Message())
		}

		return nil, httputils.NewErrorResponse(httputils.InternalError, s.Message())
	}

	return &register.Response{ID: resp.Id}, nil
}

func (c UserGrpcConnection) ResetPassword(ctx context.Context, updateInfo *gen.ResetPasswordRequest) (*gen.ResetPasswordResponse, *httputils.ErrorResponse) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.ResetPassword(ctx, updateInfo)

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.Aborted {
			return nil, httputils.NewErrorResponse(httputils.BadRequest, s.Message())
		}

		return nil, httputils.NewErrorResponse(httputils.InternalError, s.Message())
	}

	return resp, nil
}

func (c UserGrpcConnection) GetByEmail(ctx context.Context, userInfo *gen.GetUserByEmailMsg) (*gen.User, *httputils.ErrorResponse) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.GetUserByEmail(ctx, userInfo)

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, httputils.NewErrorResponse(httputils.NotFound, s.Message())
		}

		return nil, httputils.NewErrorResponse(httputils.InternalError, s.Message())
	}

	return resp, nil
}

func (c UserGrpcConnection) GetById(ctx context.Context, userInfo *gen.GetUserByIdMsg) (*gen.User, *httputils.ErrorResponse) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.GetUserById(ctx, userInfo)

	if err != nil {
		s, _ := status.FromError(err)

		if s.Code() == codes.NotFound {
			return nil, httputils.NewErrorResponse(httputils.NotFound, s.Message())
		}

		return nil, httputils.NewErrorResponse(httputils.InternalError, s.Message())
	}

	return resp, nil
}
