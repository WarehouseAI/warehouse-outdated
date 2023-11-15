package tokendata

import (
	"errors"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) Create(token *m.ResetToken) *e.DBError {
	if err := d.DB.Create(token).Error; err != nil {
		if isDuplicateKeyError(err) {
			return e.NewDBError(e.DbExist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}
	}

	return nil
}

func (d *Database) Get(conditions map[string]interface{}) (*m.ResetToken, *e.DBError) {
	var user m.ResetToken

	result := d.DB.Where(conditions).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &user, nil
}

func (d *Database) Delete(condition map[string]interface{}) *e.DBError {
	var item m.ResetToken

	if err := d.DB.Where(condition).Delete(&item).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}

		return e.NewDBError(e.DbNotFound, "Entity not found.", err.Error())
	}

	return nil
}
func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"

	}
	return false
}
