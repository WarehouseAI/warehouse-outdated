package postgresdb

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase[T All] struct {
	db *gorm.DB
}

func NewPostgresDatabase[T All](host string, user string, password string, dbName string, port string) *PostgresDatabase[T] {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the AI database.")
		panic(err)
	}

	var structure T

	db.AutoMigrate(&structure)

	return &PostgresDatabase[T]{
		db: db,
	}
}

func (cfg *PostgresDatabase[T]) Add(item *T) error {
	result := cfg.db.Create(item)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (cfg *PostgresDatabase[T]) GetOneBy(key string, value interface{}) (*T, error) {
	var item T

	result := cfg.db.Where(map[string]interface{}{key: value}).First(&item)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}

		return nil, nil
	}

	return &item, nil
}
