package datastore

import (
	"context"
	"encoding/json"
	"time"
	m "warehouse/src/services/auth/pkg/model"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionDatabaseOperations interface {
	CreateSession(context.Context, string) (*m.Session, error)
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
	TTL := 180 * time.Second
	sessionId := uuid.Must(uuid.NewV4()).String()

	authenticatedUser := m.SessionPayload{
		UserId:        userId,
		Authenticated: true,
	}

	marshaledPayload, err := json.Marshal(authenticatedUser)

	if err != nil {
		return nil, err
	}

	// Поменять потом на 3 дня
	if result := cfg.rClient.Set(ctx, sessionId, marshaledPayload, TTL); result.Err() != nil {
		return nil, result.Err()
	}

	return &m.Session{ID: sessionId, Payload: string(marshaledPayload), TTL: TTL}, nil
}
