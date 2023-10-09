package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpcutils"
	m "warehouse/src/services/auth/pkg/models"
)

func CreateUser(ctx context.Context, userInfo *gen.CreateUserMsg) (*m.RegisterResponse, error) {
	conn, err := utils.ServiceConnection(ctx, "user-service:8001")

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

func GetUserByEmail(ctx context.Context, userInfo *gen.GetUserByEmailMsg) (*gen.User, error) {
	conn, err := utils.ServiceConnection(ctx, "user-service:8001")

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

func GetUserById(ctx context.Context, userInfo *gen.GetUserByIdMsg) (*gen.User, error) {
	conn, err := utils.ServiceConnection(ctx, "user-service:8001")

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
