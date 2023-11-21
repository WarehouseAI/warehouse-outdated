package grpc

import (
	"fmt"
	"net"
	"time"
	"warehouseai/ai/adapter/grpc/gen"
	"warehouseai/ai/adapter/grpc/server"
	"warehouseai/ai/dataservice"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(host string, db dataservice.AiInterface, logger *logrus.Logger) func() {
	grpc := grpc.NewServer()
	server := newAiGrpcServer(db, logger)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Println("❌Failed to listen the GRPC host.")
		logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}

	return func() {
		gen.RegisterAiServiceServer(grpc, server)

		if err := grpc.Serve(listener); err != nil {
			fmt.Println("❌Failed to start the GRPC server.")
			logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
			panic(err)
		}
	}
}

func newAiGrpcServer(database dataservice.AiInterface, logger *logrus.Logger) *server.AiGrpcServer {
	return &server.AiGrpcServer{
		DB:     database,
		Logger: logger,
	}
}
