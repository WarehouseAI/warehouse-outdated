package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpc"
	m "warehouse/src/services/auth/pkg/model"
)

func CreateUser(ctx context.Context, userInfo *gen.CreateUserRequest) (*m.UserIdResponse, error) {
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
