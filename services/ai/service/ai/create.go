package ai

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	m "warehouseai/ai/model"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type CreateWithoutKeyRequest struct {
	Description string       `json:"description"`
	Name        string       `json:"name"`
	AuthScheme  m.AuthScheme `json:"auth_scheme"`
	Image       string       `json:"image"`
}

type CreateWithKeyRequest struct {
	Description string       `json:"description"`
	Name        string       `json:"name"`
	AuthScheme  m.AuthScheme `json:"auth_scheme"`
	ApiKey      string       `json:"api_key"`
	Image       string       `json:"image"`
}

type CreateResponse struct {
	ID         string `json:"id"`
	ApiKey     string `json:"api_key"`
	AuthScheme string `json:"auth_scheme"`
}

func CreateWithGeneratedKey(aiInfo *CreateWithoutKeyRequest, userId string, ai dataservice.AiInterface, logger *logrus.Logger) (*CreateResponse, *e.ErrorResponse) {
	key, err := generateToken(32)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	apiKey := fmt.Sprintf("wh.%s", key)

	newAI := &m.AI{
		ID:            uuid.Must(uuid.NewV4()),
		Name:          aiInfo.Name,
		Description:   aiInfo.Description,
		Owner:         uuid.Must(uuid.FromString(userId)),
		AuthScheme:    aiInfo.AuthScheme,
		ApiKey:        apiKey,
		BackgroundUrl: aiInfo.Image,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if dbErr := ai.Create(newAI); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &CreateResponse{
		ID:         newAI.ID.String(),
		ApiKey:     apiKey,
		AuthScheme: string(newAI.AuthScheme),
	}, nil
}

func CreateWithOwnKey(aiInfo *CreateWithKeyRequest, userId string, ai dataservice.AiInterface, logger *logrus.Logger) (*CreateResponse, *e.ErrorResponse) {
	newAI := &m.AI{
		ID:            uuid.Must(uuid.NewV4()),
		Name:          aiInfo.Name,
		Description:   aiInfo.Description,
		Owner:         uuid.Must(uuid.FromString(userId)),
		AuthScheme:    aiInfo.AuthScheme,
		BackgroundUrl: aiInfo.Image,
		ApiKey:        aiInfo.ApiKey,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if dbErr := ai.Create(newAI); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &CreateResponse{
		ID:         newAI.ID.String(),
		ApiKey:     aiInfo.ApiKey,
		AuthScheme: string(newAI.AuthScheme),
	}, nil
}

func generateToken(length int) (string, error) {
	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(randomBytes)
	key = key[:length]

	return key, nil
}
