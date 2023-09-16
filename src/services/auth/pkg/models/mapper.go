package models

import (
	"warehouse/gen"
)

func UserIdFromProto(m *gen.CreateUserResponse) *RegisterResponse {
	return &RegisterResponse{
		ID: m.Id,
	}
}
