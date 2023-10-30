package postgresdb

import (
	"errors"
	"fmt"
	"reflect"
	db "warehouse/src/internal/database"

	"github.com/jackc/pgx/v5/pgconn"
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

func (cfg *PostgresDatabase[T]) Add(item *T) *db.DBError {
	err := cfg.db.Create(item).Error

	if err != nil {
		if isDuplicateKeyError(err) {
			return db.NewDBError(db.Exist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return db.NewDBError(db.System, "Something went wrong.", err.Error())
		}

	}

	return nil
}

func (cfg *PostgresDatabase[T]) GetOneBy(key string, value interface{}) (*T, *db.DBError) {
	var item T

	result := cfg.db.Where(map[string]interface{}{key: value}).First(&item)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, db.NewDBError(db.System, "Something went wrong.", result.Error.Error())
		}

		return nil, db.NewDBError(db.NotFound, "Entity not found.", result.Error.Error())
	}

	return &item, nil
}

func (cfg *PostgresDatabase[T]) Update(id string, updatedFields interface{}) (*T, *db.DBError) {
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

	if len(finalFieldsMap) == 0 {
		return nil, db.NewDBError(db.Update, "Nothing to update.", gorm.ErrEmptySlice.Error())
	}

	cfg.db.Model(&item).Clauses(clause.Returning{}).Where("id", id).Updates(finalFieldsMap)

	return &item, nil
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"

	}
	return false
}
