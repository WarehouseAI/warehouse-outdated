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
	Description       string `json:"description"`
	Name              string `json:"name"`
	AuthHeaderName    string `json:"auth_header_name"`
	AuthHeaderContent string `json:"auth_header_content"`
	Image             string `json:"image"`
}

type CreateWithKeyRequest struct {
	Description       string `json:"description"`
	Name              string `json:"name"`
	AuthHeaderName    string `json:"auth_header_name"`
	AuthHeaderContent string `json:"auth_header_content"`
	Image             string `json:"image"`
}

type CreateResponse struct {
	ID                string `json:"id"`
	AuthHeaderContent string `json:"auth_header_content"`
	AuthHeaderName    string `json:"auth_header_name"`
}

func CreateWithGeneratedKey(aiInfo *CreateWithoutKeyRequest, userId string, ai dataservice.AiInterface, logger *logrus.Logger) (*CreateResponse, *e.HttpErrorResponse) {
	key, err := generateToken(32)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, e.NewErrorResponse(e.HttpInternalError, err.Error())
	}

	apiKey := fmt.Sprintf("wh.%s", key)

	newAI := &m.AiProduct{
		Name:              aiInfo.Name,
		Description:       aiInfo.Description,
		Owner:             uuid.Must(uuid.FromString(userId)),
		AuthHeaderName:    aiInfo.AuthHeaderName,
		AuthHeaderContent: apiKey,
		BackgroundUrl:     aiInfo.Image,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if dbErr := ai.Create(newAI); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &CreateResponse{
		ID:                newAI.ID.String(),
		AuthHeaderContent: newAI.AuthHeaderContent,
		AuthHeaderName:    newAI.AuthHeaderName,
	}, nil
}

func CreateWithOwnKey(aiInfo *CreateWithKeyRequest, userId string, ai dataservice.AiInterface, logger *logrus.Logger) (*CreateResponse, *e.HttpErrorResponse) {
	newAI := &m.AiProduct{
		Name:              aiInfo.Name,
		Description:       aiInfo.Description,
		Owner:             uuid.Must(uuid.FromString(userId)),
		AuthHeaderName:    aiInfo.AuthHeaderName,
		BackgroundUrl:     aiInfo.Image,
		AuthHeaderContent: aiInfo.AuthHeaderContent,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if dbErr := ai.Create(newAI); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, e.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &CreateResponse{
		ID:                newAI.ID.String(),
		AuthHeaderContent: newAI.AuthHeaderContent,
		AuthHeaderName:    newAI.AuthHeaderName,
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
