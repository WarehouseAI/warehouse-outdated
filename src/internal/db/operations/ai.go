package operations

import (
	dbm "warehouse/src/internal/db/models"

	"gorm.io/gorm"
)

type AIDatabaseOperations interface {
	Add(*dbm.AI) error
}

type AIOperationsConfig struct {
	db *gorm.DB
}

func NewAIOperations(db *gorm.DB) AIDatabaseOperations {
	return &AIOperationsConfig{
		db: db,
	}
}

func (cfg *AIOperationsConfig) Add(ai *dbm.AI) error {
	result := cfg.db.Create(&ai)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
