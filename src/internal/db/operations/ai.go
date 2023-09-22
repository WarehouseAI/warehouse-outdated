package operations

import (
	"errors"
	dbm "warehouse/src/internal/db/models"

	"gorm.io/gorm"
)

type AIDatabaseOperations[T dbm.All] interface {
	Add(T) error
	GetOneBy(key string, value interface{}) (*T, error)
}

type AIOperationsConfig[T dbm.All] struct {
	db *gorm.DB
}

func NewAIOperations[T dbm.All](db *gorm.DB) AIDatabaseOperations[T] {
	return &AIOperationsConfig[T]{
		db: db,
	}
}

func (cfg *AIOperationsConfig[T]) Add(item T) error {
	result := cfg.db.Create(item)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (cfg *AIOperationsConfig[T]) GetOneBy(key string, value interface{}) (*T, error) {
	var item T

	result := cfg.db.Where(map[string]interface{}{key: value}).First(&item)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}

		return nil, nil
	}

	return &item, nil
}
