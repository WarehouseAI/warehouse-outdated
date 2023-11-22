package dataservice

import (
	"fmt"
	"warehouseai/ai/config"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/commanddata"
	"warehouseai/ai/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewAiDatabase() *aidata.Database {
	cfg := config.NewAiDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.AI{})

	return &aidata.Database{DB: db}
}

func NewCommandDatabase() *commanddata.Database {
	cfg := config.NewAiDatabaseCfg()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the database.")
		panic(err)
	}

	db.AutoMigrate(&model.Command{})

	return &commanddata.Database{DB: db}
}
