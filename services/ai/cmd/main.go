package main

import (
	"fmt"
	"os"
	"time"
	"warehouseai/ai/cmd/adapter/grpc"
	"warehouseai/ai/cmd/dataservice"
	"warehouseai/ai/cmd/server"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	file, err := os.OpenFile("./user.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	aiDB := dataservice.NewAiDatabase()
	commandDB := dataservice.NewCommandDatabase()
	pictureStorage := dataservice.NewPictureStorage()
	fmt.Println("✅Database successfully connected.")

	grpcServer := grpc.Start("ai:8021", aiDB, log)
	go grpcServer()

	if err := server.StartServer(":8020", aiDB, commandDB, pictureStorage, log); err != nil {
		fmt.Println("❌Failed to start the HTTP Handler.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AI Microservice")
		panic(err)
	}
}
