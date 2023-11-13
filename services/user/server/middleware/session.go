package middleware

import (
	e "warehouseai/internal/errors"
	"warehouseai/user/adapter/grpc"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Session(logger *logrus.Logger, auth grpc.AuthGrpcInterface) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("sessionId")

		if sessionId == "" {
			return c.Status(e.HttpUnauthorized).JSON(e.NewErrorResponse(e.HttpUnauthorized, "Empty session key."))
		}

		userId, authErr := auth.Authenticate(sessionId)

		if authErr != nil {
			return c.Status(authErr.ErrorCode).JSON(authErr)
		}

		c.Locals("userId", userId)
		return c.Next()
	}
}
