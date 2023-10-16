package create

import (
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type Request struct {
	Name        string                 `json:"name"`
	AiID        uuid.UUID              `json:"ai_id"`
	Payload     map[string]interface{} `json:"payload"`
	PayloadType pg.PayloadType         `json:"payload_type"`
	InputType   pg.IOType              `json:"input_type"`
	OutputType  pg.IOType              `json:"output_type"`
	RequestType pg.RequestScheme       `json:"request_type"`
	URL         string                 `json:"url"`
}

type CommandProvider interface {
	GetOneBy(key string, value interface{}) (*pg.Command, error)
	Add(item *pg.Command) error
}

func CreateCommand(commandCreds *Request, commandProvider CommandProvider, logger *logrus.Logger) error {
	// TODO: Перенести Get проверку в Add
	existCommand, err := commandProvider.GetOneBy("name", commandCreds.Name)

	if existCommand != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add command")
		return dto.ExistError
	}

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add command")
		return dto.InternalError
	}

	newCommand := &pg.Command{
		ID:            uuid.Must(uuid.NewV4()),
		Name:          commandCreds.Name,
		AI:            commandCreds.AiID,
		RequestScheme: commandCreds.RequestType,
		InputType:     commandCreds.InputType,
		OutputType:    commandCreds.OutputType,
		Payload:       commandCreds.Payload,
		PayloadType:   commandCreds.PayloadType,
		URL:           commandCreds.URL,
		CreatedAt:     time.Now(),
		UpdateAt:      time.Now(),
	}

	if err := commandProvider.Add(newCommand); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add new command to AI")
		return err
	}

	return nil
}
