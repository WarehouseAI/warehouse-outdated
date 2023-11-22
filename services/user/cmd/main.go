package main

import (
	"fmt"
	"os"
	"time"
	"warehouseai/user/cmd/adapter/broker/mail"
	"warehouseai/user/cmd/adapter/grpc"
	"warehouseai/user/cmd/dataservice"
	"warehouseai/user/cmd/server"

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

	userDB := dataservice.NewUserDatabase()
	favoritesDB := dataservice.NewFavoritesDatabase()
	mailProducer := mail.NewMailProducer()
	fmt.Println("✅Database successfully connected.")

	grpcServer := grpc.Start("user:8001", userDB, favoritesDB, log)
	go grpcServer()

	if err := server.StartServer(":8000", userDB, favoritesDB, mailProducer, log); err != nil {
		fmt.Println("❌Failed to start the HTTP Handler.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}
}
