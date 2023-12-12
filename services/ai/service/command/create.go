package command

import (
	"time"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type CreateCommandRequest struct {
	Name        string                 `json:"name"`
	AiID        string                 `json:"ai_id"`
	Payload     map[string]interface{} `json:"payload"`
	PayloadType m.PayloadType          `json:"payload_type"`
	InputType   m.IOType               `json:"input_type"`
	OutputType  m.IOType               `json:"output_type"`
	RequestType m.RequestScheme        `json:"request_type"`
	URL         string                 `json:"url"`
}

func CreateCommand(request *CreateCommandRequest, command dataservice.CommandInterface, logger *logrus.Logger) *e.ErrorResponse {
	newCommand := &m.Command{
		Name:          request.Name,
		AIID:          uuid.FromStringOrNil(request.AiID),
		RequestScheme: request.RequestType,
		InputType:     request.InputType,
		OutputType:    request.OutputType,
		Payload:       request.Payload,
		PayloadType:   request.PayloadType,
		URL:           request.URL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if dbErr := command.Create(newCommand); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Add new command to AI")
		return e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return nil
}
