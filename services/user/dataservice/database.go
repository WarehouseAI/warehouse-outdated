package dataservice

import (
	"fmt"
	"warehouseai/user/config"
	"warehouseai/user/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDatabase(cfg config.DatabaseCfg) *Database {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.User{})

	return &Database{DB: db}
}
