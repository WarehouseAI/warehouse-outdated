package dataservice

import (
	"context"
	m "warehouseai/auth/model"
	e "warehouseai/internal/errors"
)

type ResetTokenInterface interface {
	Create(newResetToken *m.ResetToken) *e.DBError
	Get(condition map[string]interface{}) (*m.ResetToken, *e.DBError)
	Delete(condition map[string]interface{}) *e.DBError
}

type SessionInterface interface {
	Create(ctx context.Context, userId string) (*m.Session, *e.DBError)
	Get(ctx context.Context, sessionId string) (*m.Session, *e.DBError)
	Delete(ctx context.Context, sessionId string) *e.DBError
	Update(ctx context.Context, sessionId string) (*m.Session, *e.DBError)
}
