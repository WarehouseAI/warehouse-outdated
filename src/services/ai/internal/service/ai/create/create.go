package create

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

func CreateWithGeneratedKey(aiInfo *RequestWithoutKey, user *pg.User, aiCreator AICreator, logger *logrus.Logger, ctx context.Context) (*Response, *httputils.ErrorResponse) {
	key, err := u.GenerateKey(32)
	hasher := md5.New()

	if err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Create new AI")
		return nil, httputils.NewErrorResponse(httputils.InternalError, err.Error())
	}

	apiKey := fmt.Sprintf("wh.%s", key)
	hasher.Write([]byte(key))

	newAI := &pg.AI{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       aiInfo.Name,
		Owner:      user.ID,
		AuthScheme: aiInfo.AuthScheme,
		ApiKey:     hex.EncodeToString(hasher.Sum(nil)),
		CreatedAt:  time.Now(),
		UpdateAt:   time.Now(),
	}

	if dbErr := aiCreator.Add(newAI); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &Response{Name: aiInfo.Name, ApiKey: apiKey, AuthScheme: aiInfo.AuthScheme}, nil
}

func CreateWithOwnKey(aiInfo *RequestWithKey, user *pg.User, aiCreator AICreator, logger *logrus.Logger, ctx context.Context) (*Response, *httputils.ErrorResponse) {
	newAI := &pg.AI{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       aiInfo.Name,
		Owner:      user.ID,
		AuthScheme: aiInfo.AuthScheme,
		ApiKey:     aiInfo.AuthKey,
		CreatedAt:  time.Now(),
		UpdateAt:   time.Now(),
	}

	if dbErr := aiCreator.Add(newAI); dbErr != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Create new AI")
		return nil, httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
	}

	return &Response{Name: aiInfo.Name, ApiKey: aiInfo.AuthKey, AuthScheme: aiInfo.AuthScheme}, nil
}
