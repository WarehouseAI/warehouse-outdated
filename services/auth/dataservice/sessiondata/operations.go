package sessiondata

import (
	"context"
	"encoding/json"
	"time"
	m "warehouseai/auth/model"
	e "warehouseai/internal/errors"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	DB *redis.Client
}

func (d *Database) Create(ctx context.Context, userId string) (*m.Session, *e.DBError) {
	TTL := 24 * time.Hour
	sessionId := uuid.Must(uuid.NewV4()).String()

	sessionPayload := m.SessionPayload{
		UserId:    userId,
		CreatedAt: time.Now(),
	}

	marshaledPayload, err := json.Marshal(sessionPayload)

	if err != nil {
		return nil, e.NewDBError(e.DbSystem, "Can't marshal entity to JSON", err.Error())
	}

	if err := d.DB.Set(ctx, sessionId, marshaledPayload, TTL).Err(); err != nil {
		return nil, e.NewDBError(e.DbSystem, "Can't save JSON in DB", err.Error())
	}

	return &m.Session{ID: sessionId, Payload: sessionPayload, TTL: TTL}, nil
}

func (d *Database) Get(ctx context.Context, sessionId string) (*m.Session, *e.DBError) {
	var sessionPayload m.SessionPayload

	record := d.DB.Get(ctx, sessionId)
	recordTTL := d.DB.TTL(ctx, sessionId)

	if record.Err() != nil {
		return nil, e.NewDBError(e.DbNotFound, "Record not found.", record.Err().Error())
	}

	recordInfo, _ := record.Result()
	TTLInfo, _ := recordTTL.Result()

	if err := json.Unmarshal([]byte(recordInfo), &sessionPayload); err != nil {
		return nil, e.NewDBError(e.DbSystem, "Can't unmarhal record", err.Error())
	}

	return &m.Session{ID: sessionId, Payload: sessionPayload, TTL: TTLInfo}, nil
}

func (d *Database) Delete(ctx context.Context, sessionId string) *e.DBError {
	if err := d.DB.Del(ctx, sessionId).Err(); err != nil {
		return e.NewDBError(e.DbNotFound, "Record not found.", err.Error())
	}

	return nil
}

func (d *Database) Update(ctx context.Context, sessionId string) (*m.Session, *e.DBError) {
	session, err := d.Get(ctx, sessionId)

	if err != nil {
		return nil, err
	}

	if err := d.Delete(ctx, sessionId); err != nil {
		return nil, err
	}

	return d.Create(ctx, session.Payload.UserId)
}
