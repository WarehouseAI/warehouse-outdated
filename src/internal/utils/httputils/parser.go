package httputils

import (
	"github.com/gofiber/fiber/v2"
)

func Parser(structure interface{}) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(structure); err != nil {
			statusCode := BadRequest
			return c.Status(statusCode).JSON(NewErrorResponse(statusCode, "Invalid request body"))
		}

		return nil
	}
}
