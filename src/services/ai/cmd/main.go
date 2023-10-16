package main

import (
	"fmt"
	"os"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	r "warehouse/src/internal/database/redisdb"
	mw "warehouse/src/internal/middleware"
	"warehouse/src/services/ai/internal/handler/http"

	"github.com/sirupsen/logrus"
)

func main() {
	// -----------SETUP LOGGER-----------
	log := logrus.New()

	file, err := os.OpenFile("ai.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	// -----------CONNECT DATABASES-----------

	userDatabase, err := pg.NewPostgresDatabase[pg.User](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_USERS"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}
	fmt.Println("✅User database successfully connected.")

	aiDatabase, err := pg.NewPostgresDatabase[pg.AI](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_AI"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}
	fmt.Println("✅AI database successfully connected.")

	commandDatabase, err := pg.NewPostgresDatabase[pg.Command](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_AI"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}
	fmt.Println("✅Command database successfully connected.")

	sessionDatabase := r.NewRedisDatabase(os.Getenv("SESSION_DB_HOST"), os.Getenv("SESSION_DB_PORT"), os.Getenv("SESSION_DB_PASSWORD"))
	fmt.Println("✅Session database successfully connected.")

	// -----------START SERVER-----------
	sessionMiddleware := mw.Session(sessionDatabase, log)
	userMiddleware := mw.User(userDatabase, log)
	api := http.NewAiAPI(sessionMiddleware, userMiddleware, aiDatabase, commandDatabase, log)

	app := api.Init()

	if err := app.Listen(":8020"); err != nil {
		fmt.Println("❌Failed to start the AI Microservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AI Microservice")
		panic(err)
	}

	fmt.Println("✅AI Microservice successfully started.")
}
