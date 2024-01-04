package mapper

import (
	"encoding/json"
	"time"
	"warehouseai/ai/adapter/grpc/gen"
	m "warehouseai/ai/model"

	"github.com/gofrs/uuid"
)

func AiToProto(ai *m.AiProduct) *gen.AI {
	var commands []*gen.Command

	for _, s := range ai.Commands {
		commands = append(commands, CommandToProto(&s))
	}

	return &gen.AI{
		Id:                ai.ID.String(),
		Owner:             ai.Owner.String(),
		Name:              ai.Name,
		AuthHeaderContent: ai.AuthHeaderContent,
		AuthHeaderName:    ai.AuthHeaderName,
		CreatedAt:         ai.CreatedAt.String(),
		UpdatedAt:         ai.UpdatedAt.String(),
		Commands:          commands,
	}
}

func ProtoToAi(ai *gen.AI) m.AiProduct {
	createdAt, _ := time.Parse(time.RFC3339, ai.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, ai.UpdatedAt)
	var commands []m.AiCommand

	for _, s := range ai.Commands {
		commands = append(commands, ProtoToCommand(s))
	}

	return m.AiProduct{
		ID:                uuid.FromStringOrNil(ai.Id),
		Owner:             uuid.FromStringOrNil(ai.Owner),
		Name:              ai.Name,
		Commands:          commands,
		AuthHeaderContent: ai.AuthHeaderContent,
		AuthHeaderName:    ai.AuthHeaderName,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func CommandToProto(cmd *m.AiCommand) *gen.Command {
	jsonObject, _ := cmd.Payload.MarshalJSON()

	return &gen.Command{
		Id:          cmd.ID.String(),
		Ai:          cmd.AIID.String(),
		Name:        cmd.Name,
		Payload:     string(jsonObject),
		PayloadType: cmd.PayloadType,
		RequestType: cmd.RequestType,
		InputType:   cmd.InputType,
		OutputType:  cmd.OutputType,
		Url:         cmd.URL,
		CreatedAt:   cmd.CreatedAt.String(),
		UpdatedAt:   cmd.UpdatedAt.String(),
	}
}

func ProtoToCommand(cmd *gen.Command) m.AiCommand {
	createdAt, _ := time.Parse(time.RFC3339, cmd.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, cmd.UpdatedAt)
	var jsonObject map[string]interface{}

	json.Unmarshal([]byte(cmd.Payload), &jsonObject)

	return m.AiCommand{
		ID:          uuid.FromStringOrNil(cmd.Id),
		AIID:        uuid.FromStringOrNil(cmd.Ai),
		Name:        cmd.Name,
		Payload:     jsonObject,
		PayloadType: cmd.PayloadType,
		RequestType: cmd.RequestType,
		InputType:   cmd.InputType,
		OutputType:  cmd.OutputType,
		URL:         cmd.Url,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
