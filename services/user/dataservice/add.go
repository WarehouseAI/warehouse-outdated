package dataservice

import (
	errs "warehouseai/user/errors"
	"warehouseai/user/model"

	"github.com/jackc/pgx/v5/pgconn"
)

func (pvd *Database) Add(user *model.User) *errs.DBError {
	if err := pvd.DB.Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return errs.NewDBError(errs.DbExist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return errs.NewDBError(errs.DbSystem, "Something went wrong.", err.Error())
		}
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
