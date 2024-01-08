package dataservice

import (
	"mime/multipart"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"
)

type AiInterface interface {
	Create(token *m.AiProduct) *e.DBError
	Get(conditions map[string]interface{}) (*m.AiProduct, *e.DBError)
	GetMany(ids []string) (*[]m.AiProduct, *e.DBError)
	GetLike(field string, value string) (*[]m.AiProduct, *e.DBError)
	GetWithPreload(conditions map[string]interface{}, preload string) (*m.AiProduct, *e.DBError)
	Update(ai *m.AiProduct, updatedFields map[string]interface{}) *e.DBError
}

type CommandInterface interface {
	Create(token *m.AiCommand) *e.DBError
	Get(conditions map[string]interface{}) (*m.AiCommand, *e.DBError)
	GetWithPreload(conditions map[string]interface{}, preload string) (*m.AiCommand, *e.DBError)
}

type PictureInterface interface {
	UploadFile(file multipart.File, fileName string) (string, error)
	DeleteImage(fileName string) error
}

type RatingInterface interface {
	Update(existRate *m.AiRate, newRate int16) *e.DBError
	GetAverageAiRating(aiId string) (*float64, *e.DBError)
	GetCountAiRating(aiId string) (*int64, *e.DBError)
	Get(conditions map[string]interface{}) (*m.AiRate, *e.DBError)
	Add(rate *m.AiRate) *e.DBError
}
