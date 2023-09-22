package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"
	u "warehouse/src/internal/utils"
	m "warehouse/src/services/ai/pkg/models"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AIService interface {
	Create(context.Context, *m.CreateAIRequest, *dbm.User) (*m.CreateAIResponse, error)
	AddCommand(context.Context, *m.AddCommandRequest) error
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
		ID:          uuid.Must(uuid.NewV4()),
		Name:        commandInfo.Name,
		AI:          commandInfo.AiID,
		RequestType: commandInfo.RequestType,
		Payload:     commandInfo.Payload,
		PayloadType: commandInfo.PayloadType,
		URL:         commandInfo.URL,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}

	if err := commandOperations.Add(newCommand); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Add new command to AI")
		return err
	}

	return nil
}
