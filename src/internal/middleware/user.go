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
	GetOneBy(map[string]interface{}) (*pg.User, *db.DBError)
}

func User(userProvider UserProvider, logger *logrus.Logger) Middleware {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		user, dbErr := userProvider.GetOneBy(map[string]interface{}{"id": userId})

		if dbErr != nil {
			statusCode := httputils.InternalError
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("User middleware")
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, dbErr.Message))
		}

		if user == nil {
			statusCode := httputils.NotFound
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "User not found."))
		}

		c.Locals("user", user)

		return c.Next()
	}

}
