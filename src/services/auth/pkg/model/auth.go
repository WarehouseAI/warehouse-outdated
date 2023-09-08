package model

import (
	"context"
	"warehouse/gen"
)

type (
	AuthService interface {
		// Login(context.Context) (m.TokenPair, error)
		Register(context.Context, *gen.CreateUserRequest) (*string, error)
		// Refresh(context.Context) (m.TokenPair, error)
	}

	TokenPair struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)
