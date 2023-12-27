package main

import (
	"fmt"
	"os"
	"time"
	"warehouseai/stat/cmd/dataservice"
	"warehouseai/stat/cmd/server"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	file, err := os.OpenFile("./stat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	statDB := dataservice.NewStatDatabase()
	if err := server.StartServer(":8020", statDB, log); err != nil {
		fmt.Println("❌Failed to start the HTTP Handler.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AI Microservice")
		panic(err)
	}

	fmt.Println("Start the Statistic service...")
}
