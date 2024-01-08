package get

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
	AI                *m.AiProduct
	Command           m.AiCommand
	AuthHeaderContent string
	AuthHeaderName    string
}

func GetCommand(getRequest GetCommandRequest, aiProvider dataservice.AiInterface, logger *logrus.Logger) (*GetCommandResponse, *e.HttpErrorResponse) {
	existAI, dbErr := aiProvider.GetWithPreload(map[string]interface{}{"id": getRequest.AiID}, "Commands")

	if dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Get command")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	// Bug: Крашится если нет такой команды
	for i := 0; i < len(existAI.Commands); i++ {
		if existAI.Commands[i].Name == getRequest.Name {
			return &GetCommandResponse{existAI, existAI.Commands[i], existAI.AuthHeaderContent, existAI.AuthHeaderName}, nil
		}
	}

	return nil, e.NewErrorResponse(e.HttpNotFound, "Command not found.")
}
