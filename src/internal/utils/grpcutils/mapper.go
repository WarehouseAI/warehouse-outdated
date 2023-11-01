package grpcutils

import (
	"encoding/json"
	"time"
	"warehouse/gen"
	pg "warehouse/src/internal/database/postgresdb"

	"github.com/gofrs/uuid"
)

func UserPayloadToEntity(m *gen.CreateUserMsg) *pg.User {
	return &pg.User{
		ID:        uuid.Must(uuid.NewV4()),
		Username:  m.Username,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func UserToProto(m *pg.User) *gen.User {
	var ownedAi []*gen.AI

	for _, s := range m.OwnedAi {
		ownedAi = append(ownedAi, AiToProto(&s))
	}

	return &gen.User{
		Id:        m.ID.String(),
		Username:  m.Username,
		Password:  m.Password,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Picture:   m.Picture,
		Email:     m.Email,
		OwnedAi:   ownedAi,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdatedAt.String(),
	}
}

func ProtoToUser(m *gen.User) *pg.User {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)
	var ownedAi []pg.AI

	for _, s := range m.OwnedAi {
		ownedAi = append(ownedAi, ProtoToAi(s))
	}

	return &pg.User{
		ID:        uuid.FromStringOrNil(m.Id),
		Username:  m.Username,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		OwnedAi:   ownedAi,
	}
}

func AiToProto(m *pg.AI) *gen.AI {
	var commands []*gen.Command

	for _, s := range m.Commands {
		commands = append(commands, CommandToProto(&s))
	}

	return &gen.AI{
		Id:         m.ID.String(),
		Owner:      m.Owner.String(),
		Name:       m.Name,
		ApiKey:     m.ApiKey,
		AuthScheme: string(m.AuthScheme),
		CreatedAt:  m.CreatedAt.String(),
		UpdatedAt:  m.UpdatedAt.String(),
		Commands:   commands,
	}
}

func ProtoToAi(m *gen.AI) pg.AI {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)
	var commands []pg.Command

	for _, s := range m.Commands {
		commands = append(commands, ProtoToCommand(s))
	}

	return pg.AI{
		ID:         uuid.FromStringOrNil(m.Id),
		Owner:      uuid.FromStringOrNil(m.Owner),
		Name:       m.Name,
		Commands:   commands,
		ApiKey:     m.ApiKey,
		AuthScheme: pg.AuthScheme(m.AuthScheme),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func CommandToProto(m *pg.Command) *gen.Command {
	jsonObject, _ := m.Payload.MarshalJSON()

	return &gen.Command{
		Id:            m.ID.String(),
		Ai:            m.AIID.String(),
		Name:          m.Name,
		Payload:       string(jsonObject),
		PayloadType:   string(m.PayloadType),
		RequestScheme: string(m.RequestScheme),
		InputType:     string(m.InputType),
		OutputType:    string(m.OutputType),
		Url:           m.URL,
		CreatedAt:     m.CreatedAt.String(),
		UpdatedAt:     m.UpdatedAt.String(),
	}
}

func ProtoToCommand(m *gen.Command) pg.Command {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)
	var jsonObject map[string]interface{}

	json.Unmarshal([]byte(m.Payload), &jsonObject)

	return pg.Command{
		ID:            uuid.FromStringOrNil(m.Id),
		AIID:          uuid.FromStringOrNil(m.Ai),
		Name:          m.Name,
		Payload:       jsonObject,
		PayloadType:   pg.PayloadType(m.PayloadType),
		RequestScheme: pg.RequestScheme(m.RequestScheme),
		InputType:     pg.IOType(m.InputType),
		OutputType:    pg.IOType(m.OutputType),
		URL:           m.Url,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
