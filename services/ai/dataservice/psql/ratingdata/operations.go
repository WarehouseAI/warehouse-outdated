package ratingdata

import (
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
			return e.NewDBError(e.DbExist, "Rate with this key/keys already exists.", err.Error())

		case "23503":
			return e.NewDBError(e.DbNotFound, "Invalid ai_id value", err.Error())

		case "20000":
			return e.NewDBError(e.DbNotFound, "Rate not found", err.Error())
		}
	}

	return e.NewDBError(e.DbSystem, "Something went wrong", err.Error())
}

func (d *Database) Add(rate *m.AiRate) *e.DBError {
	if err := d.DB.Create(rate).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}

func (d *Database) Get(conditions map[string]interface{}) (*m.AiRate, *e.DBError) {
	var rate m.AiRate

	if err := d.DB.Where(conditions).First(&rate).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &rate, nil
}

func (d *Database) GetAverageAiRating(aiId string) (*float64, *e.DBError) {
	var result float64

	if err := d.DB.Model(&m.AiRate{}).Select("COALESCE(AVG(rate), 0)").Where("ai_id = ?", aiId).Scan(&result).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &result, nil
}

func (d *Database) GetCountAiRating(aiId string) (*int64, *e.DBError) {
	var result int64

	if err := d.DB.Model(&m.AiRate{}).Where("ai_id = ?", aiId).Distinct("user_id").Count(&result).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &result, nil
}

func (d *Database) Update(existRate *m.AiRate, newRate int16) *e.DBError {
	if err := d.DB.Model(existRate).Updates(map[string]interface{}{"rate": newRate}).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}
