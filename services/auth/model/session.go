package model

import "time"

type (
	Session struct {
		ID      string         `json:"id"`
		Payload SessionPayload `json:"payload"`
		TTL     time.Duration  `json:"ttl"`
	}

	SessionPayload struct {
		UserId    string    `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	}
)
