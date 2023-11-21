package command

import (
	"time"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type GetCommandRequest struct {
	AiID uuid.UUID `json:"ai_id"`
	Name string    `json:"name"`
}

type GetCommandResponse struct {
	AI         *m.AI
	Payload    m.Command
	ApiKey     string
	AuthScheme m.AuthScheme
}

func GetCommand(getRequest GetCommandRequest, aiProvider dataservice.AiInterface, logger *logrus.Logger) (*GetCommandResponse, *e.ErrorResponse) {
	existAI, dbErr := aiProvider.GetWithPreload(map[string]interface{}{"id": getRequest.AiID}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get command")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	for i := 0; i <= len(existAI.Commands); i++ {
		if existAI.Commands[i].Name == getRequest.Name {
			return &GetCommandResponse{existAI, existAI.Commands[i], existAI.ApiKey, existAI.AuthScheme}, nil
		}
	}

	return nil, e.NewErrorResponse(e.HttpNotFound, "Command not found.")
}
