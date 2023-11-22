package favoritesdata

import (
	"errors"
	e "warehouseai/user/errors"
	m "warehouseai/user/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) Add(favorite *m.UserFavorites) *e.DBError {
	if err := d.DB.Create(favorite).Error; err != nil {
		if isDuplicateKeyError(err) {
			return e.NewDBError(e.DbExist, "Entity with this key/keys already exists.", err.Error())
		} else {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}
	}

	return nil
}

func (d *Database) GetUserFavorites(userId string) (*[]m.UserFavorites, *e.DBError) {
	var favorites []m.UserFavorites

	if err := d.DB.Where(map[string]interface{}{"user_id": userId}).Find(&favorites).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entities not found.", err.Error())
	}

	return &favorites, nil
}

func (d *Database) GetFavorite(userId string, aiId string) (*m.UserFavorites, *e.DBError) {
	var favorite m.UserFavorites

	if err := d.DB.Where(map[string]interface{}{"user_id": userId, "ai_id": aiId}).First(&favorite).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}

		return nil, e.NewDBError(e.DbNotFound, "Entity not found.", err.Error())
	}

	return &favorite, nil
}

func (d *Database) Delete(userId string, aiId string) *e.DBError {
	var favorite m.UserFavorites

	if err := d.DB.Where(map[string]interface{}{"user_id": userId, "ai_id": aiId}).Delete(&favorite).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewDBError(e.DbSystem, "Something went wrong.", err.Error())
		}

		return e.NewDBError(e.DbNotFound, "Entity not found.", err.Error())
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		// unique_violation = 23505
		return pgErr.Code == "23505"

	}
	return false
}
