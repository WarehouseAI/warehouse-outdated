package logout

import (
	"context"
	"time"
	db "warehouse/src/internal/database"
	"warehouse/src/internal/utils/httputils"

	"github.com/sirupsen/logrus"
)

type SessionFlusher interface {
	Delete(context.Context, string) *db.DBError
}

func Logout(sessionId string, sessionFlusher SessionFlusher, logger *logrus.Logger, ctx context.Context) *httputils.ErrorResponse {
	if err := sessionFlusher.Delete(ctx, sessionId); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Payload}).Info("Logout user")
		return httputils.NewErrorResponseFromDBError(err.ErrorType, err.Message)
	}

	return nil
}
