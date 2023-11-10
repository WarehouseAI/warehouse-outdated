package handler

import (
	"warehouseai/user/dataservice"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	db     *dataservice.Database
	logger *logrus.Logger
}
