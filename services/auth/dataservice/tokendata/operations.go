package tokendata

import (
	"errors"
	"fmt"
	"time"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Database[T m.Tokens] struct {
	DB *gorm.DB
}

func (d *Database[T]) Create(token *T) *e.DBError {
	if err := d.DB.Create(token).Error; err != nil {
		if isDuplicateKeyError(err) {
			return e.NewDBError(e.DbExist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}
	}

	return nil
}

func (d *Database[T]) Get(conditions map[string]interface{}) (*T, *e.DBError) {
	var token T

	result := d.DB.Where(conditions).First(&token)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &token, nil
}

func (d *Database[T]) Delete(condition map[string]interface{}) *e.DBError {
	var token T

	if err := d.DB.Where(condition).Delete(&token).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}

		return e.NewDBError(e.DbNotFound, "Entity not found.", err.Error())
	}

	return nil
}

func (d *Database[T]) Flusher(duration time.Duration) func() {
	var items []T

	ticker := time.NewTicker(duration)
	done := make(chan bool)

	return func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := d.DB.Where("expires_at < ?", time.Now()).Delete(&items).Error
				fmt.Println(err)
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
