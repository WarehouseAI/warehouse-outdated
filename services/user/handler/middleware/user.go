package middleware

import (
	"time"
	"warehouseai/user/dataservice"
	errs "warehouseai/user/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func User(logger *logrus.Logger, db *dataservice.Database) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		user, dbErr := db.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

		if dbErr != nil {
			statusCode := errs.HttpInternalError
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("User middleware")
			return c.Status(statusCode).JSON(errs.NewErrorResponse(statusCode, dbErr.Message))
		}

		if user == nil {
			statusCode := errs.HttpNotFound
			return c.Status(statusCode).JSON(errs.NewErrorResponse(statusCode, "User not found."))
		}

		c.Locals("user", user)

		return c.Next()
	}

}
