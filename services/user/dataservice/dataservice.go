package dataservice

import (
	e "warehouseai/internal/errors"
	m "warehouseai/user/model"
)

type UserInterface interface {
	Add(user *m.User) *e.DBError
	RawUpdate(userId string, updatedFields interface{}) (*m.User, *e.DBError)
	GetOneByPreload(conditions map[string]interface{}, preload string) (*m.User, *e.DBError)
	GetOneBy(conditions map[string]interface{}) (*m.User, *e.DBError)
}
