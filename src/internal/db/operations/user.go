package operations

import (
	"errors"
	dbm "warehouse/src/internal/db/models"

	"gorm.io/gorm"
)

type UserDatabaseOperations interface {
	Add(user *dbm.User) error
	GetOneBy(key string, value interface{}) (*dbm.User, error)
}

type UserOperationsConfig struct {
	db *gorm.DB
}

func NewUserOperations(db *gorm.DB) UserDatabaseOperations {
	return &UserOperationsConfig{
		db: db,
	}
}

func (cfg *UserOperationsConfig) Add(user *dbm.User) error {
	result := cfg.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (cfg *UserOperationsConfig) GetOneBy(key string, value interface{}) (*dbm.User, error) {
	var user dbm.User

	result := cfg.db.Where(map[string]interface{}{key: value}).First(&user)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}

		return nil, nil
	}

	return &user, nil
}
