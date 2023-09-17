package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	dbm "warehouse/src/internal/db/models"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"
	u "warehouse/src/internal/utils"
	m "warehouse/src/services/ai/pkg/models"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type AIService interface {
	Create(context.Context, *m.CreateAIRequest, dbm.User) (*m.CreateAIResponse, error)
}

type AIServiceConfig struct {
	operations dbo.AIDatabaseOperations
	logger     *logrus.Logger
}

func NewAIService(operations dbo.AIDatabaseOperations, logger *logrus.Logger) AIService {
	return &AIServiceConfig{
		operations: operations,
		logger:     logger,
	}
}

func (cfg *AIServiceConfig) Create(ctx context.Context, aiInfo *m.CreateAIRequest, user dbm.User) (*m.CreateAIResponse, error) {
	apiKeyPayload, err := u.GenerateRandomString(32)

	if err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, dto.InternalError
	}

	apiKey := fmt.Sprintf("wh.%s", apiKeyPayload)
	hash := sha256.Sum256([]byte(apiKey))

	newAi := &dbm.AI{ID: uuid.Must(uuid.NewV4()), Name: aiInfo.Name, Owner: user.ID, AuthScheme: aiInfo.AuthScheme, ApiKey: string(hash[:]), CreatedAt: time.Now(), UpdateAt: time.Now()}

	if err := cfg.operations.Add(newAi); err != nil {
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, dto.InternalError
	}

	return &m.CreateAIResponse{Name: aiInfo.Name, ApiKey: apiKey, AuthScheme: aiInfo.AuthScheme}, nil
}
