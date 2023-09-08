package gateway

import (
	"context"
	"warehouse/gen"
	utils "warehouse/src/internal/utils/grpc"
)

func CreateUser(ctx context.Context, userInfo *gen.CreateUserRequest) (*string, error) {
	conn, err := utils.ServiceConnection(ctx, "localhost:8001")

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := gen.NewUserServiceClient(conn)
	resp, err := client.CreateUser(ctx, userInfo)

	if err != nil {
		return nil, err
	}

	return &resp.UserId, nil
}
