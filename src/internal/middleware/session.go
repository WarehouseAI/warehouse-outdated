package middleware

import (
	"context"
	"time"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type MiddlewareService interface {
	Session(c *fiber.Ctx) error
}

type MiddlewareConfig struct {
	operations dbo.SessionDatabaseOperations
	logger     *logrus.Logger
}

func NewMiddlewareService(operations dbo.SessionDatabaseOperations, logger *logrus.Logger) MiddlewareService {
	return &MiddlewareConfig{
		operations: operations,
		logger:     logger,
	}
}

func (cfg *MiddlewareConfig) Session(c *fiber.Ctx) error {
	sessionId := c.Cookies("id")

	if sessionId == "" {
		statusCode := fiber.StatusUnauthorized
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session has expired"})
	}

	session, err := cfg.operations.GetSession(context.Background(), sessionId)

	if session == nil {
		c.Cookie(&fiber.Cookie{
			Name:    "id",
			Value:   "",
			Expires: time.Now().Add(-time.Hour),
		})

		statusCode := fiber.StatusUnauthorized
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session has expired"})
	}

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	newSession, err := cfg.operations.UpdateSession(context.Background(), sessionId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "id",
		Value: newSession.ID,
	})

	return c.Next()

}
