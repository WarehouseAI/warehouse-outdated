package logout

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type SessionFlusher interface {
	Delete(context.Context, string) error
}

func Logout(sessionId string, sessionFlusher SessionFlusher, logger *logrus.Logger, ctx context.Context) error {
	if err := sessionFlusher.Delete(ctx, sessionId); err != nil {
		logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Logout user")
		return err
	}

	return nil
}
