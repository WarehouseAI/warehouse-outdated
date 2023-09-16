package operations

import (
	"context"
	"encoding/json"
	"time"
	m "warehouse/src/internal/db/models"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionDatabaseOperations interface {
	CreateSession(context.Context, string) (*m.Session, error)
	GetSession(context.Context, string) (*m.Session, error)
	DeleteSession(context.Context, string) error
	UpdateSession(context.Context, string) (*m.Session, error)
}

type SessionOperationsConfig struct {
	rClient *redis.Client
}

func NewSessionOperations(rClient *redis.Client) SessionDatabaseOperations {
	return &SessionOperationsConfig{
		rClient: rClient,
	}
}

func (cfg *SessionOperationsConfig) CreateSession(ctx context.Context, userId string) (*m.Session, error) {
	//TODO: Поменять потом на 3 дня
	TTL := 180 * time.Second
	sessionId := uuid.Must(uuid.NewV4()).String()

	sessionPayload := m.SessionPayload{
		UserId:    userId,
		CreatedAt: time.Now(),
	}

	marshaledPayload, err := json.Marshal(sessionPayload)

	if err != nil {
		return nil, err
	}

	if result := cfg.rClient.Set(ctx, sessionId, marshaledPayload, TTL); result.Err() != nil {
		return nil, result.Err()
	}

	return &m.Session{ID: sessionId, Payload: sessionPayload, TTL: TTL}, nil
}

func (cfg *SessionOperationsConfig) GetSession(ctx context.Context, sessionId string) (*m.Session, error) {
	var sessionPayload m.SessionPayload

	exist, err := cfg.rClient.Exists(ctx, sessionId).Result()

	if err != nil {
		return nil, err
	}

	if exist != 1 {
		return nil, nil
	}

	record := cfg.rClient.Get(ctx, sessionId)
	recordTTL := cfg.rClient.TTL(ctx, sessionId)

	if record.Err() != nil {
		return nil, record.Err()
	}

	recordInfo, _ := record.Result()
	TTLInfo, _ := recordTTL.Result()

	if err := json.Unmarshal([]byte(recordInfo), &sessionPayload); err != nil {
		return nil, err
	}

	return &m.Session{ID: sessionId, Payload: sessionPayload, TTL: TTLInfo}, nil
}

func (cfg *SessionOperationsConfig) DeleteSession(ctx context.Context, sessionId string) error {
	if err := cfg.rClient.Del(ctx, sessionId).Err(); err != nil {
		return err
	}

	return nil
}

func (cfg *SessionOperationsConfig) UpdateSession(ctx context.Context, sessionId string) (*m.Session, error) {
	session, err := cfg.GetSession(ctx, sessionId)

	if err != nil {
		return nil, err
	}

	if err := cfg.DeleteSession(ctx, sessionId); err != nil {
		return nil, err
	}

	return cfg.CreateSession(ctx, session.Payload.UserId)
}
