package main

import (
	"fmt"
	"os"
	"time"
	dbo "warehouse/src/internal/db/operations"
	mv "warehouse/src/internal/middleware"
	"warehouse/src/services/auth/internal/handler/http"
	svc "warehouse/src/services/auth/internal/service"

	"github.com/redis/go-redis/v9"
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

	// // -----------SETUP .ENV-----------
	// fmt.Println("Set up the enviroment variables...")
	// if err := godotenv.Load(); err != nil {
	// 	fmt.Println("❌Failed to set up the enviroment variables.")
	// 	log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Env")
	// 	panic(err)
	// }
	// fmt.Println("✅Environment successfully set up.")

	// -----------CONNECT TO DATABASE-----------
	fmt.Println("Connect to the Session database...")
	DSN := fmt.Sprintf("%s:%s", os.Getenv("SESSION_DB_HOST"), os.Getenv("SESSION_DB_PORT"))
	fmt.Println(DSN)
	rClient := redis.NewClient(&redis.Options{
		Addr:     DSN,
		Password: os.Getenv("SESSION_DB_PASSWORD"),
		DB:       0,
	})
	fmt.Println("✅Session database successfully connected.")

	// -----------START SERVER-----------
	fmt.Println("Start the Auth Microservice...")
	operations := dbo.NewSessionOperations(rClient)
	svc := svc.NewAuthService(operations, log)
	sMw := mv.NewSessionMiddleware(operations, log)
	api := http.NewAuthAPI(svc, sMw)

	app := api.Init()

	if err := app.Listen(":8010"); err != nil {
		fmt.Println("❌Failed to start the Auth Microservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
		panic(err)
	}

	fmt.Println("✅Auth Microservice successfully started.")
}
