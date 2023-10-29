package httputils

import (
	"github.com/gofiber/fiber/v2"
)

func Parser(structure interface{}) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(structure); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(Abort, "Invalid request body"))
		}

		return nil
	}
}
