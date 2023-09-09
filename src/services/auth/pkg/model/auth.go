package model

import (
	"context"
	"warehouse/gen"
)

type (
	AuthService interface {
		// Login(context.Context) (m.TokenPair, error)
		Register(context.Context, *gen.CreateUserRequest) (*UserIdResponse, error)
		// Refresh(context.Context) (m.TokenPair, error)
	}

	UserIdResponse struct {
		ID string `json:"id"`
	}

	TokenPair struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)
