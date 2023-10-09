package mapper

import (
	"encoding/json"
	"time"
	"warehouse/gen"
	dbm "warehouse/src/internal/db/models"

	"github.com/gofrs/uuid"
)

func UserPayloadToEntity(m *gen.CreateUserMsg) *dbm.User {
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

func UserToProto(m *dbm.User) *gen.User {
	var ownedAi []*gen.AI

	for _, s := range m.OwnedAi {
		ownedAi = append(ownedAi, AiToProto(s))
	}

	return &gen.User{
		Id:        m.ID.String(),
		Username:  m.Username,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		OwnedAi:   ownedAi,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: m.CreatedAt.String(),
		UpdatedAt: m.UpdateAt.String(),
	}
}

func ProtoToUser(m *gen.User) *dbm.User {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)
	var ownedAi []dbm.AI

	for _, s := range m.OwnedAi {
		ownedAi = append(ownedAi, ProtoToAi(s))
	}

	return &dbm.User{
		ID:        uuid.FromStringOrNil(m.Id),
		Username:  m.Username,
		Password:  m.Password,
		Picture:   m.Picture,
		Email:     m.Email,
		ViaGoogle: m.ViaGoogle,
		CreatedAt: createdAt,
		UpdateAt:  updatedAt,
		OwnedAi:   ownedAi,
	}
}

func AiToProto(m dbm.AI) *gen.AI {
	var commands []*gen.Command

	for _, s := range m.Commands {
		commands = append(commands, CommandToProto(s))
	}

	return &gen.AI{
		Id:         m.ID.String(),
		Owner:      m.Owner.String(),
		Name:       m.Name,
		ApiKey:     m.ApiKey,
		AuthScheme: string(m.AuthScheme),
		CreatedAt:  m.CreatedAt.String(),
		UpdatedAt:  m.UpdateAt.String(),
		Commands:   commands,
	}
}

func ProtoToAi(m *gen.AI) dbm.AI {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)
	var commands []dbm.Command

	for _, s := range m.Commands {
		commands = append(commands, ProtoToCommand(s))
	}

	return dbm.AI{
		ID:         uuid.FromStringOrNil(m.Id),
		Owner:      uuid.FromStringOrNil(m.Owner),
		Name:       m.Name,
		Commands:   commands,
		ApiKey:     m.ApiKey,
		AuthScheme: dbm.AuthScheme(m.AuthScheme),
		CreatedAt:  createdAt,
		UpdateAt:   updatedAt,
	}
}

func CommandToProto(m dbm.Command) *gen.Command {
	jsonObject, _ := m.Payload.MarshalJSON()

	return &gen.Command{
		Id:            m.ID.String(),
		Ai:            m.AI.String(),
		Name:          m.Name,
		Payload:       string(jsonObject),
		PayloadType:   string(m.PayloadType),
		RequestScheme: string(m.RequestScheme),
		InputType:     string(m.InputType),
		OutputType:    string(m.OutputType),
		Url:           m.URL,
		CreatedAt:     m.CreatedAt.String(),
		UpdatedAt:     m.UpdateAt.String(),
	}
}

func ProtoToCommand(m *gen.Command) dbm.Command {
	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, m.UpdatedAt)
	var jsonObject map[string]interface{}

	json.Unmarshal([]byte(m.Payload), &jsonObject)

	return dbm.Command{
		ID:            uuid.FromStringOrNil(m.Id),
		AI:            uuid.FromStringOrNil(m.Ai),
		Name:          m.Name,
		Payload:       jsonObject,
		PayloadType:   dbm.PayloadType(m.PayloadType),
		RequestScheme: dbm.RequestScheme(m.RequestScheme),
		InputType:     dbm.IOType(m.InputType),
		OutputType:    dbm.IOType(m.OutputType),
		URL:           m.Url,
		CreatedAt:     createdAt,
		UpdateAt:      updatedAt,
	}
}
