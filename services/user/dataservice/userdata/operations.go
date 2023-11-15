package userdata

import (
	"errors"
	"reflect"
	e "warehouseai/user/errors"
	m "warehouseai/user/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	DB *gorm.DB
}

func (pvd *Database) Add(user *m.User) *e.DBError {
	if err := pvd.DB.Create(user).Error; err != nil {
		if isDuplicateKeyError(err) {
			return e.NewDBError(e.DbExist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}
	}

	return nil
}

func (pvd *Database) GetOneBy(conditions map[string]interface{}) (*m.User, *e.DBError) {
	var user m.User

	result := pvd.DB.Where(conditions).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &user, nil
}

func (pvd *Database) GetOneByPreload(conditions map[string]interface{}, preload string) (*m.User, *e.DBError) {
	var user m.User

	result := pvd.DB.Where(conditions).Preload(preload).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", result.Error.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", result.Error.Error())
	}

	return &user, nil
}

func (pvd *Database) RawUpdate(userId string, updatedFields interface{}) (*m.User, *e.DBError) {
	var item m.User

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
		return nil, e.NewDBError(e.DbUpdate, "Nothing to update.", gorm.ErrEmptySlice.Error())
	}

	pvd.DB.Model(&item).Clauses(clause.Returning{}).Where(map[string]interface{}{"id": userId}).Updates(finalFieldsMap)

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
