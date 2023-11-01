package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"warehouse/gen"

	pg "warehouse/src/internal/database/postgresdb"
	mw "warehouse/src/internal/middleware"
	pvtAPI "warehouse/src/services/user/internal/handler/grpc"
	pubAPI "warehouse/src/services/user/internal/handler/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// -----------SETUP LOGGER-----------
	log := logrus.New()

	file, err := os.OpenFile("user.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	// -----------CONNECT TO DATABASE-------------
	userDatabase, err := pg.NewPostgresDatabase[pg.User](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_NAME"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}

	if err := userDatabase.GrantPrivileges("users", os.Getenv("DATA_DB_USER")); err != nil {
		panic(err)
	}
	fmt.Println("✅Database successfully connected.")

	// -----------START SERVER-----------
	fmt.Println("Start the User Microservice...")

	sessionMiddleware := mw.Session(log)
	userMiddleware := mw.User(log)

	pvtApi := pvtAPI.NewUserPrivateAPI(userDatabase, log)
	pubApi := pubAPI.NewUserPublicAPI(userDatabase, userMiddleware, sessionMiddleware, log)

	publicApp := pubApi.Init()
	privateApp := grpc.NewServer()

	lis, err := net.Listen("tcp", "user-service:8001")
	if err != nil {
		fmt.Println("❌Failed to listen the Private API port.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}

	go func() {
		gen.RegisterUserServiceServer(privateApp, pvtApi)
		if err := privateApp.Serve(lis); err != nil {
			fmt.Println("❌Failed to start the Private API.")
			log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
			panic(err)
		}
	}()

	if err := publicApp.Listen(":8000"); err != nil {
		fmt.Println("❌Failed to start the Public API.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("User Microservice")
		panic(err)
	}

	fmt.Println("✅Microservice started.")
}
