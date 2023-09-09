package model

import (
	"time"
	"warehouse/gen"
	d "warehouse/src/services/user/internal/datastore"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func UserToProto(m d.User) *gen.User {
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

func UserPayloadToEntity(m *gen.CreateUserRequest) *d.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(m.Password), 12)

	return &d.User{
		ID:        uuid.Must(uuid.NewV4()),
		Username:  m.Username,
		Password:  string(hash),
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}
