package server

import (
	"fmt"
	"net"
	"time"
	"warehouseai/user/adapter/grpc/gen"
	"warehouseai/user/dataservice"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UserGrpcServer struct {
	gen.UnimplementedUserServiceServer
	db     *dataservice.Database
	logger *logrus.Logger
}

func NewUserGrpcServer(database *dataservice.Database, logger *logrus.Logger) *UserGrpcServer {
	return &UserGrpcServer{
		db:     database,
		logger: logger,
	}
}

func (p *UserGrpcServer) Start(host string) func() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Println("❌Failed to listen the GRPC host.")
		p.logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}

	return func() {
		gen.RegisterUserServiceServer(grpcServer, p)

		if err := grpcServer.Serve(listener); err != nil {
			fmt.Println("❌Failed to start the GRPC server.")
			p.logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
			panic(err)
		}
	}
}
