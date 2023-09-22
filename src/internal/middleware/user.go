package middleware

import (
	"fmt"
	"time"
	dbo "warehouse/src/internal/db/operations"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserMiddleware struct {
	database *gorm.DB
	logger   *logrus.Logger
}

func NewUserMiddleware(database *gorm.DB, logger *logrus.Logger) *UserMiddleware {
	return &UserMiddleware{
		database: database,
		logger:   logger,
	}
}

func (cfg *UserMiddleware) User(c *fiber.Ctx) error {
	userOperations := dbo.NewUserOperations(cfg.database)
	userId := c.Locals("userId")
	user, err := userOperations.GetOneBy("id", userId)

	fmt.Println("user, err")
	fmt.Println(user, err)

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
