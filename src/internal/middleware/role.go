package middleware

import (
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Role(preferRole pg.UserRole, logger *logrus.Logger) Middleware {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*pg.User)

		if user == nil {
			statusCode := httputils.NotFound
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "User not found."))
		}

		if user.Role != preferRole {
			statusCode := httputils.Forbidden
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "You're not allowed to do this."))
		}

		return c.Next()
	}

}
