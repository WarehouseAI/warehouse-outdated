package middleware

import (
	"time"
	db "warehouse/src/internal/database"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserProvider interface {
	GetOneBy(key string, value interface{}) (*pg.User, *db.DBError)
}

func User(userProvider UserProvider, logger *logrus.Logger) Middleware {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		user, dbErr := userProvider.GetOneBy("id", userId)

		if dbErr != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("User middleware")
			return c.Status(fiber.StatusInternalServerError).JSON(httputils.NewErrorResponse(httputils.ServerError, dbErr.Message))
		}

		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(httputils.NewErrorResponse(httputils.Abort, "User not found."))
		}

		c.Locals("user", user)

		return c.Next()
	}

}
