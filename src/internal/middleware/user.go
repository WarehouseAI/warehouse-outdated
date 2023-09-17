package middleware

import (
	"time"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserMiddleware struct {
	userOperations dbo.UserDatabaseOperations
	logger         *logrus.Logger
}

func NewUserMiddleware(userOperations dbo.UserDatabaseOperations, logger *logrus.Logger) *UserMiddleware {
	return &UserMiddleware{
		userOperations: userOperations,
		logger:         logger,
	}
}

func (cfg *UserMiddleware) User(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	user, err := cfg.userOperations.GetOneBy("id", userId)

	if err != nil {
		statusCode := fiber.StatusInternalServerError
		cfg.logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("User middleware")
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
	}

	if user == nil {
		statusCode := fiber.StatusNotFound
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.NotFoundError.Error()})
	}

	c.Locals("user", user)

	return c.Next()
}
