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
		Username:  m.Username,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
	}
}

func UserPayloadToEntity(m *gen.CreateUserRequest) *d.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(m.User.Password), 12)

	return &d.User{
		ID:        uuid.Must(uuid.NewV4()),
		Username:  m.User.Username,
		Password:  string(hash),
		Picture:   m.User.Picture,
		Email:     m.User.Email,
		ViaGoogle: m.User.ViaGoogle,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}
}
