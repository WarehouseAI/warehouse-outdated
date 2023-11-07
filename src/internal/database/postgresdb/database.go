package postgresdb

import (
	"errors"
	"fmt"
	"reflect"
	"time"
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
		fmt.Println("❌Failed to connect to the database.")
		return nil, err
	}

	var structure T

	db.AutoMigrate(&structure)

	return &PostgresDatabase[T]{
		db: db,
	}, nil
}

func (cfg *PostgresDatabase[T]) GrantPrivileges(table string, username string) error {
	return cfg.db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON TABLE %s TO %s;", table, username)).Error
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

func (cfg *PostgresDatabase[T]) GetOneBy(conditions map[string]interface{}) (*T, *db.DBError) {
	var item T

	result := cfg.db.Where(conditions).First(&item)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, db.NewDBError(db.System, "Something went wrong.", result.Error.Error())
		}

		return nil, db.NewDBError(db.NotFound, "Entity not found.", result.Error.Error())
	}

	return &item, nil
}

func (cfg *PostgresDatabase[T]) GetOneByPreload(conditions map[string]interface{}, preload string) (*T, *db.DBError) {
	var item T

	result := cfg.db.Where(conditions).Preload(preload).First(&item)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, db.NewDBError(db.System, "Something went wrong.", result.Error.Error())
		}

		return nil, db.NewDBError(db.NotFound, "Entity not found.", result.Error.Error())
	}

	return &item, nil
}

func (cfg *PostgresDatabase[T]) Transaction(transaction func(tx *gorm.DB) error) *db.DBError {
	err := cfg.db.Transaction(transaction)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return db.NewDBError(db.System, "Something went wrong.", err.Error())
		}

		return db.NewDBError(db.NotFound, "Entity not found.", err.Error())
	}

	return nil
}

func (cfg *PostgresDatabase[T]) Update(item *T, updatedFields map[string]interface{}) *db.DBError {
	if err := cfg.db.Model(item).Updates(updatedFields).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return db.NewDBError(db.System, "Something went wrong.", err.Error())
		}

		return db.NewDBError(db.NotFound, "Entity not found.", err.Error())
	}

	return nil
}

func (cfg *PostgresDatabase[T]) RawUpdate(condition map[string]interface{}, updatedFields interface{}) (*T, *db.DBError) {
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
			finalFieldsMap[genericField.Name] = value
		}
	}

	if len(finalFieldsMap) == 0 {
		return nil, db.NewDBError(db.Update, "Nothing to update.", gorm.ErrEmptySlice.Error())
	}

	cfg.db.Model(&item).Clauses(clause.Returning{}).Where(condition).Updates(finalFieldsMap)

	return &item, nil
}

func (cfg *PostgresDatabase[T]) DeleteAssociation(parent *T, deleteable interface{}, association string) *db.DBError {
	if err := cfg.db.Model(parent).Association(association).Delete(deleteable); err != nil {
		return db.NewDBError(db.System, "Unable to delete from entity", err.Error())
	}

	return nil
}

func (cfg *PostgresDatabase[T]) DeleteEntity(condition map[string]interface{}) *db.DBError {
	var item T

	if err := cfg.db.Where(condition).Delete(&item).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return db.NewDBError(db.System, "Something went wrong.", err.Error())
		}

		return db.NewDBError(db.NotFound, "Entity not found.", err.Error())
	}

	return nil
}

func (cfg *PostgresDatabase[T]) AutoDeleteSQL(duration time.Duration, query string, arg interface{}) func() {
	var item T

	ticker := time.NewTicker(duration)
	quit := make(chan struct{})

	return func() {
		for {
			select {
			case <-ticker.C:
				cfg.db.Where(query, arg).Delete(&item)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"

	}
	return false
}
