package database

import (
	"fmt"
	"warehouseai/user/dataservice"
	"warehouseai/user/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(host string, user string, password string, dbName string, port string) *dataservice.Database {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.User{})

	return &dataservice.Database{DB: db}
}
