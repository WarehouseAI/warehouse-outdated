package middleware

import (
	"context"
	"time"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type SessionMiddleware struct {
	database *redis.Client
	logger   *logrus.Logger
}

func NewSessionMiddleware(database *redis.Client, logger *logrus.Logger) *SessionMiddleware {
	return &SessionMiddleware{
		database: database,
		logger:   logger,
	}
}

func (cfg *SessionMiddleware) Session(c *fiber.Ctx) error {
	sessionOperations := dbo.NewSessionOperations(cfg.database)
	sessionId := c.Cookies("sessionId")

	if sessionId == "" {
		statusCode := fiber.StatusUnauthorized
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session is invalid"})
	}

	session, err := sessionOperations.GetSession(context.Background(), sessionId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("Session middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	if session == nil {
		c.ClearCookie("sessionId")

		statusCode := fiber.StatusForbidden
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: "Your session has expired"})
	}

	newSession, err := sessionOperations.UpdateSession(context.Background(), sessionId)

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
