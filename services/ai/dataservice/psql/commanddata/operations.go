package commanddata

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
			return e.NewDBError(e.DbExist, "Command with this key/keys already exists.", err.Error())

		case "23503":
			return e.NewDBError(e.DbNotFound, "Invalid ai_id value", err.Error())

		case "20000":
			return e.NewDBError(e.DbNotFound, "Command not found", err.Error())
		}
	}

	return e.NewDBError(e.DbSystem, "Something went wrong", err.Error())
}

func (d *Database) Create(token *m.AiCommand) *e.DBError {
	if err := d.DB.Create(token).Error; err != nil {
		return d.errorHandle(err)
	}

	return nil
}

func (d *Database) Get(conditions map[string]interface{}) (*m.AiCommand, *e.DBError) {
	var cmd m.AiCommand

	if err := d.DB.Where(conditions).First(&cmd).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &cmd, nil
}

func (d *Database) GetWithPreload(conditions map[string]interface{}, preload string) (*m.AiCommand, *e.DBError) {
	var cmd m.AiCommand

	if err := d.DB.Where(conditions).Preload(preload).First(&cmd).Error; err != nil {
		return nil, d.errorHandle(err)
	}

	return &cmd, nil
}
