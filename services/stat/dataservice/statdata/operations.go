package statdata

import (
	e "warehouseai/stat/errors"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (pvd *Database) GetNumOfUsers() (uint, *e.DBError) {
	return 0, nil
}

func (pvd *Database) GetNumOfDevelopers() (uint, *e.DBError) {
	return 0, nil
}

func (pvd *Database) GetNumOfAiUses(id uuid.UUID) (uint, *e.DBError) {
	return 0, nil
}
