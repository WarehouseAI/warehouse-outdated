package server

import (
	"warehouseai/user/adapter/grpc/gen"
	"warehouseai/user/dataservice"

	"github.com/sirupsen/logrus"
)

type UserGrpcProducer struct {
	gen.UnimplementedUserServiceServer
	db     *dataservice.Database
	logger *logrus.Logger
}

func NewUserGrpcProducer(database *dataservice.Database, logger *logrus.Logger) *UserGrpcProducer {
	return &UserGrpcProducer{
		db:     database,
		logger: logger,
	}
}
