package models

import "time"

type (
	Session struct {
		ID      string         `json:"id"`
		Payload SessionPayload `json:"payload"`
		TTL     time.Duration  `json:"ttl"`
	}

	SessionPayload struct {
		UserId    string    `json:"userId"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
