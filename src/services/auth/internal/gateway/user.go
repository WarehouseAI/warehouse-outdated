package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpcutils"
	m "warehouse/src/services/auth/pkg/models"
)

type UserGrpcConnection struct {
	grpcUrl string
}

func NewUserGrpcConnection(grpcUrl string) *UserGrpcConnection {
	return &UserGrpcConnection{
		grpcUrl: grpcUrl,
	}
}

func (c UserGrpcConnection) Create(ctx context.Context, userInfo *gen.CreateUserMsg) (*m.RegisterResponse, error) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.CreateUser(ctx, userInfo)

	if err != nil {
		return nil, err
	}

	return &m.RegisterResponse{ID: resp.Id}, nil
}

func (c UserGrpcConnection) GetByEmail(ctx context.Context, userInfo *gen.GetUserByEmailMsg) (*gen.User, error) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.GetUserByEmail(ctx, userInfo)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c UserGrpcConnection) GetById(ctx context.Context, userInfo *gen.GetUserByIdMsg) (*gen.User, error) {
	conn, err := utils.ServiceConnection(ctx, c.grpcUrl)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.GetUserById(ctx, userInfo)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
