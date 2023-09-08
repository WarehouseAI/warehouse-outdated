package datastore

import (
	"errors"

	"gorm.io/gorm"
)

type Operations interface {
	Add(user *User) error
	GetOneBy(key string, value interface{}) (*User, error)
}

type OperationsConfig struct {
	db *gorm.DB
}

func NewOperations(db *gorm.DB) Operations {
	return &OperationsConfig{
		db: db,
	}
}

func (cfg *OperationsConfig) Add(user *User) error {
	result := cfg.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (cfg *OperationsConfig) GetOneBy(key string, value interface{}) (*User, error) {
	var user User

	result := cfg.db.Where(map[string]interface{}{key: value}).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}

		return nil, nil
	}

	return &user, nil
}
