package model

import (
	"warehouse/gen"
)

func UserIdFromProto(m *gen.CreateUserResponse) *UserIdResponse {
	return &UserIdResponse{
		ID: m.Id,
	}
}
