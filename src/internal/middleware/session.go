package middleware

import (
	"context"
	"time"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SessionMiddleware struct {
	sessionOperations dbo.SessionDatabaseOperations
	logger            *logrus.Logger
}

func NewSessionMiddleware(sessionOperations dbo.SessionDatabaseOperations, logger *logrus.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		sessionOperations: sessionOperations,
		logger:            logger,
	}
}

func (cfg *SessionMiddleware) Session(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if sessionId == "" {
		statusCode := fiber.StatusUnauthorized
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session is invalid"})
	}

	session, err := cfg.sessionOperations.GetSession(context.Background(), sessionId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	if session == nil {
		c.Cookie(&fiber.Cookie{
			Name:    "sessionId",
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
		})

		statusCode := fiber.StatusForbidden
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session has expired"})
	}

	newSession, err := cfg.sessionOperations.UpdateSession(context.Background(), sessionId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	c.Locals("userId", newSession.Payload.UserId)
	c.Cookie(&fiber.Cookie{
		Name:  "sessionId",
		Value: newSession.ID,
	})

	return c.Next()
}
