package dataservice

import (
	"fmt"
	"warehouseai/user/config"
	"warehouseai/user/dataservice/favoritesdata"
	"warehouseai/user/dataservice/userdata"
	"warehouseai/user/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewUserDatabase() *userdata.Database {
	cfg := config.NewUserDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.User{})

	return &userdata.Database{DB: db}
}

func NewFavoritesDatabase() *favoritesdata.Database {
	cfg := config.NewUserDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.UserFavorite{})

	return &favoritesdata.Database{DB: db}
}
