package middleware

import (
	"fmt"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(key string, value interface{}) (*pg.User, error)
}

func User(userProvider UserProvider, logger *logrus.Logger) Middleware {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		user, err := userProvider.GetOneBy("id", userId)

		fmt.Println("user, err")
		fmt.Println(user, err)

		if err != nil {
			statusCode := fiber.StatusInternalServerError
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": err.Error()}).Info("User middleware")
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.InternalError.Error()})
		}

		if user == nil {
			statusCode := fiber.StatusNotFound
			return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: dto.NotFoundError.Error()})
		}

		c.Locals("user", user)

		return c.Next()
	}

}
