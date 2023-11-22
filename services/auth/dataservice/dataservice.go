package dataservice

import (
	"context"
	"mime/multipart"
	e "warehouseai/auth/errors"
	m "warehouseai/auth/model"
)

type ResetTokenInterface interface {
	Create(newResetToken *m.ResetToken) *e.DBError
	Get(condition map[string]interface{}) (*m.ResetToken, *e.DBError)
	Delete(condition map[string]interface{}) *e.DBError
}

type VerificationTokenInterface interface {
	Create(newVerificationToken *m.VerificationToken) *e.DBError
	Get(condition map[string]interface{}) (*m.VerificationToken, *e.DBError)
	Delete(condition map[string]interface{}) *e.DBError
}

type SessionInterface interface {
	Create(ctx context.Context, userId string) (*m.Session, *e.DBError)
	Get(ctx context.Context, sessionId string) (*m.Session, *e.DBError)
	Delete(ctx context.Context, sessionId string) *e.DBError
	Update(ctx context.Context, sessionId string) (*string, *m.Session, *e.DBError)
}

type PictureInterface interface {
	UploadFile(file multipart.File, fileName string) (string, error)
	DeleteImage(fileName string) error
}
