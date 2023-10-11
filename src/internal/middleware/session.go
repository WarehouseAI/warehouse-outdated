package middleware

import (
	"context"
	"time"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SessionProvider interface {
	Get(ctx context.Context, sessionId string) (*r.Session, error)
	Update(ctx context.Context, sessionId string) (*r.Session, error)
}

type SessionMiddlewareProvider struct {
	sessionProvider SessionProvider
	logger          *logrus.Logger
}

func NewSessionMiddleware(sessionProvider SessionProvider, logger *logrus.Logger) *SessionMiddlewareProvider {
	return &SessionMiddlewareProvider{
		sessionProvider: sessionProvider,
		logger:          logger,
	}
}

func (pvd *SessionMiddlewareProvider) Session(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if sessionId == "" {
		statusCode := fiber.StatusUnauthorized
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session is invalid"})
	}

	session, err := pvd.sessionProvider.Get(context.Background(), sessionId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		pvd.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	if session == nil {
		c.ClearCookie("sessionId")

		statusCode := fiber.StatusForbidden
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session has expired"})
	}

	newSession, err := pvd.sessionProvider.Update(context.Background(), sessionId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		pvd.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	c.Locals("userId", newSession.Payload.UserId)
	c.Cookie(&fiber.Cookie{
		Name:  "sessionId",
		Value: newSession.ID,
	})

	return c.Next()
}
