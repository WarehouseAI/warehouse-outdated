package grpcutils

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {
	return grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
