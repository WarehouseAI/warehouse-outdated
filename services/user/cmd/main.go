package cmd

import (
	"fmt"
	"os"
	"time"
	"warehouseai/user/adapter/grpc/server"
	"warehouseai/user/config"
	"warehouseai/user/dataservice"
	"warehouseai/user/handler"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	file, err := os.OpenFile("../user.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	db := dataservice.InitDatabase(config.NewDatabaseCfg())
	fmt.Println("✅Database successfully connected.")

	grpcProducer := server.NewUserGrpcServer(db, log)
	go grpcProducer.Start("user-service:8001")

	httpHandler := handler.NewHttpHandler(db, log)
	if err := httpHandler.Start(":8000"); err != nil {
		fmt.Println("❌Failed to start the HTTP Handler.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}
}
