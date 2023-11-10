package dataservice

import (
	"reflect"
	errs "warehouseai/user/errors"
	"warehouseai/user/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (pvd *Database) RawUpdate(userId string, updatedFields interface{}) (*model.User, *errs.DBError) {
	var item model.User

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
		return nil, errs.NewDBError(errs.DbUpdate, "Nothing to update.", gorm.ErrEmptySlice.Error())
	}

	pvd.DB.Model(&item).Clauses(clause.Returning{}).Where(map[string]interface{}{"id": userId}).Updates(finalFieldsMap)

	return &item, nil
}
