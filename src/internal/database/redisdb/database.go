package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionDatabase struct {
	rClient *redis.Client
}

func NewRedisDatabase(host string, port string, password string) *SessionDatabase {
	DSN := fmt.Sprintf("%s:%s", host, port)

	rClient := redis.NewClient(&redis.Options{
		Addr:     DSN,
		Password: password,
		DB:       0,
	})
	defer rClient.Close()

	return &SessionDatabase{
		rClient: rClient,
	}
}

func (cfg *SessionDatabase) Create(ctx context.Context, userId string) (*Session, error) {
	//TODO: Поменять потом на 3 дня
	TTL := 180 * time.Second
	sessionId := uuid.Must(uuid.NewV4()).String()

	sessionPayload := SessionPayload{
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

	return &Session{ID: sessionId, Payload: sessionPayload, TTL: TTL}, nil
}

func (cfg *SessionDatabase) Get(ctx context.Context, sessionId string) (*Session, error) {
	var sessionPayload SessionPayload

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

	return &Session{ID: sessionId, Payload: sessionPayload, TTL: TTLInfo}, nil
}

func (cfg *SessionDatabase) Delete(ctx context.Context, sessionId string) error {
	if err := cfg.rClient.Del(ctx, sessionId).Err(); err != nil {
		return err
	}

	return nil
}

func (cfg *SessionDatabase) Update(ctx context.Context, sessionId string) (*Session, error) {
	session, err := cfg.Get(ctx, sessionId)

	if err != nil {
		return nil, err
	}

	if err := cfg.Delete(ctx, sessionId); err != nil {
		return nil, err
	}

	return cfg.Create(ctx, session.Payload.UserId)
}
