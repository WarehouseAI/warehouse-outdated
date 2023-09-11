package model

import (
	"context"
	"time"
	"warehouse/gen"
)

type (
	AuthService interface {
		Login(context.Context, *LoginRequest) (*Session, error)
		Register(context.Context, *gen.CreateUserRequest) (*UserIdResponse, error)
		// Refresh(context.Context) (m.TokenPair, error)
	}

	UserIdResponse struct {
		ID string `json:"id"`
	}

	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Session struct {
		ID      string        `json:"id"`
		Payload string        `json:"payload"`
		TTL     time.Duration `json:"ttl"`
	}

	SessionPayload struct {
		UserId        string `json:"userId"`
		Authenticated bool   `json:"authenticated"`
	}
)
