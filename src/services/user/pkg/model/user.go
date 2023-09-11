package model

import (
	"context"
	"warehouse/gen"
	d "warehouse/src/services/user/internal/datastore"
)

type (
	UserService interface {
		Create(context.Context, *gen.CreateUserRequest) (*d.User, error)
		Get(context.Context, *gen.GetUserRequest) (*d.User, error)
	}

	CreateUserDTO struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Picture   string `json:"picture"`
		Email     string `json:"email"`
		ViaGoogle bool   `json:"viaGoogle"`
	}
)
