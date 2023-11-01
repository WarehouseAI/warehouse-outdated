package create

import (
	"context"
	"fmt"
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	u "warehouse/src/internal/utils"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type RequestWithoutKey struct {
	Name       string        `json:"name"`
	AuthScheme pg.AuthScheme `json:"auth_scheme"`
}

type RequestWithKey struct {
	Name       string        `json:"name"`
	AuthScheme pg.AuthScheme `json:"auth_scheme"`
	AuthKey    string        `json:"auth_key"`
}

type Response struct {
	Name       string        `json:"name"`
	AuthScheme pg.AuthScheme `json:"auth_scheme"`
	ApiKey     string        `json:"api_key"`
}

type AICreator interface {
	Add(item *pg.AI) *db.DBError
}

func CreateWithGeneratedKey(aiInfo *RequestWithoutKey, userId string, aiCreator AICreator, logger *logrus.Logger, ctx context.Context) (*pg.AI, *httputils.ErrorResponse) {
	key, err := u.GenerateKey(32)

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	apiKey := fmt.Sprintf("wh.%s", key)

	newAI := &pg.AI{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       aiInfo.Name,
		Owner:      uuid.Must(uuid.FromString(userId)),
		AuthScheme: aiInfo.AuthScheme,
		ApiKey:     apiKey,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if dbErr := aiCreator.Add(newAI); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return newAI, nil
}

func CreateWithOwnKey(aiInfo *RequestWithKey, userId string, aiCreator AICreator, logger *logrus.Logger, ctx context.Context) (*pg.AI, *httputils.ErrorResponse) {
	newAI := &pg.AI{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       aiInfo.Name,
		Owner:      uuid.Must(uuid.FromString(userId)),
		AuthScheme: aiInfo.AuthScheme,
		ApiKey:     aiInfo.AuthKey,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if dbErr := aiCreator.Add(newAI); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return newAI, nil
}
