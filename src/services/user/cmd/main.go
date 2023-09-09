package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"warehouse/gen"
	d "warehouse/src/services/user/internal/datastore"

	pvtAPI "warehouse/src/services/user/internal/handler/grpc"
	pubAPI "warehouse/src/services/user/internal/handler/http"
	svc "warehouse/src/services/user/internal/service/user"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// -----------SETUP LOGGER-----------
	fmt.Println("Set up the logger...")
	log := logrus.New()

	file, err := os.OpenFile("user.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("❌Failed to set up the logger")
		panic(err)
	}

	log.Out = file
	fmt.Println("✅Logger successfully set up.")

	// -----------SETUP DATABASE-------------
	fmt.Println("Set up the database...")
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_NAME"), os.Getenv("POSTGRES_PORT"))
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(DSN)
		fmt.Println("❌Failed to set up the database.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Database")
		panic(err)
	}
	db.AutoMigrate(&d.User{})

	fmt.Println("✅Database successfully set up.")

	// -----------START SERVER-----------
	fmt.Println("Start the UserMicroservice...")
	operations := d.NewOperations(db)
	svc := svc.NewUserService(operations, log)
	pvtApi := pvtAPI.NewUserPrivateAPI(svc)
	pubApi := pubAPI.NewUserPublicAPI(svc)

	publicApp := pubApi.Init()
	privateApp := grpc.NewServer()

	fmt.Println("Start Private API")
	lis, err := net.Listen("tcp", "user-service:8001")
	if err != nil {
		fmt.Println("❌Failed to listen the Private API port.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("UserMicroservice")
		panic(err)
	}

	go func() {
		gen.RegisterUserServiceServer(privateApp, pvtApi)
		if err := privateApp.Serve(lis); err != nil {
			fmt.Println("❌Failed to start the Private API.")
			log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("UserMicroservice")
			panic(err)
		}
	}()

	fmt.Println("Start Public API")
	if err := publicApp.Listen(":8000"); err != nil {
		fmt.Println("❌Failed to start the Public API.")
		log.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("UserMicroservice")
		panic(err)
	}
}
