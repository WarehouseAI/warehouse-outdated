package dataservice

import (
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
)

type AiInterface interface {
	Create(token *m.AI) *e.DBError
	Get(conditions map[string]interface{}) (*m.AI, *e.DBError)
	GetMany(ids []string) (*[]m.AI, *e.DBError)
	GetLike(field string, value string) (*[]m.AI, *e.DBError)
	GetWithPreload(conditions map[string]interface{}, preload string) (*m.AI, *e.DBError)
	Update(ai *m.AI, updatedFields map[string]interface{}) *e.DBError
}

type CommandInterface interface {
	Create(token *m.Command) *e.DBError
	Get(conditions map[string]interface{}) (*m.Command, *e.DBError)
	GetWithPreload(conditions map[string]interface{}, preload string) (*m.Command, *e.DBError)
}
