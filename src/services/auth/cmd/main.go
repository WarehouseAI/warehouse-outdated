package main

import (
	"fmt"
	"os"
	"time"
	"warehouse/src/internal/database/redisdb"
	mv "warehouse/src/internal/middleware"
	"warehouse/src/internal/s3"
	"warehouse/src/services/auth/internal/handler/http"

	"github.com/sirupsen/logrus"
)

func main() {
	// -----------SETUP LOGGER-----------
	log := logrus.New()

	file, err := os.OpenFile("auth.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	// -----------CONNECT TO DATABASE-----------
	session := redisdb.NewRedisDatabase(os.Getenv("SESSION_DB_HOST"), os.Getenv("SESSION_DB_PORT"), os.Getenv("SESSION_DB_PASSWORD"))
	fmt.Println("✅Session database successfully connected.")

	s3, err := s3.NewS3(os.Getenv("S3_LINK"), os.Getenv("S3_BUCKET"))
	if err != nil {
		panic(err)
	}

	// -----------START SERVER-----------
	sessionMiddleware := mv.Session(log)
	api := http.NewAuthAPI(session, sessionMiddleware, log, s3)

	app := api.Init()

	if err := app.Listen(":8010"); err != nil {
		fmt.Println("❌Failed to start the Auth Microservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
		panic(err)
	}

	fmt.Println("✅Microservice started.")
}
