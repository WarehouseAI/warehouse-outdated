package receiver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcReceiver struct {
	TargetUrl string
}

func NewGrpcReceiver(grpcUrl string) *GrpcReceiver {
	return &GrpcReceiver{
		TargetUrl: grpcUrl,
	}
}

func (c *GrpcReceiver) Connect() (*grpc.ClientConn, error) {
	return grpc.Dial(c.TargetUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
