package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"
	u "warehouse/src/internal/utils"
	"warehouse/src/internal/utils/httputils"
	m "warehouse/src/services/ai/pkg/models"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AIService interface {
	Create(ctx context.Context, aiInfo *m.CreateAIRequest, user *dbm.User) (*m.CreateAIResponse, error)
	Get(ctx context.Context, aiID uuid.UUID) (*dbm.AI, error)
	AddCommand(ctx context.Context, commandInfo *m.AddCommandRequest) error
	GetCommand(ctx context.Context, aiId string, commandName string) (*dbm.Command, error)
	ExecuteJSONCommand(ctx context.Context, jsonData map[string]interface{}, command *dbm.Command) (*bytes.Buffer, error)
	ExecuteFormDataCommand(ctx context.Context, formData *multipart.Form, command *dbm.Command) (*bytes.Buffer, error)
}

type AIServiceConfig struct {
	database *gorm.DB
	logger   *logrus.Logger
}

func NewAIService(database *gorm.DB, logger *logrus.Logger) AIService {
	return &AIServiceConfig{
		database: database,
		logger:   logger,
	}
}

func (cfg *AIServiceConfig) Create(ctx context.Context, aiInfo *m.CreateAIRequest, user *dbm.User) (*m.CreateAIResponse, error) {
	aiOperations := dbo.NewAIOperations[dbm.AI](cfg.database)
	apiKeyPayload, err := u.GenerateRandomString(32)
	hasher := md5.New()

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, dto.InternalError
	}

	apiKey := fmt.Sprintf("wh.%s", apiKeyPayload)
	hasher.Write([]byte(apiKey))

	newAI := dbm.AI{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       aiInfo.Name,
		Owner:      user.ID,
		AuthScheme: aiInfo.AuthScheme,
		ApiKey:     hex.EncodeToString(hasher.Sum(nil)),
		CreatedAt:  time.Now(),
		UpdateAt:   time.Now(),
	}

	if err := aiOperations.Add(newAI); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, dto.InternalError
	}

	return &m.CreateAIResponse{Name: aiInfo.Name, ApiKey: apiKey, AuthScheme: aiInfo.AuthScheme}, nil
}

func (cfg *AIServiceConfig) Get(ctx context.Context, aiID uuid.UUID) (*dbm.AI, error) {
	aiOperations := dbo.NewAIOperations[dbm.AI](cfg.database)

	existAI, err := aiOperations.GetOneBy("id", aiID)

	if existAI == nil {
		return nil, dto.NotFoundError
	}

	if err != nil {
		return nil, dto.InternalError
	}

	return existAI, nil
}

func (cfg *AIServiceConfig) AddCommand(ctx context.Context, commandInfo *m.AddCommandRequest) error {
	commandOperations := dbo.NewAIOperations[dbm.Command](cfg.database)

	existCommand, err := commandOperations.GetOneBy("name", commandInfo.Name)

	if existCommand != nil {
		return dto.ExistError
	}

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add command")
		return dto.InternalError
	}

	newCommand := dbm.Command{
		ID:            uuid.Must(uuid.NewV4()),
		Name:          commandInfo.Name,
		AI:            commandInfo.AiID,
		RequestScheme: commandInfo.RequestType,
		InputType:     commandInfo.InputType,
		OutputType:    commandInfo.OutputType,
		Payload:       commandInfo.Payload,
		PayloadType:   commandInfo.PayloadType,
		URL:           commandInfo.URL,
		CreatedAt:     time.Now(),
		UpdateAt:      time.Now(),
	}

	if err := commandOperations.Add(newCommand); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add new command to AI")
		return err
	}

	return nil
}

func (cfg *AIServiceConfig) GetCommand(ctx context.Context, aiID string, commandName string) (*dbm.Command, error) {
	commandOperations := dbo.NewAIOperations[dbm.Command](cfg.database)

	existCommand, err := commandOperations.GetOneBy("name", commandName)

	if existCommand == nil {
		return nil, nil
	}

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add command")
		return nil, dto.InternalError
	}

	return existCommand, nil
}

func (cfg *AIServiceConfig) ExecuteFormDataCommand(ctx context.Context, formData *multipart.Form, command *dbm.Command) (*bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range command.Payload {
		if value == "file" {
			fileHeader := formData.File[key][0]
			file, err := fileHeader.Open()

			if err != nil {
				cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
				return nil, dto.InternalError
			}

			field, _ := writer.CreateFormFile(key, fileHeader.Filename)
			io.Copy(field, file)

			defer file.Close()
		} else {
			writer.WriteField(key, formData.Value[key][0])
		}
	}

	writer.Close()

	headers := map[string]string{
		"Content-Type": "multipart/form-data",
	}

	response, err := httputils.MakeHTTPRequest(command.URL, string(command.RequestScheme), headers, url.Values{}, body)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
		return nil, err
	}

	buffer, err := httputils.DecodeHTTPResponse(response, command.OutputType)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
		return nil, dto.InternalError
	}

	return buffer, nil
}

func (cfg *AIServiceConfig) ExecuteJSONCommand(ctx context.Context, jsonData map[string]interface{}, command *dbm.Command) (*bytes.Buffer, error) {
	json, err := json.Marshal(jsonData)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, dto.InternalError
	}

	ai, err := cfg.Get(ctx, command.AI)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("%s %s", string(ai.AuthScheme), ai.ApiKey),
	}

	body := bytes.NewBuffer(json)
	response, err := httputils.MakeHTTPRequest(command.URL, string(command.RequestScheme), headers, url.Values{}, body)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute JSON-Command")
		return nil, err
	}

	buffer, err := httputils.DecodeHTTPResponse(response, command.OutputType)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Execute FormData-Command")
		return nil, dto.InternalError
	}

	return buffer, nil
}
