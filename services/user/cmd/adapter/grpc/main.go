package grpc

import (
	"fmt"
	"net"
	"time"
	"warehouseai/user/adapter/grpc/gen"
	"warehouseai/user/adapter/grpc/server"
	"warehouseai/user/dataservice"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(host string, userDB dataservice.UserInterface, favoriteDB dataservice.FavoritesInterface, logger *logrus.Logger) func() {
	grpc := grpc.NewServer()
	server := newUserGrpcServer(userDB, favoriteDB, logger)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Println("❌Failed to listen the GRPC host.")
		logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}

	return func() {
		gen.RegisterUserServiceServer(grpc, server)

		if err := grpc.Serve(listener); err != nil {
			fmt.Println("❌Failed to start the GRPC server.")
			logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
			panic(err)
		}
	}
}

func newUserGrpcServer(user dataservice.UserInterface, favorites dataservice.FavoritesInterface, logger *logrus.Logger) *server.UserGrpcServer {
	return &server.UserGrpcServer{
		UserDB:     user,
		FavoriteDB: favorites,
		Logger:     logger,
	}
}
