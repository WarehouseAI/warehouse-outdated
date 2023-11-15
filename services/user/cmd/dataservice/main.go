package dataservice

import (
	"fmt"
	"warehouseai/user/config"
	"warehouseai/user/dataservice/userdata"
	"warehouseai/user/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewUserDatabase() *userdata.Database {
	config := config.NewUserDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Password, config.Name, config.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.User{})

	return &userdata.Database{DB: db}
}
