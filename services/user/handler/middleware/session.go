package middleware

import (
	"warehouseai/user/adapter/grpc/client"
	"warehouseai/user/adapter/grpc/client/auth"
	errs "warehouseai/user/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Session(logger *logrus.Logger, authReceiver *client.GrpcClient) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("sessionId")

		if sessionId == "" {
			return c.Status(errs.HttpUnauthorized).JSON(errs.NewErrorResponse(errs.HttpUnauthorized, "Empty session key."))
		}

		userId, authErr := auth.Authenticate(authReceiver, sessionId)

		if authErr != nil {
			return c.Status(authErr.ErrorCode).JSON(authErr)
		}

		c.Locals("userId", userId)
		return c.Next()
	}
}
