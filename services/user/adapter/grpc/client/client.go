package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	TargetUrl string
}

func NewGrpcClient(grpcUrl string) *GrpcClient {
	return &GrpcClient{
		TargetUrl: grpcUrl,
	}
}

func (c *GrpcClient) Connect() (*grpc.ClientConn, error) {
	return grpc.Dial(c.TargetUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
