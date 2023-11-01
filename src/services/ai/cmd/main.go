package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"
	mw "warehouse/src/internal/middleware"
	pvtAPI "warehouse/src/services/ai/internal/handler/grpc"
	pubApi "warehouse/src/services/ai/internal/handler/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	aiDatabase, err := pg.NewPostgresDatabase[pg.AI](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_NAME"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}

	if err := aiDatabase.GrantPrivileges("ais", os.Getenv("DATA_DB_USER")); err != nil {
		panic(err)
	}
	fmt.Println("✅AI database successfully connected.")

	commandDatabase, err := pg.NewPostgresDatabase[pg.Command](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_NAME"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}

	if err := commandDatabase.GrantPrivileges("commands", os.Getenv("DATA_DB_USER")); err != nil {
		panic(err)
	}
	fmt.Println("✅Command database successfully connected.")

	// -----------START SERVER-----------
	sessionMiddleware := mw.Session(log)
	userMiddleware := mw.User(log)
	pvtApi := pvtAPI.NewAiPrivateAPI(aiDatabase, log)
	pubApi := pubApi.NewAiAPI(sessionMiddleware, userMiddleware, aiDatabase, commandDatabase, log)

	publicApp := pubApi.Init()
	privateApp := grpc.NewServer()

	lis, err := net.Listen("tcp", "ai-service:8021")

	go func() {
		gen.RegisterAiServiceServer(privateApp, pvtApi)
		if err := privateApp.Serve(lis); err != nil {
			fmt.Println("❌Failed to start the Private API.")
			log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AI Microservice")
			panic(err)
		}
	}()

	if err := publicApp.Listen(":8020"); err != nil {
		fmt.Println("❌Failed to start the AI Microservice.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("AI Microservice")
		panic(err)
	}

	fmt.Println("✅AI Microservice successfully started.")
}
