package main

import (
	"fmt"
	"os"
	"time"
	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"
	mw "warehouse/src/internal/middleware"
	"warehouse/src/services/ai/internal/handler/http"
	"warehouse/src/services/ai/internal/service"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// -----------SETUP LOGGER-----------
	fmt.Println("Set up the logger...")
	log := logrus.New()

	file, err := os.OpenFile("ai.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	// -----------CONNECT TO DATABASE-----------
	fmt.Println("Connect to the User database...")
	pgDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_AI"), os.Getenv("DATA_DB_PORT"))
	psqlClient, err := gorm.Open(postgres.Open(pgDSN), &gorm.Config{})
	if err != nil {
		fmt.Println(pgDSN)
		fmt.Println("❌Failed to set up the database.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Database")
		panic(err)
	}
	psqlClient.Exec("CREATE TYPE authscheme AS ENUM ('Bearer', 'Basic','ApiKey');")
	psqlClient.AutoMigrate(&dbm.AI{})

	fmt.Println("✅User database successfully connected.")

	fmt.Println("Connect to the Session database...")
	rDSN := fmt.Sprintf("%s:%s", os.Getenv("SESSION_DB_HOST"), os.Getenv("SESSION_DB_PORT"))
	fmt.Println(rDSN)
	rClient := redis.NewClient(&redis.Options{
		Addr:     pgDSN,
		Password: os.Getenv("SESSION_DB_PASSWORD"),
		DB:       0,
	})
	fmt.Println("✅Session database successfully connected.")

	// -----------START SERVER-----------
	fmt.Println("Start the AI Microservice...")
	aiOperations := dbo.NewAIOperations(psqlClient)
	userOperations := dbo.NewUserOperations(psqlClient)
	sessionOperations := dbo.NewSessionOperations(rClient)

	svc := service.NewAIService(aiOperations, log)
	sMw := mw.NewSessionMiddleware(sessionOperations, log)
	uMw := mw.NewUserMiddleware(userOperations, log)
	api := http.NewAiAPI(svc, sMw, uMw)

	app := api.Init()

	if err := app.Listen(":8020"); err != nil {
		fmt.Println("❌Failed to start the AI Microservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AI Microservice")
		panic(err)
	}

	fmt.Println("✅AI Microservice successfully started.")
}
