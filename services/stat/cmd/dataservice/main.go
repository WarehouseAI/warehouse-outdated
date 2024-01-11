package dataservice

import (
	"fmt"
	"warehouseai/stat/config"
	"warehouseai/stat/dataservice/statdata"

	m "warehouseai/stat/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewStatDatabase() *statdata.Database {
	cfg := config.NewStatDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&m.User{}, &m.AI{}, &m.Command{}, &m.RatingPerUser{}, &m.UserFavorite{})

	return &statdata.Database{DB: db}
}
