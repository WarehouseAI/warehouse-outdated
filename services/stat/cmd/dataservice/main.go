package dataservice

import (
	"fmt"
	"warehouseai/stat/config"
	"warehouseai/stat/dataservice/statdata"

	e "warehouseai/stat/errors"
	m "warehouseai/stat/model"

	"github.com/jackc/pgx/v5/pgconn"
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

	//TODO: подумать над тем, чтобы сделать запускать репликации по запросу от реплицируемого

	db.AutoMigrate(&m.User{}, &m.AI{}, &m.Command{}, &m.RatingPerUser{}, &m.UserFavorite{})

	if err := createSub(config.NewAiDatabaseCfg(), db); err != nil {
		panic(err)
	}

	if err := createSub(config.NewUserDatabaseCfg(), db); err != nil {
		panic(err)
	}

	return &statdata.Database{DB: db}
}

func createSub(c config.DatabaseCfg, db *gorm.DB) *e.DBError {
	rawSql := fmt.Sprintf("CREATE SUBSCRIPTION %s CONNECTION 'port=%s user=%s dbname=%s host=%s password=%s' PUBLICATION %s", c.SubName, c.Port, c.User, c.Name, c.Host, c.Password, c.PubName)

	if err := db.Exec(rawSql).Error; err != nil {
		if err.(*pgconn.PgError).Code != "42710" {
			return e.NewDBError(e.DbSystem, "Something went wrong with subscription.", err.Error())
		}
	}
	return nil
}
