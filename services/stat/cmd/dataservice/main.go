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

	db.AutoMigrate(&m.User{}, &m.AI{})

	//TODO:Переместить эту ересь в sh-файлы

	db.Exec("ALTER SYSTEM SET wal_level = logical")

	userDbCfg := config.UserDatabaseCfg()
	rawSqlUserSub := fmt.Sprintf("CREATE SUBSCRIPTION users_sub \n CONNECTION 'port=%s user=%s dbname=%s host=%s password=%s' \n PUBLICATION users_pub", userDbCfg.Port, userDbCfg.User, userDbCfg.Name, userDbCfg.Host, userDbCfg.Password)
	db.Exec(rawSqlUserSub)

	aiDbCfg := config.AiDatabaseCfg()
	rawSqlAiSub := fmt.Sprintf("CREATE SUBSCRIPTION ais_sub \n CONNECTION 'port=%s user=%s dbname=%s host=%s password=%s' \n PUBLICATION ais_pub", aiDbCfg.Port, aiDbCfg.User, aiDbCfg.Name, aiDbCfg.Host, aiDbCfg.Password)
	db.Exec(rawSqlAiSub)

	return &statdata.Database{DB: db}
}
