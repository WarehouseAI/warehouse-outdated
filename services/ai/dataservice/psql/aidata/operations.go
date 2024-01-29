package aidata

import (
	"fmt"
	"strings"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) errorHandle(err error) *e.DBError {
	if err == nil {
		return nil
	}

	// Добавлять новые ошибки в этот свитч и использовать потом внутри if с ошибкой
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		switch pgErr.Code {
		case "23505":
			return e.NewDBError(e.DbExist, "AI with this key/keys already exists.", err.Error())

		case "20000":
			return e.NewDBError(e.DbNotFound, "AI not found", err.Error())

		case "42703":
			return e.NewDBError(e.DbNotFound, "Invalid AI field property", err.Error())
		}
	}

	return e.NewDBError(e.DbSystem, "Something went wrong", err.Error())
}

func (d *Database) Create(ai *m.AiProduct) *e.DBError {
	if err := d.DB.Create(ai).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}

func (d *Database) Get(conditions map[string]interface{}) (*m.AiProduct, *e.DBError) {
	var ai m.AiProduct

	if err := d.DB.Where(conditions).First(&ai).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &ai, nil
}

func (d *Database) GetWithPreload(conditions map[string]interface{}, preload string) (*m.AiProduct, *e.DBError) {
	var ai m.AiProduct

	if err := d.DB.Where(conditions).Preload(preload).First(&ai).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &ai, nil
}

func (d *Database) GetMany(ids []string) (*[]m.AiProduct, *e.DBError) {
	var ais []m.AiProduct

	if err := d.DB.Where("id IN ?", ids).Preload("Commands").Find(&ais).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	if len(ais) == 0 {
		return nil, e.NewDBError(e.DbNotFound, "AIs not found.", "Empty favorites")
	}

	return &ais, nil
}

func (d *Database) GetLike(field string, value string) (*[]m.AiProduct, *e.DBError) {
	var ais []m.AiProduct

	if err := d.DB.Where(fmt.Sprintf("LOWER(%s) LIKE ?", field), strings.ToLower(value)).Preload("Commands").Find(&ais).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &ais, nil
}

func (d *Database) Update(ai *m.AiProduct, updatedFields map[string]interface{}) *e.DBError {
	if err := d.DB.Model(ai).Updates(updatedFields).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}
