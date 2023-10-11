package main

import (
	"fmt"
	"os"
	"time"
	"warehouse/src/internal/database/redisdb"
	"warehouse/src/services/auth/internal/handler/http"

	"github.com/sirupsen/logrus"
)

func main() {
	// -----------SETUP LOGGER-----------
	fmt.Println("Set up the logger...")
	log := logrus.New()

	file, err := os.OpenFile("auth.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	// -----------CONNECT TO DATABASE-----------
	fmt.Println("Connect to the Session database...")
	session := redisdb.NewRedisDatabase(os.Getenv("SESSION_DB_HOST"), os.Getenv("SESSION_DB_PORT"), os.Getenv("SESSION_DB_PASSWORD"))
	fmt.Println("✅Session database successfully connected.")

	// -----------START SERVER-----------
	fmt.Println("Start the Auth Microservice...")
	api := http.NewAuthAPI(session, log)

	app := api.Init()

	if err := app.Listen(":8010"); err != nil {
		fmt.Println("❌Failed to start the Auth Microservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
		panic(err)
	}

	fmt.Println("✅Auth Microservice successfully started.")
}
