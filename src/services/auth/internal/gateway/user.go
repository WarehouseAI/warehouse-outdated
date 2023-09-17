package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpcutils"
	m "warehouse/src/services/auth/pkg/models"
)

func CreateUser(ctx context.Context, userInfo *gen.CreateUserRequest) (*m.RegisterResponse, error) {
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

	return m.UserIdFromProto(resp), nil
}

func GetUser(ctx context.Context, userInfo *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	conn, err := utils.ServiceConnection(ctx, "user-service:8001")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.GetUser(ctx, userInfo)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
