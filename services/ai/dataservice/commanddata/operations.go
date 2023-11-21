package commanddata

import (
	"errors"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) Create(token *m.Command) *e.DBError {
	if err := d.DB.Create(token).Error; err != nil {
		if isDuplicateKeyError(err) {
			return e.NewDBError(e.DbExist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}
	}

	return nil
}

func (d *Database) Get(conditions map[string]interface{}) (*m.Command, *e.DBError) {
	var cmd m.Command

	if err := d.DB.Where(conditions).First(&cmd).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", err.Error())
	}

	return &cmd, nil
}

func (d *Database) GetWithPreload(conditions map[string]interface{}, preload string) (*m.Command, *e.DBError) {
	var cmd m.Command

	result := d.DB.Where(conditions).Preload(preload).First(&cmd)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &cmd, nil
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"

	}
	return false
}
