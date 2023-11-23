package mapper

import (
	"time"
	"warehouseai/user/adapter/grpc/gen"
	"warehouseai/user/model"

	"github.com/gofrs/uuid"
)

func UserToProto(m *model.User) *gen.User {

	return &gen.User{
		Id:        m.ID.String(),
		Username:  m.Username,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Picture:   m.Picture,
		Email:     m.Email,
		Password:  m.Password,
		Verified:  m.Verified,
		Role:      string(m.Role),
		ViaGoogle: m.ViaGoogle,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
	}
}

func ProtoToUser(m *gen.User) *model.User {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)

	return &model.User{
		ID:        uuid.FromStringOrNil(m.Id),
		Username:  m.Username,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Picture:   m.Picture,
		Email:     m.Email,
		Password:  m.Password,
		Verified:  m.Verified,
		Role:      model.UserRole(m.Role),
		ViaGoogle: m.ViaGoogle,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
