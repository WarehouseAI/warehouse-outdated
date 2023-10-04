package models

import (
	dbm "warehouse/src/internal/db/models"

	"github.com/gofrs/uuid"
)

type (
	CreateAIRequest struct {
		Name       string         `json:"name"`
		AuthScheme dbm.AuthScheme `json:"auth_scheme"`
		AuthKey    string         `json:"auth_key"`
	}

	AddCommandRequest struct {
		Name        string                 `json:"name"`
		AiID        uuid.UUID              `json:"ai_id"`
		Payload     map[string]interface{} `json:"payload"`
		PayloadType dbm.PayloadType        `json:"payload_type"`
		InputType   dbm.IOType             `json:"input_type"`
		OutputType  dbm.IOType             `json:"output_type"`
		RequestType dbm.RequestScheme      `json:"request_type"`
		URL         string                 `json:"url"`
	}
)
