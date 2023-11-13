package middleware

import (
	"time"
	e "warehouseai/internal/errors"
	"warehouseai/user/dataservice/userdata"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func User(logger *logrus.Logger, db *userdata.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		user, dbErr := db.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

		if dbErr != nil {
			statusCode := e.HttpInternalError
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("User middleware")
			return c.Status(statusCode).JSON(e.NewErrorResponse(statusCode, dbErr.Message))
		}

		if user == nil {
			statusCode := e.HttpNotFound
			return c.Status(statusCode).JSON(e.NewErrorResponse(statusCode, "User not found."))
		}

		c.Locals("user", user)

		return c.Next()
	}

}
