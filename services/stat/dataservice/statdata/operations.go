package statdata

import (
	"errors"
	e "warehouseai/stat/errors"
	m "warehouseai/stat/model"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (pvd *Database) GetNumOfUsers() (uint, *e.DBError) {
	var num uint

	if err := pvd.DB.Model(m.User{}).Select("COUNT(*)").Scan(&num).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, e.NewDBError(e.DbSystem, "Something went wrong", err.Error())
		}

		return 0, e.NewDBError(e.DbNotFound, "User not found.", err.Error())
	}

	return num, nil
}

func (pvd *Database) GetNumOfDevelopers() (uint, *e.DBError) {
	var num uint

	if err := pvd.DB.Model(m.User{}).Select("COUNT(*)").Where("role = ?", "Developer").Find(&num).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, e.NewDBError(e.DbSystem, "Something went wrong", err.Error())
		}

		return 0, e.NewDBError(e.DbNotFound, "User not found.", err.Error())
	}

	return num, nil
}

func (pvd *Database) GetNumOfAiUses(id uuid.UUID) (uint, *e.DBError) {
	var ai m.AI
	if err := pvd.DB.Model(m.AI{}).Where("id = ?", id).Find(&ai).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, e.NewDBError(e.DbSystem, "Something went wrong", err.Error())
		}

		return 0, e.NewDBError(e.DbNotFound, "User not found.", err.Error())
	}

	return uint(ai.Used), nil
}

// func isDuplicateKeyError(err error) bool {
// 	pgErr, ok := err.(*pgconn.PgError)
// 	if ok {
// 		// unique_violation = 23505
// 		return pgErr.Code == "23505"

// 	}
// 	return false
// }
