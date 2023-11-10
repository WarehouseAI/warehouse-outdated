package dataservice

import (
	"errors"
	errs "warehouseai/user/errors"
	"warehouseai/user/model"

	"gorm.io/gorm"
)

func (pvd *Database) GetOneBy(conditions map[string]interface{}) (*model.User, *errs.DBError) {
	var user model.User

	result := pvd.DB.Where(conditions).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewDBError(errs.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, errs.NewDBError(errs.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &user, nil
}

func (pvd *Database) GetOneByPreload(conditions map[string]interface{}, preload string) (*model.User, *errs.DBError) {
	var user model.User

	result := pvd.DB.Where(conditions).Preload(preload).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewDBError(errs.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, errs.NewDBError(errs.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &user, nil
}
