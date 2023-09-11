package main

import (
	"fmt"
	"os"
	"time"
	d "warehouse/src/services/auth/internal/datastore"
	"warehouse/src/services/auth/internal/handler/api"
	svc "warehouse/src/services/auth/internal/service/auth"

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

	// -----------SETUP DATABASE-----------
	fmt.Println("Set up the database...")
	DSN := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	fmt.Println(DSN)
	rClient := redis.NewClient(&redis.Options{
		Addr:     DSN,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	fmt.Println("✅Database successfully set up.")

	// -----------START SERVER-----------
	fmt.Println("Start the AuthMicroservice...")
	operations := d.NewSessionOperations(rClient)
	svc := svc.NewAuthService(operations, log)
	api := api.NewAuthAPI(svc)

	app := api.Init()

	if err := app.Listen(":8010"); err != nil {
		fmt.Println("❌Failed to start the AuthMicroservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AuthMicroservice")
		panic(err)
	}

	fmt.Println("✅AuthMicroservice successfully started.")
}
