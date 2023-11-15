package grpc

import (
	"fmt"
	"net"
	"time"
	"warehouseai/auth/adapter/grpc/gen"
	"warehouseai/auth/adapter/grpc/server"
	"warehouseai/auth/dataservice"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(host string, db dataservice.SessionInterface, logger *logrus.Logger) func() {
	grpc := grpc.NewServer()
	server := server.NewAuthGrpcServer(db, logger)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Println("❌Failed to listen the GRPC host.")
		logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
		panic(err)
	}

	return func() {
		gen.RegisterAuthServiceServer(grpc, server)

		if err := grpc.Serve(listener); err != nil {
			fmt.Println("❌Failed to start the GRPC server.")
			logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
			panic(err)
		}
	}
}
