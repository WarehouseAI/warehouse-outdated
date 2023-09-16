package models

import (
	"time"
	"warehouse/gen"
	dbm "warehouse/src/internal/db/models"

	"github.com/gofrs/uuid"
)

func UserToProto(m *dbm.User) *gen.User {
	return &gen.User{
		Id:        m.ID.String(),
		Username:  m.Username,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdateAt.String(),
	}
}

func UserPayloadToEntity(m *gen.CreateUserRequest) *dbm.User {
	return &dbm.User{
		ID:        uuid.Must(uuid.NewV4()),
		Username:  m.Username,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}
