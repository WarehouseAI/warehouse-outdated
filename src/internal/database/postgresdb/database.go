package postgresdb

import (
	"errors"
	"fmt"
	"reflect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresDatabase[T All] struct {
	db *gorm.DB
}

func NewPostgresDatabase[T All](host string, user string, password string, dbName string, port string) (*PostgresDatabase[T], error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("❌Failed to connect to the AI database.")
		return nil, err
	}

	var structure T

	db.AutoMigrate(&structure)

	return &PostgresDatabase[T]{
		db: db,
	}, nil
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

		return nil, gorm.ErrRecordNotFound
	}

	return &item, nil
}

func (cfg *PostgresDatabase[T]) Update(id string, updatedFields interface{}) (*T, error) {
	var item T

	updatedFieldsReflect := reflect.ValueOf(updatedFields)
	itemReflect := reflect.ValueOf(item)

	finalFieldsMap := make(map[string]interface{})

	for i := 0; i < updatedFieldsReflect.NumField(); i++ {
		field := updatedFieldsReflect.Type().Field(i).Name
		value := updatedFieldsReflect.Field(i).Interface()

		genericField, exist := itemReflect.Type().FieldByName(field)

		// TODO: работает только со строками, добавить поддержку других типов
		if exist {
			if value != "" {
				finalFieldsMap[genericField.Name] = value
			}
		}
	}

	cfg.db.Model(&item).Clauses(clause.Returning{}).Where("id", id).Updates(finalFieldsMap)

	return &item, nil
}
