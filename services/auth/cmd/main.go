package main

import (
	"fmt"
	"os"
	"time"
	"warehouseai/auth/cmd/adapter/broker"
	"warehouseai/auth/cmd/adapter/grpc"
	"warehouseai/auth/cmd/dataservice"
	"warehouseai/auth/cmd/server"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	file, err := os.OpenFile("./auth.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	sessionDB := dataservice.NewSessionDatabase()
	resetTokenDB := dataservice.NewResetTokenDatabase()
	verificationTokenDB := dataservice.NewVerificationTokenDatabase()
	pictureStorage := dataservice.NewPictureStorage()
	broker := broker.NewBroker()

	fmt.Println("✅Database successfully connected.")

	grpcServer := grpc.Start("auth:8041", sessionDB, log)
	go grpcServer()

	if err := server.StartServer(":8040", resetTokenDB, verificationTokenDB, sessionDB, pictureStorage, broker, log); err != nil {
		fmt.Println("❌Failed to start the HTTP Handler.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}

	defer func() {
		broker.Channel.Close()
		broker.Connection.Close()
	}()
}
