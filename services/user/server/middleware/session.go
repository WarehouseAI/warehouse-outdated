package middleware

import (
	"warehouseai/user/adapter"
	e "warehouseai/user/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Session(logger *logrus.Logger, auth adapter.AuthGrpcInterface) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("sessionId")

		if sessionId == "" {
			return c.Status(e.HttpUnauthorized).JSON(e.NewErrorResponse(e.HttpUnauthorized, "Empty session key."))
		}

		userId, newSessionId, authErr := auth.Authenticate(sessionId)

		if authErr != nil {
			return c.Status(authErr.ErrorCode).JSON(authErr)
		}

		c.ClearCookie("sessionId")
		c.Cookie(&fiber.Cookie{
			Name:     "sessionId",
			Value:    newSessionId,
			SameSite: fiber.CookieSameSiteNoneMode,
			Secure:   true,
		})

		c.Locals("userId", userId)
		return c.Next()
	}
}
